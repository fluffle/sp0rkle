package bot

import (
	"flag"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"sync"
)

// This is here because I'm not sure where better to put it...
var httpHost *string = flag.String("http_host", "http://sp0rk.ly",
	"Hostname for HTTP paths served by bot.")

func HttpHost() string {
	return *httpHost
}

type botData struct {
	connected bool
	servers   ServerSet
	rewriters RewriteSet
	commands  CommandSet
}

var bot *botData
var lock sync.Mutex

func Init() {
	lock.Lock()
	defer lock.Unlock()
	if bot != nil {
		return
	}

	bot = &botData{
		servers:   newServerSet(),
		commands:  newCommandSet(),
		rewriters: newRewriteSet(),
	}

	// This is a special handler that dispatches commands from the command set
	bot.servers.HandleAll(client.PRIVMSG, bot.commands)

	// These three in handlers.go
	Handle(connected, client.CONNECTED)
	Handle(rebuild, client.NOTICE)
	Handle(shutdown, client.NOTICE)

	// These two in commands.go
	Command(ignore, "ignore", "ignore <nick>  -- "+
		"make the bot ignore <nick> completely.")
	Command(unignore, "unignore", "unignore <nick>  -- "+
		"make the bot unignore <nick> again.")
}

func Connect() chan bool {
	lock.Lock()
	defer lock.Unlock()
	if bot == nil {
		logging.Fatal("Called Connect() before Init().")
	}
	if bot.connected {
		logging.Warn("Already connected to servers.")
	}
	bot.connected = true
	return bot.servers.Connect()
}

func Shutdown() {
	lock.Lock()
	defer lock.Unlock()
	if bot == nil {
		logging.Fatal("Called Shutdown() before Init().")
	}
	if !bot.connected {
		logging.Warn("Not connected to servers.")
	}
	bot.connected = false
	bot.servers.Shutdown(false)
}

func Handle(fn HandlerFunc, events ...string) {
	for _, ev := range events {
		bot.servers.HandleAll(ev, fn)
	}
}

func Command(fn HandlerFunc, prefix, help string) {
	bot.commands.Add(&command{fn, help}, prefix)
}

func Rewrite(fn RewriteFunc) {
	bot.rewriters.Add(fn)
}
