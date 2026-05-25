package main

// sp0rkle will live again!

import (
	"context"
	_ "expvar"
	"flag"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fluffle/goirc/logging/golog"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/bot"
	"github.com/fluffle/sp0rkle/db/conf"
	"github.com/fluffle/sp0rkle/collections/factoids"
	"github.com/fluffle/sp0rkle/collections/karma"
	"github.com/fluffle/sp0rkle/collections/markov"
	"github.com/fluffle/sp0rkle/collections/pushes"
	"github.com/fluffle/sp0rkle/collections/quotes"
	"github.com/fluffle/sp0rkle/collections/reminders"
	"github.com/fluffle/sp0rkle/collections/seen"
	"github.com/fluffle/sp0rkle/collections/stats"
	"github.com/fluffle/sp0rkle/collections/urls"
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
	httpPort        = flag.String("http", ":6666", "Port to serve HTTP requests on.")
	boltDB          = flag.String("boltdb", "sp0rkle.boltdb", "Path to boltdb file.")
	backupDir       = flag.String("backup_dir", "backup", "Where to write BoltDB backups to")
	backupEvery     = flag.Duration("backup_every", 24*time.Hour, "How often to write backups.")
	timezone        = flag.String("timezone", "Europe/London", "Default timezone for date/time.")
	backupOnStartup = flag.Bool("backup_on_startup", true, "Run initial backup on startup.")
	fsckOnStartup   = flag.Bool("fsck_on_startup", true, "Run Fsck on all collections after registration.")
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

	// Connect to database
	dbInst, err := db.New(*boltDB, *backupDir, *backupEvery)
	if err != nil {
		logging.Fatal("Unable to open BoltDB file %q: %v", *boltDB, err)
	}
	defer dbInst.Close()

	// Initialize collections
	karmaC := karma.New(db.RegisterKeyed[*karma.Karma](dbInst))
	quotesC := quotes.New(db.RegisterIndexed[*quotes.Quote](dbInst))
	seenC := seen.New(db.RegisterIndexed[*seen.Nick](dbInst))
	urlsC := urls.New(db.RegisterIndexed[*urls.Url](dbInst))
	statsC := stats.New(db.RegisterIndexed[*stats.NickStat](dbInst))
	pushesC := pushes.New(db.RegisterIndexed[*pushes.State](dbInst))
	factoidsC := factoids.New(db.RegisterIndexed[*factoids.Factoid](dbInst))
	remindersC := reminders.New(db.RegisterIndexed[*reminders.Reminder](dbInst))
	markovC := markov.New(dbInst.DB())

	// Wire up the bot's ignore namespace
	config := conf.New(db.RegisterKeyed[*conf.Entry](dbInst))

	// Initialise bot state
	ctx := context.Background()
	botInst := bot.New(ctx, config)

	// Register drivers
	calcdriver.New(botInst, config)
	decisiondriver.New(botInst)
	factdriver.New(botInst, factoidsC)
	karmadriver.New(botInst, karmaC)
	markovdriver.New(botInst, markovC, config)
	netdriver.New(botInst, remindersC, pushesC, config)
	quotedriver.New(botInst, quotesC)
	reminddriver.New(botInst, remindersC, pushesC, config)
	seendriver.New(botInst, seenC)
	statsdriver.New(botInst, statsC)
	urldriver.New(botInst, urlsC)

	// Start up the HTTP server
	go http.ListenAndServe(*httpPort, nil)

	// Set up a signal handler to shut things down gracefully.
	go func() {
		called := new(int32)
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		for range sigint {
			if atomic.AddInt32(called, 1) > 1 {
				dbInst.Close()
				logging.Fatal("Recieved multiple interrupts, dying.")
			}
			botInst.Shutdown()
		}
	}()

	// Connect the bot to IRC and wait; reconnects are handled automatically.
	<-botInst.Connect()
	logging.Info("Shutting down cleanly.")
}
