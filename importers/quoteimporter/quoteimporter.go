package main

// Imports perlfu's SQLite quote database into a quotes collection

import (
	"flag"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/collections/quotes"
	"github.com/fluffle/sp0rkle/db"
	"github.com/kuroneko/gosqlite3"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var file *string = flag.String("db", "Quotes.db",
	"SQLite database to import quotes from.")

const (
	// The Quotes table columns are:
	cID = iota
	cNick
	cChannel
	cQuote
	cTime
)

func parseQuote(row []interface{}, out chan *quotes.Quote) {
	out <- &quotes.Quote{
		Quote:     row[cQuote].(string),
		QID:       int(row[cID].(int64)),
		Nick:      bot.Nick(row[cNick].(string)),
		Chan:      bot.Chan(row[cChannel].(string)),
		Accessed:  0,
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
	qc := quotes.Init()

	// A communication channel of Quotes.
	quotes := make(chan *quotes.Quote)
	rows := make(chan []interface{})

	// Function to feed rows into the rows channel.
	row_feeder := func(sth *sqlite3.Statement, row ...interface{}) {
		rows <- row
	}

	// Function to execute a query on the SQLite db.
	db_query := func(dbh *sqlite3.Database) {
		n, err := dbh.Execute("SELECT * FROM Quotes;", row_feeder)
		if err == nil {
			logging.Info("Read %d rows from database.\n", n)
		} else {
			logging.Error("DB error: %s\n", err)
		}
	}

	// Open up the quote database in a goroutine and feed rows
	// in on the input_rows channel.
	go func() {
		sqlite3.Session(*file, db_query)
		// once we've done the query, close the channel to indicate this
		close(rows)
	}()

	// Another goroutine to munge the rows into quotes.
	// This was originally done inside the SQLite callbacks, but
	// cgo or sqlite3 obscures runtime panics and makes fail happen.
	go func() {
		for row := range rows {
			parseQuote(row, quotes)
		}
		close(quotes)
	}()

	// And finally...
	count := 0
	var err error
	for quote := range quotes {
		// ... push each quote into mongo
		err = qc.Insert(quote)
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
	logging.Info("Inserted %d quotes.\n", count)
}
