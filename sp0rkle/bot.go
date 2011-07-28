package main

import (
	"github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/event"
	"log"
)

// The bot is called sp0rkle...
const _BOT_NAME string = "sp0rkle"

type sp0rkle struct {
	// it's got a bunch of drivers that register event handlers
	drivers map[string]Driver

	//channel to join on start up
	channels []string

	// and we need to kill it occasionally.
	quit chan bool
}

func Bot() *sp0rkle {
	return &sp0rkle{
		drivers:  make(map[string]Driver),
		channels: make([]string, 0, 1),
		quit:     make(chan bool),
	}
}

func (bot *sp0rkle) RegisterHandlers(r event.EventRegistry) {
	r.AddHandler("connected", client.IRCHandler(bot_connected))
	r.AddHandler("disconnected", client.IRCHandler(bot_disconnected))
}

func (bot *sp0rkle) Name() string {
	return _BOT_NAME
}

func (bot *sp0rkle) RegisterAll(r event.EventRegistry) {
	for _, d := range bot.drivers {
		d.RegisterHandlers(r)
	}
}

func (bot *sp0rkle) AddDriver(d Driver) {
	bot.drivers[d.Name()] = d
}

func (bot *sp0rkle) AddChannel(c string) {
	bot.channels = append(bot.channels, c)
}

func bot_connected(irc *client.Conn, line *client.Line) {
	bot := getState(irc)
	for _, c := range bot.channels {
		log.Printf("Joining %s on startup.\n", c)
		irc.Join(c)
	}
}

func bot_disconnected(irc *client.Conn, line *client.Line) {
	log.Println("Disconnected...")
	bot := getState(irc)
	bot.quit <- true
}

func getState(irc *client.Conn) *sp0rkle {
	return irc.State.(*sp0rkle)
}
