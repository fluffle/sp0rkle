package netdriver

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"net"
	"strconv"
	"strings"
	"time"
)

type mcStatus struct {
	motd string
	nump, maxp string
	players []string
	version string
}

const playerdata = "\x00\x00\x01player_\x00\x00"

var (
	status chan *mcStatus
	quit   chan struct{}
	mcHandshake = []byte("\xfe\xfd\x09\x00\x00\x00\x00")
	mcGetStatus = []byte("\xfe\xfd\x00\x00\x00\x00\x00")
	mcServer    = flag.String("mc_server", "", "Minecraft server to poll")
	mcPollFreq  = flag.Duration("mc_poll_freq", 5*time.Minute,
		"How regularly to poll server")
	mcChan      = flag.String("mc_chan", "#minecraft",
		"Channel whose topic poller should keep updated")
)

func mcStartPoller(ctx *bot.Context) {
	if *mcServer == "" {
		return
	}
	status = make(chan *mcStatus, 1)
	quit = make(chan struct{})
	t := time.NewTicker(*mcPollFreq)
	go func() {
		for {
			select {
			case <-t.C:
				logging.Debug("polling minecraft server at %s", *mcServer)
				if st, err := poll(*mcServer); err == nil {
					status <- st
					ctx.Topic(*mcChan)
				} else {
					logging.Error("poll failed: %v", err)
				}
			case <-quit:
				close(status)
				return
			}
		}
	}()
}

func mcStopPoller(ctx *bot.Context) {
	if *mcServer == "" {
		return
	}
	close(quit)
	for _ = range status {}
}

func mcChanTopic(ctx *bot.Context) {
	if ctx.Args[1] != *mcChan {
		return
	}
	logging.Debug("got minecraft topic: %#v", ctx)
	var st *mcStatus
	select {
	case st = <-status:
	default:
		logging.Debug("skipping read for empty chan")
		return
	}
	topic := ctx.Text()
	if idx := strings.Index(topic, " || "); idx == -1 {
		topic = ""
	} else {
		topic = topic[idx:]
	}
	players := ""
	if len(st.players) > 0 {
		players = ": " + strings.Join(st.players, ", ")
	}
	topic = fmt.Sprintf("%s %s v%s [%s/%s%s]%s", st.motd, *mcServer,
		st.version, st.nump, st.maxp, players, topic)
	if topic != ctx.Text() {
		ctx.Topic(*mcChan, topic)
	}
}

func poll(server string) (*mcStatus, error) {
	nc, err := net.Dial("udp", server)
	if err != nil {
		return nil, err
	}

	// Send initial handshake
	var n int
	if n, err = nc.Write(mcHandshake); err != nil || n != len(mcHandshake) {
		if n != len(mcHandshake) {
			return nil, fmt.Errorf("short write in handshake")
		}
		return nil, err
	}

	// Read response, convert ascii integer to actual integer.
	buf := make([]byte, 1024)
	if n, err = nc.Read(buf); err != nil {
		return nil, err
	}
	handshake := string(buf[:n])
	idx := strings.Index(handshake[5:], "\x00") + 5
	if idx < 5 {
		idx = len(handshake)
	}
	challenge, err := strconv.Atoi(handshake[5:idx])
	if err != nil {
		return nil, err
	}

	// Send status request, 
	b := bytes.NewBuffer(mcGetStatus)
	binary.Write(b, binary.BigEndian, int32(challenge))
	b.WriteString("\x00\x00\x00\x00")
	if n, err = nc.Write(b.Bytes()); err != nil || n != b.Len() {
		if n != len(handshake) {
			return nil, fmt.Errorf("short write in status")
		}
		return nil, err
	}

	// Read and parse response
	if n, err = nc.Read(buf); err != nil {
		return nil, err
	}
	data := string(buf[11:n-1]) // skip splitnum + int, and strip trailing \x00
	idx = strings.Index(data, playerdata)
	if idx == -1 {
		return nil, fmt.Errorf("could not find player data")
	}
	// The first part is a list of null-terminated strings, key=>value
	kvs := strings.Split(data[:idx], "\x00")
	items := make(map[string]string)
	for i := 0; i < len(kvs); i += 2 {
		items[kvs[i]] = kvs[i+1]
	}
	st := &mcStatus{
		motd:    items["hostname"], // meh.
		nump:    items["numplayers"],
		maxp:    items["maxplayers"],
		version: items["version"],
	}
	// And the second part is the null-terminated string list of players
	if idx + len(playerdata) < len(data) {
		st.players = strings.Split(data[idx+len(playerdata):len(data)-1], "\x00")
	}
	return st, nil
}
