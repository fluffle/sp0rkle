package main

// sp0rkle will live again!

import (
	"flag"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/lib/db"
	"github.com/fluffle/sp0rkle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/sp0rkle/drivers/calcdriver"
	"github.com/fluffle/sp0rkle/sp0rkle/drivers/decisiondriver"
	"github.com/fluffle/sp0rkle/sp0rkle/drivers/factdriver"
	"github.com/fluffle/sp0rkle/sp0rkle/drivers/netdriver"
	"github.com/fluffle/sp0rkle/sp0rkle/drivers/quotedriver"
	"github.com/fluffle/sp0rkle/sp0rkle/drivers/seendriver"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

var (
	host *string = flag.String("host", "", "IRC server to connect to.")
	port *string = flag.String("port", "6667", "Port to connect to.")
	ssl *bool = flag.Bool("ssl", false, "Use SSL when connecting.")
	nick *string = flag.String("nick", "sp0rklf",
		"Name of bot, defaults to 'sp0rklf'")
	channels *string = flag.String("channels", "#sp0rklf",
		"Comma-separated list of channels to join, defaults to '#sp0rklf'")
	rebuilder *string = flag.String("rebuilder", "",
		"Nick[:password] to accept rebuild command from.")
)

func main() {
	flag.Parse()
	log := logging.NewFromFlags()
	reg := event.NewRegistry()

	if *host == "" {
		//Don't call log.Fatal as we don't want a backtrace in this case
		log.Error("--host option required. \nOptions are:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Connect to mongo
	db, err := db.Connect("localhost")
	if err != nil {
		log.Fatal("mongo dial failed: %v\n", err)
	}
	defer db.Session.Close()

	// Initialise the factoid driver (which currently acts as a plugin mgr too).
	fd := factdriver.FactoidDriver(db, log)

	// Configure IRC client
	irc := client.Client(*nick, "boing", "not really sp0rkle", reg, log)
	irc.SSL = *ssl

	// Initialise bot state
	bot := bot.Bot(irc, fd, log)
	bot.AddChannels(strings.Split(*channels, ","))
	if *rebuilder != "" {
		bot.Rebuilder(*rebuilder)
	}

	// Add drivers
	bot.AddDriver(bot)
	bot.AddDriver(fd)
	bot.AddDriver(calcdriver.CalcDriver(log))
	bot.AddDriver(decisiondriver.DecisionDriver())
	bot.AddDriver(quotedriver.QuoteDriver(db, log))
	bot.AddDriver(netdriver.NetDriver(log))
	bot.AddDriver(seendriver.SeenDriver(db, log))

	// Register everything
	bot.RegisterAll()

	// Connect loop.
	hp := strings.Join([]string{*host, *port}, ":")
	quit := false
	for !quit {
		if err := irc.Connect(hp); err != nil {
			log.Fatal("Connection error: %s", err)
		}
		quit = <-bot.Quit
	}
	if bot.ReExec() {
		// Calling syscall.Exec probably means deferred functions won't get
		// called, so disconnect from mongodb first for politeness' sake.
		db.Session.Close()
		// If sp0rkle was run from PATH, we need to do that lookup manually.
		fq, _ := exec.LookPath(os.Args[0])
		log.Warn("Re-executing sp0rkle with args '%v'.", os.Args)
		err := syscall.Exec(fq, os.Args, os.Environ())
		if err != nil {
			// hmmmmmm
			log.Fatal("Couldn't re-exec sp0rkle: %v", err)
		}
	}
}
