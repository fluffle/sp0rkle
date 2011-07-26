package main

// sp0rkle will live again!

import (
	"flag"
	"fmt"
	"github.com/fluffle/goirc/client"
	"lib/db"
	"lib/factoids"
	"lib/util"
	"log"
	"strings"
)

var host *string = flag.String("host", "", "IRC server to connect to.")
var port *string = flag.String("port", "6667", "Port to connect to.")
var ssl  *bool   = flag.Bool("ssl", false, "Use SSL when connecting.")

// The bot is called sp0rkle...
type sp0rkle struct {
	// and it has a Factoid Collection...
	fc   *factoids.FactoidCollection

	// and we need to kill it occasionally.
	quit chan bool
}

var handlers = map[string]func(*client.Conn, *client.Line) {
	"connected":    h_connected,
	"privmsg":      h_privmsg,
	"action":       h_action,
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
	bot := &sp0rkle{fc: fc, quit: make(chan bool)}
	
	// Configure IRC client
	irc := client.New("sp0rklf", "boing", "not really sp0rkle")
	irc.SSL = *ssl
	irc.State = bot
	
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
		case quit = <-bot.quit:
			log.Println("Shutting down...")
		}
	}
}

func h_connected(irc *client.Conn, line *client.Line) {
	log.Println("Connected, joining #sp0rklf...")
	irc.Join("#sp0rklf")
}

func h_privmsg(irc *client.Conn, line *client.Line) {
	bot := getState(irc)
	key := strings.ToLower(strings.TrimSpace(line.Args[1]))
	key = util.RemovePrefixedNick(key, irc.Me.Nick)
	if fact := bot.fc.GetPseudoRand(key); fact != nil {
		switch fact.Type {
		case factoids.F_ACTION:
			irc.Action(line.Args[0], fact.Value)
		default:
			irc.Privmsg(line.Args[0], fact.Value)
		}
	}
}

func h_action(irc *client.Conn, line *client.Line) {
	bot := getState(irc)
	key := strings.ToLower(strings.TrimSpace(line.Args[1]))
	var fact *factoids.Factoid
	
	if fact = bot.fc.GetPseudoRand(key); fact == nil {
		// Support sp0rkle's habit of stripping off it's own nick
		// but only for actions, not privmsgs.
		if strings.HasSuffix(key, irc.Me.Nick) {
			key = strings.TrimSpace(key[:len(key)-len(irc.Me.Nick)])
			if fact = bot.fc.GetPseudoRand(key); fact == nil {
				return
			}
		}
	}
	switch fact.Type {
	case factoids.F_ACTION:
		irc.Action(line.Args[0], fact.Value)
	default:
		irc.Privmsg(line.Args[0], fact.Value)
	}
}

func h_disconnected(irc *client.Conn, line *client.Line) {
	log.Println("Disconnected...")
	bot := getState(irc)
	bot.quit <- true
}

func getState(irc *client.Conn) *sp0rkle {
	return irc.State.(*sp0rkle)
}
