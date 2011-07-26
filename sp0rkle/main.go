package main

// sp0rkle will live again!

import (
	"flag"
	"fmt"
	"github.com/fluffle/goirc/client"
	"lib/db"
	"lib/factoids"
	"log"
	"strings"
)

var host *string = flag.String("host", "", "IRC server to connect to.")
var port *string = flag.String("port", "6667", "Port to connect to.")
var ssl  *bool   = flag.Bool("ssl", false, "Use SSL when connecting.")

type botState struct {
	fc   *factoids.FactoidCollection
	quit chan bool
}

var handlers = map[string]func(*client.Conn, *client.Line) {
	"connected":    h_connected,
	"privmsg":      h_privmsg,
	"disconnected": h_disconnected,
}

func main() {
	flag.Parse()

	if *host == "" {
		log.Fatalln("need a --host, retard")
	}
	
	// Connect to mongo and initialise state
	db, err := db.Connect("localhost")
	if err != nil {
		log.Fatalf("mongo dial failed: %v\n", err)
	}
	defer db.Session.Close()
	fc, err := factoids.Collection(db)
	if err != nil {
		log.Fatalf("factoid collection failed: %v\n", err)
	}
	state := &botState{fc: fc, quit: make(chan bool)}
	
	// Configure IRC client
	irc := client.New("sp0rklf", "boing", "not really sp0rkle")
	irc.SSL = *ssl
	irc.State = state
	
	for event, handler := range handlers {
		irc.AddHandler(event, handler)
	}
	
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
		case quit = <-state.quit:
			log.Println("Shutting down...")
		}
	}
}

func h_connected(irc *client.Conn, line *client.Line) {
	log.Println("Connected, joining #sp0rklf...")
	irc.Join("#sp0rklf")
}

func h_privmsg(irc *client.Conn, line *client.Line) {
	state := irc.State.(*botState)
	text := line.Args[1]
	if fact := state.fc.GetPseudoRand(strings.ToLower(text)); fact != nil {
		switch fact.Type {
		case factoids.F_ACTION:
			irc.Action(line.Args[0], fact.Value)
		default:
			irc.Privmsg(line.Args[0], fact.Value)
		}
	}
}
	
func h_disconnected(irc *client.Conn, line *client.Line) {
	log.Println("Disconnected...")
	state := irc.State.(*botState)
	state.quit <- true
}

	
