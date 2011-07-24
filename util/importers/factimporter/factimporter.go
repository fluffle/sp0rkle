package main

// Imports perlfu's SQLite factoid database into mongodb using lib/factoids

import (
	"flag"
	"fmt"
	"github.com/garyburd/go-mongo"
	"github.com/kuroneko/gosqlite3"
	"lib/factoids"
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
		Timestamp: parseTimestamp(row[cCreated]),
		Nick: toString(row[cCreator]),
		// We don't know these things :-(
		Ident: "", Host: "", Chan: "", Count: 1,
	}
	m := &factoids.FactoidStat{
		Ident: "", Host: "", Chan: "", Count: 0,
	}
	if ts := parseTimestamp(row[cModified]); ts != nil {
		m.Timestamp = ts
		m.Nick = toString(row[cModifier])
	} else {
		m.Timestamp = c.Timestamp
		m.Nick = c.Nick
	}
	p := &factoids.FactoidPerms{
		ReadOnly: parseReadOnly(row[cAccess]),
		Owner: toString(row[cCreator]),
	}
	for _, val := range values {
		t, v := parseValue(toString(row[cKey]), toString(row[cRel]), val)
		out <- &factoids.Factoid{
			Key: toString(row[cKey]), Value: v, Type: t,
			Created: c, Modified: m, Accessed: nil, Perms: p,
		}
	}
}

// Parse out multiple entry values for a factoid key.
// This involves copying strings around a fair bit :-/
// Also, pipe-separated with escaped \| but not escaped \\
// is REALLY FUCKING STUPID and occasionally bad to parse.
func parseMultipleValues(v string) []string {
	temp_vals := strings.Split(v, "|", -1)
	vals := make([]string, 0, len(temp_vals))
	for i:= 0; i < len(temp_vals); i++ {
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
	if strings.HasPrefix(v, "<me>") {
		// <me>does something
		ft, fv = factoids.F_ACTION, v[4:]
	} else if strings.HasPrefix(v, "<reply>") {
		// <reply> 
		ft, fv = factoids.F_REPLY, v[7:]
	} else {
		fv = v
	}
	if looksURLish(fv) {
		// Quite a few factoids are just <reply>http://some.url/
		// it's helpful to detect this so we can do useful things
		ft = factoids.F_URL
	} else {
		// Just a normal factoid whose value is actually "key relation value"
		ft, fv = factoids.F_FACT, strings.Join([]string{k,r,v}, " ")
	}
	return
}

// Does this string look like a URL to you?
// This should be fairly conservative, I hope:
//   s starts with http:// or https:// and contains no spaces
func looksURLish(s string) bool {
	return ((strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://")) &&
		strings.Index(s, " ") == -1)
}

// Parse the Created field with a type switch, cos it varies :-/
func parseTimestamp(ts interface{}) *time.Time {
	var tm int64
	switch ts.(type) {
		case float64:
			tm = int64(ts.(float64))
		case int64:
			tm = ts.(int64)
		case string:
			tm, _ = strconv.Atoi64(ts.(string))
		default:
			return nil
	}
	return time.SecondsToLocalTime(tm)
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
				return strconv.Ftoa64(s.(float64), 'f', -1)
			}
		case int64:
			return strconv.Itoa64(s.(int64))
		case string:
			return s.(string)
	}
	return ""
}

func main() {
	// Let's go find some mongo.
	conn, err := mongo.Dial("localhost")
	if err != nil {
		fmt.Printf("Oh no: %v", err)
		return
	}
	defer conn.Close()
	fc, err := factoids.Collection(conn)
	if err != nil {
		fmt.Printf("Oh no: %v", err)
		return
	}

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
			fmt.Printf("Read %d rows from database.\n", n)
		} else {
			fmt.Printf("DB error: %s\n", err)
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
	for fact := range facts {
		// ... push each fact into mongo
		err = fc.Insert(fact)
		count++
		if err != nil {
			fmt.Printf("Awww: %v", err)
			continue
		}
	}
	fmt.Printf("Inserted %d factoids.\n", count)
}
