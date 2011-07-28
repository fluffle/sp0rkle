package main

// sp0rkle will live again!

import (
	"flag"
	"fmt"
	"github.com/fluffle/goirc/client"
	"lib/db"
	"log"
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

	if *host == "" {
		log.Fatalln("need a --host, retard")
	}

	// Initialise bot state
	bot := Bot()
	bot.AddChannel(*channel)

	// Connect to mongo
	db, err := db.Connect("localhost")
	if err != nil {
		log.Fatalf("mongo dial failed: %v\n", err)
	}
	defer db.Session.Close()

	// Add drivers
	bot.AddDriver(bot)
	bot.AddDriver(FactoidDriver(db))

	// Configure IRC client
	irc := client.New(*nick, "boing", "not really sp0rkle")
	irc.SSL = *ssl
	irc.State = bot
	bot.RegisterAll(irc.Registry)

	hp := strings.Join([]string{*host, *port}, ":")
	if err := irc.Connect(hp); err != nil {
		fmt.Printf("Connection error: %s", err)
		return
	}

	quit := false
	for !quit {
		select {
		case err := <-irc.Err:
			log.Printf("goirc error: %s\n", err)
		case quit = <-bot.quit:
			log.Println("Shutting down...")
		}
	}
}
