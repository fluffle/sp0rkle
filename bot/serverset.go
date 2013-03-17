package bot

import (
	"flag"
	"fmt"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
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
			case <-time.After(*pause * time.Second):
			}
		}
	}
	// Decrement wait group when connectLoop exits.
	s.wg.Done()
}

type ServerSet interface {
	client.Handler
	Connect() bool
	HandleAll(event string, h client.Handler)
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
		// May not need to be buffered but seems sensible here
		rebuild: make(chan bool, 1),
		wg:      &sync.WaitGroup{},
	}
	for _, hostport := range list {
		// Configure IRC client
		cfg := client.NewConfig(*nick, "boing", "slowly becoming sp0rkle")
		cfg.Flood = true
		cfg.SSL = *ssl
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

func (ss *serverSet) Connect() bool {
	for _, server := range ss.servers {
		go server.connectLoop()
		ss.wg.Add(1)
	}
	// Wait for all connectLoops to terminate before reading the rebuild flag.
	ss.wg.Wait()
	return <-ss.rebuild
}

func (ss *serverSet) Shutdown(rebuild bool) {
	message := "Shutting down."
	if rebuild {
		message = "Restarting with new build."
	}
	logging.Warn(message)
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
	// Buffered write so that Shutdown can finish quickly without waiting for
	// all connectLoops to die and Connect to read.
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

// Catch, log, and complain about panics in handlers.
func unfail(conn *client.Conn, line *client.Line) {
	if err := recover(); err != nil {
		_, f, l, _ := runtime.Caller(2)
		i := strings.Index(f, "sp0rkle/")
		logging.Error("panic at %s:%d: %v", f[i:], l, err)
		conn.Privmsg(line.Target(), fmt.Sprintf(
			"panic at %s:%d: %v", f[i:], l, err))
	}
}
