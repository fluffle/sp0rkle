package bot

import (
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
)

var (
	nick *string = flag.String("nick", "sp0rklf",
		"Name of bot, defaults to 'sp0rklf'")
	servers *string = flag.String("servers", "",
		"Comma-separated list of IRC servers to connect to.")
	ssl *bool = flag.Bool("ssl", false,
		"Use SSL when connecting to servers.")
	pause *time.Duration = flag.Duration("pause", 300*time.Second,
		"Wait time between server reconnection attempts.")
)

type server struct {
	*client.Conn
	hostport string
	shutdown bool

	wg   *sync.WaitGroup
	wait chan struct{}
}

func (s *server) connectLoop() {
	for {
		logging.Info("Connecting to %s.", s.hostport)
		if err := s.Connect(); err == nil {
			// Wait here for a disconnect signal
			<-s.wait
			if s.shutdown {
				break
			}
		} else {
			logging.Error("Connection error: %s", err)
			select {
			case <-s.wait:
				// If we are waiting for a reconnect to this server
				// and someone calls Shutdown, we need to shut down
				if s.shutdown {
					break
				}
			case <-time.After(*pause):
			}
		}
	}
	// Decrement wait group when connectLoop exits.
	s.wg.Done()
}

type ServerSet interface {
	client.Handler
	Connect() chan bool
	HandleAll(event string, h client.Handler)
	HandleAllBG(event string, h client.Handler)
	Shutdown(rebuild bool)
}

type serverSet struct {
	servers map[*client.Conn]*server
	wg      *sync.WaitGroup
	rebuild chan bool
}

func newServerSet() *serverSet {
	list := strings.Split(*servers, ",")
	if len(list) == 0 {
		// Don't call logging.Fatal as we don't want a backtrace in this case
		logging.Error("--server option required. \nOptions are:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	ss := &serverSet{
		servers: make(map[*client.Conn]*server),
		rebuild: make(chan bool),
		wg:      &sync.WaitGroup{},
	}
	for _, hostport := range list {
		// Configure IRC client
		cfg := client.NewConfig(*nick, "boing", "slowly becoming sp0rkle")
		cfg.Flood = true
		if *ssl {
			cfg.SSL = true
			cfg.SSLConfig = &tls.Config{
				ServerName: strings.Split(hostport, ":")[0],
			}
		}
		cfg.Recover = unfail
		cfg.Server = hostport
		conn := client.Client(cfg)
		ss.servers[conn] = &server{
			Conn:     conn,
			hostport: hostport,
			wg:       ss.wg,
			wait:     make(chan struct{}),
		}
	}
	ss.HandleAll(client.DISCONNECTED, ss)
	return ss
}

func (ss *serverSet) Connect() chan bool {
	for _, server := range ss.servers {
		go server.connectLoop()
		ss.wg.Add(1)
	}
	return ss.rebuild
}

func (ss *serverSet) Shutdown(rebuild bool) {
	message := "Shutting down."
	if rebuild {
		message = "Restarting with new build."
	}
	logging.Info(message)
	for _, server := range ss.servers {
		server.shutdown = true
		if server.Connected() {
			// If we're connected to this server, disconnect gracefully
			// and send wait strobe to connectLoop from Handle()
			server.Quit(message)
		} else {
			// If we're not connected to this server, connectLoop
			// is waiting in the select{} to reconnect, so strobe now.
			server.wait <- struct{}{}
		}
	}
	// Wait for all connectLoops to terminate
	ss.wg.Wait()
	ss.rebuild <- rebuild
}

// serverSet's Handle() deals with disconnects from individual servers
func (ss *serverSet) Handle(conn *client.Conn, line *client.Line) {
	server := ss.servers[conn]
	logging.Info("Disconnected from %s...", server.hostport)
	server.wait <- struct{}{}
}

// HandleAll() registers Handlers with all the servers in the set
func (ss *serverSet) HandleAll(ev string, h client.Handler) {
	for conn, _ := range ss.servers {
		conn.Handle(ev, h)
	}
}

// HandleAllBG() registers background Handlers with all the servers in the set
func (ss *serverSet) HandleAllBG(ev string, h client.Handler) {
	for conn, _ := range ss.servers {
		conn.HandleBG(ev, h)
	}
}

// Catch, log, and complain about panics in handlers.
func unfail(conn *client.Conn, line *client.Line) {
	if err := recover(); err != nil {
		// Depth 4 is where our code usually starts.
		// But if the panic is somewhere in the depths of the standard
		// library or dependency code it's helpful to know what
		// our code was up to at the time too, so walk up the stack.
		callers := make([]uintptr, 20)
		depth := 4
		n := runtime.Callers(depth, callers)
		frames := runtime.CallersFrames(callers[:n])
		msgs := []string{fmt.Sprintf("panic: %v frames:", err)}
		for frame, ok := frames.Next(); ok; frame, ok = frames.Next() {
			msgs = append(msgs, fmt.Sprintf("(%d) %s at %s:%d",
				depth, trimPath(frame.Function), trimPath(frame.File), frame.Line))
			depth++
		}
		msg := strings.Join(msgs, ", ")
		logging.Error(msg)
		conn.Privmsg(line.Target(), msg)
	}
}

func trimPath(path string) string {
	prefixes := []string{"sp0rkle/", "fluffle/", runtime.GOROOT() + "/src/"}
	for _, prefix := range prefixes {
		trimmed, ok := trim(path, prefix)
		if ok {
			return trimmed
		}
	}
	return path
}

func trim(path string, prefix string) (string, bool) {
	i := strings.Index(path, prefix)
	if i < 0 {
		return path, false
	}
	return path[i+len(prefix):], true
}
