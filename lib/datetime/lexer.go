package datetime

import (
	"fmt"
	"github.com/fluffle/sp0rkle/lib/util"
	"strconv"
	"strings"
	"time"
	"unicode"
)

// Indexes for relTime
type offset int
const (
	O_SEC offset = iota
	O_MIN
	O_HOUR
	O_DAY
	O_MONTH
	O_YEAR
)
var offsets = [...]string{
	"seconds",
	"minutes",
	"hours",
	"days",
	"months",
	"years",
}
func (r offset) String() string {
	return offsets[r]
}

type relTime [6]int
func (rt relTime) String() string {
	s := make([]string, 0, len(rt))
	for off, val := range rt {
		if val != 0 {
			s = append(s, fmt.Sprintf("%d %s", val, offsets[off]))
		}
	}
	return strings.Join(s, " ")
}

type relDays struct {
	day time.Weekday
	num int
	year int
}
func (rd relDays) String() string {
	s := fmt.Sprintf("%d %s", rd.num, rd.day)
	if rd.year != 0 {
		s += fmt.Sprintf(" of %d", rd.year)
	}
	return s
}

type relMonths struct {
	month time.Month
	num int
	year int
}
func (rm relMonths) String() string {
    if rm.month == 0 {
        return ""
    }
	s := fmt.Sprintf("%d %s", rm.num, rm.month)
	if rm.year != 0 {
		s += fmt.Sprintf(" of %d", rm.year)
	}
	return s
}

type lexerState int
const (
    HAVE_TIME lexerState = 1 << iota
    HAVE_DATE
    HAVE_DAYS
    HAVE_MONTHS
    HAVE_OFFSET
    HAVE_AGO
)
var lexerStates = [...]string{
    "time", "date", "days", "months", "offset", "ago"}
func (ls lexerState) String() string {
    s := make([]string, 0, len(lexerStates))
    for i := range lexerStates {
        if (ls & lexerState(1<<uint32(i))) != 0 {
            s = append(s, lexerStates[i])
        }
    }
    return strings.Join(s, " ")
}

type dateLexer struct {
	*util.Lexer
	hourfmt, ampmfmt, zonefmt string
	time, date time.Time  // takes care of absolute time and date specs
	offsets relTime       // takes care of +- ymd hms
    days    relDays       // takes care of specific days into future
	months  relMonths     // takes care of specific months into future
    states  lexerState
}

func (l *dateLexer) state(s lexerState, v ...bool) bool {
    ret := (l.states & s) != 0
    if len(v) > 0 {
        if v[0] {
            l.states |= s
        } else {
            l.states &= ^s
        }
    }
    return ret
}

func (l *dateLexer) Lex(lval *yySymType) int {
	l.Scan(unicode.IsSpace)
	c := l.Peek()
	
	switch {
	case c == '+':
		lval.intval = 1
		l.Next()
		return T_PLUS
	case c == '-':
		lval.intval = -1
		l.Next()
		return T_MINUS
	case unicode.IsDigit(c):
        s := l.Scan(unicode.IsDigit)
        i, _ := strconv.Atoi(s)
        lval.tval = textint{i, len(s), s}
		return T_INTEGER
	case unicode.IsLetter(c):
		input := strings.ToUpper(l.Scan(unicode.IsLetter))
        if tok, ok := tokenMaps.Lookup(input, lval); ok {
            return tok
        }
        // No token recognised, rewind and try the current character instead
        // as long as the original input was longer than that one character
		l.Rewind()
        if len(input) > 1 {
            input = strings.ToUpper(l.Next())
            if tok, ok := tokenMaps.Lookup(input, lval); ok {
                return tok
            }
            // Still not recognised.
            l.Rewind()
        }
	}
	l.Next()
    // At no time do we want to be case-sensitive
	return int(unicode.ToUpper(c))
}

func (l *dateLexer) Error(e string) {
	fmt.Println(e)
}

func (l *dateLexer) setTime(h, m, s int, loc *time.Location) {
	if loc == nil {
		loc = time.Local
	}
	fmt.Printf("Setting time to %d:%d:%d (%s)\n", h, m, s, loc)
	if l.state(HAVE_TIME, true) {
		l.Error("Parsed two times")
		return
	}
	l.time = time.Date(1, 1, 1, h, m, s, 0, loc)
}

func (l *dateLexer) setHMS(hms int, ln int, loc *time.Location) {
    hour, min, sec := 0, 0, 0
    if ln == 2 {
        // HH
        hour = hms
    } else if ln == 4 {
        // HHMM
        hour, min = hms / 100, hms % 100
    } else {
        // HHMMSS
        hour, min, sec = hms / 10000, (hms / 100) % 100, hms % 100
    }
    l.setTime(hour, min, sec, loc)
}

func (l *dateLexer) setDate(y, m, d int) {
	fmt.Printf("Setting date to %d-%d-%d\n", y, m, d)
	if l.state(HAVE_DATE, true) {
		l.Error("Parsed two dates")
        return
	}
	l.date = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}

func (l *dateLexer) setDay(d, n int, year ...int) {
	fmt.Printf("Setting day to %d %s\n", n, time.Weekday(d))
	if l.state(HAVE_DAYS, true) {
		l.Error("Parsed two days")
        return
	}
	l.days = relDays{time.Weekday(d), n, 0}
	if len(year) > 0 {
		l.days.year = year[0]
	}
}

func (l *dateLexer) setWeek(year, week, wday int) {
    // Week and wday are ISO numbers: week == 1-53, wday == 1-7, Monday == 1
    // http://en.wikipedia.org/wiki/ISO_week_date#Calculating_a_date_given_the_year.2C_week_number_and_weekday
    jan4 := int(time.Date(year, 1, 4, 0, 0, 0, 0, time.UTC).Weekday())
    if jan4 == 0 {
        // Go weekdays are 0-6, with Sunday == 0
        jan4 = 7
    }
    ord := week * 7 + wday - jan4 - 3
    l.setDate(year, 1, ord)
}

func (l *dateLexer) setMonth(m, n int, year ...int) {
	fmt.Printf("Setting month to %d %s\n", n, time.Month(m))
	if l.state(HAVE_MONTHS, true) {
		l.Error("Parsed two months")
        return
	}
	l.months = relMonths{time.Month(m), n, 0}
	if len(year) > 0 {
		l.months.year = year[0]
	}
}

func (l *dateLexer) setYMD(ymd int, ln int) {
    year, month, day := ymd / 10000, (ymd / 100) % 100, ymd % 100
    if ln == 6 {
        // YYMMDD not YYYYMMDD
        if year > 68 {
            year += 1900
        } else {
            year += 2000
        }
    }
    l.setDate(year, month, day)
}

func (l *dateLexer) addOffset(off offset, rel int) {
	fmt.Printf("Adding relative offset of %d %s\n", rel, off)
	l.offsets[off] += rel
    l.state(HAVE_OFFSET, true)
}

func (l *dateLexer) setAgo() {
    if l.state(HAVE_AGO, true) {
        l.Error("Parsed two agos")
        return
    }
	for i := range l.offsets {
		l.offsets[i] *= -1
	}
}

// Replaces rel's hour, minute, second and time.Location with the lexer's time
func (l *dateLexer) resolveTime(rel time.Time) time.Time {
    y, m, d := rel.Date()
    h, n, s := l.time.Clock()
    rel = time.Date(y, m, d, h, n, s, 0, l.time.Location())
    fmt.Printf("Parsed time as %s %s\n", rel.Weekday(), rel)
    // check if >24h has been given. Results of this may be *very* sketchy.
    // We can:
    //   a) drop >24h info completely
    //   b) save the integer number of hours as "days" and add that
    // Currently, do (a), but (b) would be nice.
    return rel
}

// Replaces rel's year, month and day with the lexer's date
func (l *dateLexer) resolveDate(rel time.Time) time.Time {
    y, m, d := l.date.Date()
    if y == 0 {
        y = rel.Year()
    }
    h, n, s := rel.Clock()
    rel = time.Date(y, m, d, h, n, s, 0, rel.Location())
    fmt.Printf("Parsed date as %s %s\n", rel.Weekday(), rel)
    return rel
}

// Applies a relative number of weeks to either the provided base time or
// Jan 1 of an absolute year to reach +-N <day>s or the nth <day> of <year>
func (l *dateLexer) resolveDays(rel time.Time) time.Time {
    if l.days.year != 0 {
        // this is num'th weekday of year, so start by finding jan 1
        h, n, s := rel.Clock()
        rel = time.Date(l.days.year, 1, 1, h, n, s, 0, rel.Location())
    }
    // Correct for the assumption that "<day>" or "this <day>" generally
    // refers to the coming <day> (unless it refers to *today*), while
    // "<month>" or "this <month>" refers to <month> of the current year.
    diff := int(l.days.day - rel.Weekday())
    if (diff != 0 && l.days.num == 0) || (diff < 0 && l.days.num < 0) {
        l.days.num++
    } else if (diff == 0 && l.days.year != 0) || (diff > 0 && l.days.num > 0) {
        l.days.num--
    }
    rel = rel.AddDate(0, 0, l.days.num * 7 + diff)
    fmt.Printf("Parsed days as %s %s\n", rel.Weekday(), rel)
    return rel
}

// Applies a relative number of months to the provided base time, or
// sets the month of an absolute year.
func (l *dateLexer) resolveMonths(rel time.Time) time.Time {
    h, n, s := rel.Clock()
    if l.months.year != 0 {
        // this is the month of an absolute year
        rel = time.Date(l.months.year, l.months.month, rel.Day(),
            h, n, s, 0, rel.Location())
        return rel
    }
    // this is relative months
    diff := int(l.months.month - rel.Month())
    rel = rel.AddDate(0, l.months.num*12+diff, 0)
    fmt.Printf("Parsed months as %s %s\n", rel.Weekday(), rel)
    return rel
}

// Applies the lexer's relative offset to the provided base time.
func (l *dateLexer) resolveOffset(rel time.Time) time.Time {
    rel = rel.AddDate(l.offsets[O_YEAR], l.offsets[O_MONTH], l.offsets[O_DAY])
    rel = rel.Add(time.Hour * time.Duration(l.offsets[O_HOUR]) +
        time.Minute * time.Duration(l.offsets[O_MIN]) +
        time.Second * time.Duration(l.offsets[O_SEC]))
    fmt.Printf("Parsed offset as %s %s\n", rel.Weekday(), rel)
    return rel
}

func lexAndParse(input string) (*dateLexer, int) {
	lexer := &dateLexer{Lexer: &util.Lexer{Input: input}}
	yyDebug = 5
	if ret := yyParse(lexer); ret != 0 {
        return nil, ret
	}
    fmt.Println("state: ", lexer.states)
    fmt.Println("time: ", lexer.time)
    fmt.Println("date: ", lexer.date)
    fmt.Println("days: ", lexer.days)
    fmt.Println("months: ", lexer.months)
    fmt.Println("offset: ", lexer.offsets)
    return lexer, 0
}

func resolve(l *dateLexer, rel time.Time) (time.Time, bool) {
    if l.state(HAVE_DATE) && (l.days.year != 0 || l.months.year != 0) {
        // DATE is absolute, another absolute DAYS or MONTHS is ambiguous
        return time.Time{}, false
    }

    // Resolve any absolute time and date first
    if l.state(HAVE_TIME) {
        rel = l.resolveTime(rel)
    }
    if l.state(HAVE_DATE) {
        rel = l.resolveDate(rel)
    }
    // We need to resolve relative or absolute months before relative days
    if l.state(HAVE_MONTHS) {
        rel = l.resolveMonths(rel)
    }
    if l.state(HAVE_DAYS) {
        rel = l.resolveDays(rel)
    }
    // And then apply any offset
    if l.state(HAVE_OFFSET) {
        rel = l.resolveOffset(rel)
    }
    return rel, true
}

func Parse(input string) (time.Time, bool) {
    lexer, ret := lexAndParse(input)
    if lexer == nil || lexer.states == 0 {
        fmt.Println("Parse error: ", ret)
    	return time.Time{}, false
    }
    // return time.Time{}, false
    return resolve(lexer, time.Now())
}
