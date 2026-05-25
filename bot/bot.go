package bot

import (
	"context"
	"flag"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/db/conf"
)

// This is here because I'm not sure where better to put it...
var httpHost *string = flag.String("http_host", "https://s.pl0rt.org",
	"Hostname for HTTP paths served by bot.")

func HttpHost() string {
	return *httpHost
}

// CommandRegistry defines the interface for registering handlers, commands,
// rewrites, and pollers. *Bot implements this interface so it can be passed
// to drivers and collections for dependency injection.
type CommandRegistry interface {
	Handle(fn HandlerFunc, events ...string)
	HandleBG(fn HandlerFunc, events ...string)
	Command(fn HandlerFunc, prefix, help string)
	Rewrite(fn RewriteFunc)
	Poll(p Poller)
	Ctx() context.Context
}

// Bot holds all bot state and provides methods for connecting,
// shutdown, handler registration, and command management.
type Bot struct {
	mu         sync.Mutex
	isConnected  bool
	servers    ServerSet
	rewriters  RewriteSet
	commands   CommandSet
	pollers    PollerSet
	filters    *FilterPipeline
	ctx        context.Context
	ignoreNs   conf.Namespace
}

// Ensure *Bot implements CommandRegistry.
var _ CommandRegistry = (*Bot)(nil)

// New creates a new Bot instance with the given context.
func New(ctx context.Context, config *conf.Registry) *Bot {
	b := &Bot{
		ctx:        ctx,
		ignoreNs:   config.Ns("ignore"),
		servers:   newServerSet(),
		rewriters: newRewriteSet(),
		filters:   &FilterPipeline{},
	}
	b.commands = newCommandSet(b.filters, b.rewriters)
	b.pollers = newPollerSet(b.rewriters)

	// This is a special handler that dispatches commands from the command set
	b.servers.HandleAll(client.PRIVMSG, b.commands)

	// The poller set handles these two to start and stop registered pollers
	b.servers.HandleAll(client.CONNECTED, b.pollers)
	b.servers.HandleAll(client.DISCONNECTED, b.pollers)

	// Internal event handlers
	b.Handle(func(ctx *Context) { b.connected(ctx) }, client.CONNECTED)
	b.Handle(func(ctx *Context) { b.shutdown(ctx) }, client.NOTICE)

	// Internal commands
	b.Command(func(ctx *Context) { b.ignore(ctx) }, "ignore", "ignore <nick>  -- "+
		"make the bot ignore <nick> completely.")
	b.Command(func(ctx *Context) { b.unignore(ctx) }, "unignore", "unignore <nick>  -- "+
		"make the bot unignore <nick> again.")

	return b
}

// Connect starts the bot's server connections and returns a channel
// that signals when the bot should rebuild or shut down.
func (b *Bot) Connect() chan bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.isConnected {
		logging.Warn("Already connected to servers.")
	}
	b.isConnected = true
	b.filters.Add(&nickIgnoreFilter{ns: b.ignoreNs})
	return b.servers.Connect()
}

// Shutdown disconnects the bot from all servers.
func (b *Bot) Shutdown() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if !b.isConnected {
		logging.Warn("Not connected to servers.")
	}
	b.isConnected = false
	b.filters = &FilterPipeline{}
	b.servers.Shutdown(false)
}

// Handle registers a handler function for the given IRC events.
func (b *Bot) Handle(fn HandlerFunc, events ...string) {
	h := &handler{fn: fn, filters: b.filters, rws: b.rewriters}
	for _, ev := range events {
		b.servers.HandleAll(ev, h)
	}
}

// HandleBG registers a background handler function for the given IRC events.
func (b *Bot) HandleBG(fn HandlerFunc, events ...string) {
	h := &handler{fn: fn, filters: b.filters, rws: b.rewriters}
	for _, ev := range events {
		b.servers.HandleAllBG(ev, h)
	}
}

// Command registers a command with the given prefix and help text.
func (b *Bot) Command(fn HandlerFunc, prefix, help string) {
	b.commands.Add(&command{fn, help}, prefix)
}

// Rewrite registers a rewrite function.
func (b *Bot) Rewrite(fn RewriteFunc) {
	b.rewriters.Add(fn)
}

// Poll registers a poller.
func (b *Bot) Poll(p Poller) {
	b.pollers.Add(p)
}

// Ctx returns the bot's context.
func (b *Bot) Ctx() context.Context {
	return b.ctx
}

// GetSecret retrieves a secret value, supporting $ENV_VAR and <file_path> syntax.
func GetSecret(s string) string {
	if strings.HasPrefix(s, "$") {
		return os.ExpandEnv(s)
	} else if strings.HasPrefix(s, "<") {
		if bytes, err := ioutil.ReadFile(s[1:]); err == nil {
			return strings.TrimSuffix(string(bytes), "\n")
		}
		return ""
	}
	return s
}
