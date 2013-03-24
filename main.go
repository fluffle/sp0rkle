package main

// sp0rkle will live again!

import (
	_ "expvar"
	"flag"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db"
	"github.com/fluffle/sp0rkle/drivers/calcdriver"
	"github.com/fluffle/sp0rkle/drivers/decisiondriver"
	"github.com/fluffle/sp0rkle/drivers/factdriver"
	"github.com/fluffle/sp0rkle/drivers/karmadriver"
	"github.com/fluffle/sp0rkle/drivers/markovdriver"
	"github.com/fluffle/sp0rkle/drivers/netdriver"
	"github.com/fluffle/sp0rkle/drivers/quotedriver"
	"github.com/fluffle/sp0rkle/drivers/reminddriver"
	"github.com/fluffle/sp0rkle/drivers/seendriver"
	"github.com/fluffle/sp0rkle/drivers/statsdriver"
	"github.com/fluffle/sp0rkle/drivers/urldriver"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

var (
	httpPort *string = flag.String("http", ":6666", "Port to serve HTTP requests on.")
)

func main() {
	flag.Parse()
	logging.InitFromFlags()

	// Initialise bot state
	bot.Init()

	// Connect to mongo
	db.Init()
	defer db.Close()

	// Add drivers
	calcdriver.Init()
	decisiondriver.Init()
	factdriver.Init()
	karmadriver.Init()
	markovdriver.Init()
	netdriver.Init()
	quotedriver.Init()
	reminddriver.Init()
	seendriver.Init()
	statsdriver.Init()
	urldriver.Init()

	// Start up the HTTP server
	go http.ListenAndServe(*httpPort, nil)

	// Set up a signal handler to shut things down gracefully.
	// NOTE: net/http doesn't provide for graceful shutdown :-/
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		if syscall.SIGINT == <-sigint {
			bot.Shutdown()
		}
	}()

	// Connect the bot to IRC and wait; reconnects are handled automatically.
	// If we get true back from the bot, re-exec the (rebuilt) binary.
	if <-bot.Connect() {
		// Calling syscall.Exec probably means deferred functions won't get
		// called, so disconnect from mongodb first for politeness' sake.
		db.Close()
		// If sp0rkle was run from PATH, we need to do that lookup manually.
		fq, _ := exec.LookPath(os.Args[0])
		logging.Warn("Re-executing sp0rkle with args '%v'.", os.Args)
		err := syscall.Exec(fq, os.Args, os.Environ())
		if err != nil {
			// hmmmmmm
			logging.Fatal("Couldn't re-exec sp0rkle: %v", err)
		}
	}
	logging.Info("Shutting down cleanly.")
}
