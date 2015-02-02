package netdriver

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/collections/conf"
	"net"
	"strconv"
	"strings"
	"time"
)

type mcStatus struct {
	motd       string
	nump, maxp string
	players    []string
	version    string
}

const (
	playerdata = "\x00\x00\x01player_\x00\x00"
	mcServer   = "server"
	mcFreq     = "freq"
	mcChan     = "chan"
)

var (
	mcConf      conf.Namespace
	mcHandshake = []byte("\xfe\xfd\x09\x00\x00\x00\x00")
	mcGetStatus = []byte("\xfe\xfd\x00\x00\x00\x00\x00")
)

func mcSet(ctx *bot.Context) {
	kv := strings.Fields(ctx.Text())
	if len(kv) < 2 {
		ctx.ReplyN("I need a key and a value.")
		return
	}
	switch kv[0] {
	case mcServer:
		mcConf.String(mcServer, kv[1])
	case mcChan:
		if !strings.HasPrefix(kv[1], "#") {
			ctx.ReplyN("Channel '%s' doesn't start with #.", kv[1])
			return
		}
		mcConf.String(mcChan, kv[1])
	case mcFreq:
		freq, err := strconv.Atoi(kv[1])
		if err != nil {
			ctx.ReplyN("Couldn't convert '%s' to an integer.", kv[1])
			return
		}
		mcConf.Int(mcFreq, freq)
	default:
		ctx.ReplyN("Valid keys are: %s, %s, %s", mcServer, mcFreq, mcChan)
	}
}

func (mcs *mcStatus) Poll(ctxs []*bot.Context) {
	srv := mcConf.String(mcServer)
	logging.Debug("polling minecraft server at %s", srv)
	st, err := pollServer(srv)
	if err != nil {
		logging.Error("minecraft poll failed: %v", err)
		return
	}
	*mcs = *st
	for _, ctx := range ctxs {
		ctx.Topic(mcConf.String(mcChan))
	}
}

func (mcs *mcStatus) Start() { /* empty */ }
func (mcs *mcStatus) Stop()  { /* empty */ }
func (mcs *mcStatus) Tick() time.Duration {
	return time.Duration(mcConf.Int(mcFreq)) * time.Minute
}

func (mcs *mcStatus) Topic(ctx *bot.Context) {
	ch := mcConf.String(mcChan)
	if ctx.Args[1] != ch {
		return
	}
	topic := ctx.Text()
	if idx := strings.Index(topic, " || "); idx == -1 {
		topic = ""
	} else {
		topic = topic[idx:]
	}
	players := ""
	if len(mcs.players) > 0 {
		players = ": " + strings.Join(mcs.players, ", ")
	}
	topic = fmt.Sprintf("%s %s v%s [%s/%s%s]%s", mcs.motd,
		mcConf.String(mcServer), mcs.version, mcs.nump, mcs.maxp, players, topic)
	if topic != ctx.Text() {
		ctx.Topic(ch, topic)
	}
}

// Conducts a single poll of a server.
func pollServer(server string) (*mcStatus, error) {
	nc, err := net.Dial("udp", server)
	if err != nil {
		return nil, err
	}
	defer nc.Close()
	// If we've not finished doing this lot after a minute, bail out.
	if err = nc.SetDeadline(time.Now().Add(time.Minute)); err != nil {
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
		if n != b.Len() {
			return nil, fmt.Errorf("short write in status")
		}
		return nil, err
	}

	// Read and parse response
	if n, err = nc.Read(buf); err != nil {
		return nil, err
	}
	data := string(buf[11 : n-1]) // skip splitnum + int, and strip trailing \x00
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
	if idx+len(playerdata) < len(data) {
		st.players = strings.Split(data[idx+len(playerdata):len(data)-1], "\x00")
	}
	return st, nil
}
