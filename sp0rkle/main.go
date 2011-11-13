package main

// sp0rkle will live again!

import (
	"flag"
	"github.com/fluffle/goevent/event"
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/golog/logging"
	"lib/db"
	"sp0rkle/bot"
	"sp0rkle/drivers/decisiondriver"
	"sp0rkle/drivers/factdriver"
	"strings"
)

var host *string = flag.String("host", "", "IRC server to connect to.")
var port *string = flag.String("port", "6667", "Port to connect to.")
var ssl *bool = flag.Bool("ssl", false, "Use SSL when connecting.")
var nick *string = flag.String("nick", "sp0rklf",
	"Name of bot, defaults to 'sp0rklf'")
var channel *string = flag.String("channel", "#sp0rklf",
	"Channel to join, defaults to '#sp0rklf'")


func main() {
	flag.Parse()
	log := logging.NewFromFlags()
	reg := event.NewRegistry()

	if *host == "" {
		log.Fatal("need a --host, retard")
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
	bot.AddChannel(*channel)

	// Add drivers
	bot.AddDriver(bot)
	bot.AddDriver(fd)
	bot.AddDriver(decisiondriver.DecisionDriver())

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
}
