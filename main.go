package main

// sp0rkle will live again!

import (
	_ "expvar"
	"flag"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fluffle/goirc/logging/golog"
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
	"github.com/fluffle/sp0rkle/util/datetime"
)

var (
	httpPort = flag.String("http", ":6666", "Port to serve HTTP requests on.")
	boltDB   = flag.String("boltdb", "sp0rkle.boltdb", "Path to boltdb file.")
	mongoDB  = flag.String("mongodb", "localhost",
		"Address of MongoDB server to connect to, defaults to localhost.")
	backupDir   = flag.String("backup_dir", "backup", "Where to write BoltDB backups to.")
	backupEvery = flag.Duration("backup_every", 24*time.Hour, "How often to write backups.")
	timezone    = flag.String("timezone", "Europe/London", "Default timezone for date/time.")
)

func main() {
	flag.Parse()
	logging.InitFromFlags()
	golog.Init()
	if err := datetime.SetTZ(*timezone); err != nil {
		logging.Fatal("Failed to set default timezone from --timezone=%q: %v", *timezone, err)
	}

	// Slightly more random than 1.
	rand.Seed(time.Now().UnixNano() * int64(os.Getpid()))

	// Initialise bot state
	bot.Init()

	// Connect to databases

	if err := db.Mongo.Init(*mongoDB); err != nil {
		logging.Fatal("Unable to connect to MongoDB at %q: %v", *mongoDB, err)
	}
	defer db.Mongo.Close()
	if err := db.Bolt.Init(*boltDB, *backupDir, *backupEvery); err != nil {
		logging.Fatal("Unable to open BoltDB file %q: %v", *boltDB, err)
	}
	defer db.Bolt.Close()

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
		called := new(int32)
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		for _ = range sigint {
			if atomic.AddInt32(called, 1) > 1 {
				logging.Fatal("Recieved multiple interrupts, dying.")
			}
			bot.Shutdown()
		}
	}()

	// Connect the bot to IRC and wait; reconnects are handled automatically.
	// If we get true back from the bot, re-exec the (rebuilt) binary.
	if <-bot.Connect() {
		// Calling syscall.Exec probably means deferred functions won't get
		// called, so disconnect from DBs first for politeness' sake.
		db.Mongo.Close()
		db.Bolt.Close()
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
