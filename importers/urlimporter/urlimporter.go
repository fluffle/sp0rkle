package main

import (
	"flag"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/collections/urls"
	"github.com/fluffle/sp0rkle/db"
	"github.com/kuroneko/gosqlite3"
	"labix.org/v2/mgo/bson"
	"net/http"
	"time"
)

var file *string = flag.String("db", "URL.db",
	"SQLite database to import URLs from.")
var check *bool = flag.Bool("check", false,
	"Check each url for continued existence with a HEAD request.")
var workq *int = flag.Int("workers", 8,
	"How many HEAD requests to run in parallel.")

const (
	// The URL table columns are:
	cNick = iota
	cChannel
	cUrl
	cTime
)

func parseUrl(row []interface{}) *urls.Url {
	return &urls.Url{
		Nick:      bot.Nick(row[cNick].(string)),
		Chan:      bot.Chan(row[cChannel].(string)),
		Url:       row[cUrl].(string),
		Timestamp: time.Unix(row[cTime].(int64), 0),
		Id:        bson.NewObjectId(),
	}
}

func main() {
	flag.Parse()
	logging.InitFromFlags()

	// Let's go find some mongo.
	db.Init()
	defer db.Close()
	uc := urls.Init()

	work := make(chan *urls.Url)
	quit := make(chan bool)
	urls := make(chan *urls.Url)
	rows := make(chan []interface{})
	failed := 0

	// If we're checking, spin up some workers
	if *check {
		for i := 1; i <= *workq; i++ {
			go func(n int) {
				count := 0
				for u := range work {
					count++
					logging.Debug("w%02d r%04d: Fetching '%s'", n, count, u.Url)
					res, err := http.Head(u.Url)
					logging.Debug("w%02d r%04d: Response '%s'", n, count, res.Status)
					if err == nil && res.StatusCode == 200 {
						urls <- u
					} else {
						failed++
					}
				}
				quit <- true
			}(i)
		}
	}

	// Function to feed rows into the rows channel.
	row_feeder := func(sth *sqlite3.Statement, row ...interface{}) {
		rows <- row
	}

	// Function to execute a query on the SQLite db.
	db_query := func(dbh *sqlite3.Database) {
		n, err := dbh.Execute("SELECT * FROM urls;", row_feeder)
		if err == nil {
			logging.Info("Read %d rows from database.\n", n)
		} else {
			logging.Error("DB error: %s\n", err)
		}
	}

	// Open up the URL database in a goroutine and feed rows
	// in on the input_rows channel.
	go func() {
		sqlite3.Session(*file, db_query)
		// once we've done the query, close the channel to indicate this
		close(rows)
	}()

	// Another goroutine to munge the rows into Urls and optionally send
	// them to the pool of checker goroutines.
	go func() {
		for row := range rows {
			u := parseUrl(row)
			if *check {
				work <- u
			} else {
				urls <- u
			}
		}
		if *check {
			// Close work channel and wait for all workers to quit.
			close(work)
			for i := 0; i < *workq; i++ {
				<-quit
			}
		}
		close(urls)
	}()

	// And finally...
	count := 0
	var err error
	for u := range urls {
		// ... push each url into mongo
		err = uc.Insert(u)
		if err != nil {
			logging.Error("Awww: %v\n", err)
		} else {
			if count%1000 == 0 {
				fmt.Printf("%d...", count)
			}
			count++
		}
	}
	fmt.Println("done.")
	if *check {
		logging.Info("Dropped %d non-200 urls.", failed)
	}
	logging.Info("Inserted %d urls.", count)
}
