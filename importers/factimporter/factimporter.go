package main

// Imports perlfu's SQLite factoid database into a factoids collection

import (
	"flag"
	"fmt"
	"github.com/fluffle/golog/logging"
	"github.com/fluffle/sp0rkle/base"
	"github.com/fluffle/sp0rkle/collections/factoids"
	"github.com/fluffle/sp0rkle/db"
	"github.com/kuroneko/gosqlite3"
	"labix.org/v2/mgo/bson"
	"strconv"
	"strings"
	"time"
)

var file *string = flag.String("db", "Facts.db",
	"SQLite database to import factoids from.")

const (
	// The Factoids table columns are:
	cKey = iota
	cRel
	cValue
	cCreated
	cCreator
	cModified
	cModifier
	cAccess
)

// Parse a row returned from SQLite into a Factoid. 
func parseFactoid(row []interface{}, out chan *factoids.Factoid) {
	values := parseMultipleValues(toString(row[cValue]))
	c := &factoids.FactoidStat{
		Nick: base.Nick(toString(row[cCreator])),
		Chan: "",
		Count: 1,
	}
	c.Timestamp, _ = parseTimestamp(row[cCreated])
	m := &factoids.FactoidStat{Chan: "", Count: 0}
	if ts, ok := parseTimestamp(row[cModified]); ok {
		m.Timestamp = ts
		m.Nick = base.Nick(toString(row[cModifier]))
		m.Count = 1
	} else {
		m.Timestamp = c.Timestamp
		m.Nick = c.Nick
	}
	p := &factoids.FactoidPerms{
		parseReadOnly(row[cAccess]),
		base.Nick(toString(row[cCreator])),
	}
	for _, val := range values {
		t, v := parseValue(toString(row[cKey]), toString(row[cRel]), val)
		if v == "" {
			// skip the many factoids with empty values.
			continue
		}
		out <- &factoids.Factoid{
			Key: toString(row[cKey]), Value: v, Type: t, Chance: 1.0,
			Created: c, Modified: m, Accessed: c, Perms: p,
			Id: bson.NewObjectId(),
		}
	}
}

// Parse out multiple entry values for a factoid key.
// This involves copying strings around a fair bit :-/
// Also, pipe-separated with escaped \| but not escaped \\
// is REALLY FUCKING STUPID and occasionally bad to parse.
func parseMultipleValues(v string) []string {
	temp_vals := strings.Split(v, "|")
	vals := make([]string, 0, len(temp_vals))
	for i := 0; i < len(temp_vals); i++ {
		str := temp_vals[i]
		for strings.HasSuffix(str, "\\") {
			// This | separator was escaped!
			i++
			if i < len(temp_vals) {
				str = strings.Join([]string{str, temp_vals[i]}, "|")
			} else {
				break
			}
		}
		vals = append(vals, str)
	}
	return vals
}

// Parse a single factoid value, stripping <me>/<reply>
func parseValue(k, r, v string) (ft factoids.FactoidType, fv string) {
	v = strings.TrimSpace(v)
	ft, fv = factoids.ParseValue(v)
	if ft == factoids.F_FACT && fv == v {
		// If fv == v, ParseValue hasn't stripped off a <reply>, so this is
		// just a normal factoid whose value is actually "key relation value"
		// as that's how they're stored in the old SQLite database...
		fv = strings.Join([]string{k, r, v}, " ")
	}
	return
}

// Parse the Created field with a type switch, cos it varies :-/
func parseTimestamp(ts interface{}) (time.Time, bool) {
	switch ts.(type) {
	case float64:
		return time.Unix(int64(ts.(float64)), 0), true
	case int64:
		return time.Unix(ts.(int64), 0), true
	case string:
		if tm, err := strconv.ParseInt(ts.(string), 10, 64); err == nil {
			return time.Unix(tm, 0), true
		} else {
			logging.Warn("parseTimestamp: %s", err)
		}
	}
	return time.Now(), false
}

// Ditto for the Access field.
func parseReadOnly(b interface{}) bool {
	switch b.(type) {
	case float64:
		return b.(float64) > 0
	case int64:
		return b.(int64) > 0
	case string:
		i, _ := strconv.Atoi(b.(string))
		return i > 0
	}
	// default to ReadOnly == false
	return false
}

// And in many other fields that *really* should be strings.
func toString(s interface{}) string {
	switch s.(type) {
	case float64:
		if float64(int(s.(float64))) == s.(float64) {
			return strconv.Itoa(int(s.(float64)))
		} else {
			return strconv.FormatFloat(s.(float64), 'f', -1, 64)
		}
	case int64:
		return strconv.FormatInt(s.(int64), 10)
	case string:
		return s.(string)
	}
	return ""
}

func main() {
	flag.Parse()
	logging.InitFromFlags()

	// Let's go find some mongo.
	db.Init()
	defer db.Close()
	fc := factoids.Init()

	// A communication channel of Factoids.
	facts := make(chan *factoids.Factoid)
	rows := make(chan []interface{})

	// Function to feed rows into the rows channel.
	row_feeder := func(sth *sqlite3.Statement, row ...interface{}) {
		rows <- row
	}

	// Function to execute a query on the SQLite db.
	db_query := func(dbh *sqlite3.Database) {
		n, err := dbh.Execute("SELECT * FROM Factoids;", row_feeder)
		if err == nil {
			logging.Info("Read %d rows from database.\n", n)
		} else {
			logging.Error("DB error: %s\n", err)
		}
	}

	// Open up the factoid database in a goroutine and feed rows
	// in on the input_rows channel.
	go func() {
		sqlite3.Session(*file, db_query)
		// once we've done the query, close the channel to indicate this
		close(rows)
	}()

	// Another goroutine to munge the rows into factoids.
	// This was originally done inside the SQLite callbacks, but
	// cgo or sqlite3 obscures runtime panics and makes fail happen.
	go func() {
		for row := range rows {
			parseFactoid(row, facts)
		}
		close(facts)
	}()

	// And finally...
	count := 0
	var err error
	for fact := range facts {
		// ... push each fact into mongo
		err = fc.Insert(fact)
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
	logging.Info("Inserted %d factoids.\n", count)
}
