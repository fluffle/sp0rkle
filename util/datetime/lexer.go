//go:generate go tool yacc datetime.y
package datetime

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/fluffle/sp0rkle/util"
)

func DPrintf(f string, args ...interface{}) {
	if yyDebug > 0 {
		fmt.Printf(f, args...)
	}
}

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
	day  time.Weekday
	num  int
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
	num   int
	year  int
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

func ampm(hour, offset int) int {
	// Take care of the fact that 12am is midnight and 12pm midday
	switch hour + offset {
	case 12:
		return 0
	case 24:
		return 12
	default:
		return hour + offset
	}
}

type lexerState int

const (
	HAVE_TIME lexerState = 1 << iota
	HAVE_DATE
	HAVE_DAY
	HAVE_DAYS
	HAVE_DYEAR
	HAVE_MONTHS
	HAVE_MYEAR
	HAVE_OFFSET
	HAVE_AGO
	HAVE_ABSYEAR  = HAVE_DYEAR | HAVE_MYEAR
	HAVE_DMY      = HAVE_DAYS | HAVE_MONTHS | HAVE_ABSYEAR
	HAVE_DATETIME = HAVE_DATE | HAVE_TIME
)

var lexerStates = [...]string{
	"time", "date", "day", "days", "day-year", "months", "month-year", "offset", "ago"}

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
	rel                       time.Time // base time any relative offsets are computed against
	time, date                time.Time // takes care of absolute time and date specs
	day                       int       // takes care of absolute day of relative month
	offsets                   relTime   // takes care of +- ymd hms
	days                      relDays   // takes care of specific days into future
	months                    relMonths // takes care of specific months into future
	states                    lexerState
	errors                    []string
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
		pos := l.Pos()
		input := l.Scan(unicode.IsLetter)
		if tok, ok := tokenMaps.Lookup(strings.ToUpper(input), lval); ok {
			return tok
		}
		// No token recognised, but it could be a zone in IANA format!
		zstr := input + l.Not(unicode.IsSpace)
		if z := zone(zstr); z != nil {
			lval.zoneval = z
			return T_ZONE
		}
		l.Pos(pos)
		// No token recognised, rewind and try the current character instead
		// as long as the original input was longer than that one character
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
	l.errors = append(l.errors, e)
}

func (l *dateLexer) setUnix(epoch int64) {
	DPrintf("Setting unix timestamp to %d.\n", epoch)
	if l.state(HAVE_DATETIME, true) {
		l.Error("unix timestamp plus other time/date specifier")
		return
	}
	l.time = time.Unix(epoch, 0)
	l.date = l.time
}

func (l *dateLexer) setTime(h, m, s int, loc *time.Location) {
	if loc == nil {
		loc = l.rel.Location()
	}
	h, m, s = h%24, m%60, s%60
	DPrintf("Setting time to %d:%d:%d (%s)\n", h, m, s, loc)
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
		hour, min = hms/100, hms%100
	} else {
		// HHMMSS
		hour, min, sec = hms/10000, (hms/100)%100, hms%100
	}
	l.setTime(hour, min, sec, loc)
}

func (l *dateLexer) setDate(y, m, d int) {
	DPrintf("Setting date to %d-%d-%d\n", y, m, d)
	if l.state(HAVE_DATE, true) {
		l.Error("Parsed two dates")
		return
	}
	l.date = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}

func (l *dateLexer) setDay(d int) {
	DPrintf("Setting day to %d\n", d)
	if l.state(HAVE_DAY, true) {
		l.Error("Parsed two absolute days")
		return
	}
	l.day = d
}

func (l *dateLexer) setDays(d, n int, year ...int) {
	DPrintf("Setting days to %d %s\n", n, time.Weekday(d))
	if l.state(HAVE_DAYS, true) {
		l.Error("Parsed two days")
	}
	l.days = relDays{time.Weekday(d), n, 0}
	if len(year) > 0 {
		l.state(HAVE_DYEAR, true)
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
	ord := week*7 + wday - jan4 - 3
	DPrintf("Setting week to %d week %d day %d (ord=%d, jan4=%d)\n",
		year, week, wday, ord, jan4)
	l.setDate(year, 1, ord)
}

func (l *dateLexer) setMonths(m, n int, year ...int) {
	DPrintf("Setting month to %d %s\n", n, time.Month(m))
	if l.state(HAVE_MONTHS, true) {
		l.Error("Parsed two months")
		return
	}
	l.months = relMonths{time.Month(m), n, 0}
	if len(year) > 0 {
		l.state(HAVE_MYEAR, true)
		l.months.year = year[0]
	}
}

func (l *dateLexer) setYear(year int) {
	DPrintf("Setting year to %d\n", year)
	if l.state(HAVE_DATE) {
		l.date = time.Date(year, l.date.Month(), l.date.Day(),
			0, 0, 0, 0, time.Local)
	} else if l.state(HAVE_MONTHS) {
		l.state(HAVE_MYEAR, true)
		l.months.year = year
	} else {
		l.state(HAVE_DYEAR, true)
		l.days.year = year
	}
}

func (l *dateLexer) setYMD(ymd int, ln int) {
	year, month, day := ymd/10000, (ymd/100)%100, ymd%100
	if ln == 6 {
		// YYMMDD not YYYYMMDD
		if year > 68 {
			year += 1900
		} else {
			year += 2000
		}
	}
	DPrintf("Setting YYYY-MM-DD to %04d-%02d-%02d\n", year, month, day)
	l.setDate(year, month, day)
}

func (l *dateLexer) addOffset(off offset, rel int) {
	DPrintf("Adding relative offset of %d %s\n", rel, off)
	l.offsets[off] += rel
	l.state(HAVE_OFFSET, true)
}

func (l *dateLexer) setAgo() {
	DPrintf("Setting all offsets to negative because of 'ago'.\n")
	if l.state(HAVE_AGO, true) {
		l.Error("Parsed two agos")
		return
	}
	for i := range l.offsets {
		l.offsets[i] *= -1
	}
}

// Replaces rel's hour, minute, second and time.Location with the lexer's time
func (l *dateLexer) resolveTime() {
	y, m, d := l.rel.Date()
	h, n, s := l.time.Clock()
	// We can:
	//   a) drop >24h info completely
	//   b) save the integer number of hours as "days" and add that
	// Currently, do (a), but (b) would be nice.
	l.rel = time.Date(y, m, d, h, n, s, 0, l.time.Location())
}

// Replaces rel's year, month and day with the lexer's date
func (l *dateLexer) resolveDate() {
	y, m, d := l.date.Date()
	if y == 0 {
		y = l.rel.Year()
	}
	h, n, s := l.rel.Clock()
	l.rel = time.Date(y, m, d, h, n, s, 0, l.rel.Location())
}

func (l *dateLexer) resolveDay() {
	y, m, _ := l.rel.Date()
	h, n, s := l.rel.Clock()
	l.rel = time.Date(y, m, l.day, h, n, s, 0, l.rel.Location())
}

func (l *dateLexer) dayOffset() {
	// Correct for the assumption that "<day>" or "this <day>" generally
	// refers to the coming <day> unless it refers to the day of this year
	// or *today* whilst "next <day>" *always* refers to the coming <day>.
	diff := int(l.days.day - l.rel.Weekday())
	if diff < 0 && l.days.num <= 0 {
		DPrintf("Day offset %d->%d, diff=%d.\n", l.days.num, l.days.num+1, diff)
		l.days.num++
	} else if diff > 0 && l.days.num > 0 {
		DPrintf("Day offset %d->%d, diff=%d.\n", l.days.num, l.days.num-1, diff)
		l.days.num--
	}
	l.rel = l.rel.AddDate(0, 0, l.days.num*7+diff)
}

func (l *dateLexer) monthOffset() {
	diff := int(l.months.month - l.rel.Month())
	if l.months.num == 0 {
		// If just "march" or "this march" find closest month
		// preferring 6 months in future over 6 months in past
		diff = ((diff + 5) % 12) - 5
		DPrintf("Month offset %d months\n", diff)
		l.rel = l.rel.AddDate(0, diff, 0)
		return
	}
	if diff < 0 && l.months.num < 0 {
		DPrintf("Month offset %d->%d, diff=%d.\n", l.months.num, l.months.num+1, diff)
		l.months.num++
	} else if diff > 0 && l.months.num > 0 {
		DPrintf("Month offset %d->%d, diff=%d.\n", l.months.num, l.months.num-1, diff)
		l.months.num--
	}
	l.rel = l.rel.AddDate(0, l.months.num*12+diff, 0)
}

func (l *dateLexer) resolveDMY() {
	h, n, s := l.rel.Clock()
	mkrel := func(y int, m time.Month, d int) time.Time {
		return time.Date(y, m, d, h, n, s, 0, l.rel.Location())
	}
	switch l.states & HAVE_DMY {
	case HAVE_MYEAR + HAVE_MONTHS:
		// this is month year, so just set those and bail out
		DPrintf("MYEAR & MONTHS\n")
		l.rel = mkrel(l.months.year, l.months.month, l.rel.Day())
	case HAVE_DYEAR + HAVE_DAYS:
		// this is num'th weekday of year, so compute day offset from "jan 0"
		DPrintf("DYEAR & DAYS\n")
		l.rel = mkrel(l.days.year, 1, 0)
		l.dayOffset()
	case HAVE_MYEAR + HAVE_MONTHS + HAVE_DAYS:
		// this is num'th weekday of month year, so offset from "0th" of month in year
		DPrintf("MYEAR & MONTHS & DAYS\n")
		l.rel = mkrel(l.months.year, l.months.month, 0)
		l.dayOffset()
	case HAVE_DAYS:
		DPrintf("DAYS\n")
		l.dayOffset()
	case HAVE_MONTHS:
		DPrintf("MONTHS\n")
		l.monthOffset()
	case HAVE_MONTHS + HAVE_DAYS:
		DPrintf("MONTHS & DAYS\n")
		if l.months.month == 0 {
			// just num'th weekday (of this month, implicitly)
			l.months.month = l.rel.Month()
		}
		l.monthOffset()
		// this is num'th weekday of month, so we need to offset from "0th"
		l.rel = mkrel(l.rel.Year(), l.rel.Month(), 0)
		l.dayOffset()
	case HAVE_MYEAR:
		DPrintf("MYEAR\n")
		// These on their own are a little odd but probably due to the hack at
		// datetime.y:163 so just set the explicit year and return
		l.rel = mkrel(l.months.year, l.rel.Month(), l.rel.Day())
	case HAVE_DYEAR:
		DPrintf("DYEAR\n")
		l.rel = mkrel(l.days.year, l.rel.Month(), l.rel.Day())
	default:
		panic("oh fuck this is too complicated :-(\n" + l.states.String())
	}
}

// Applies the lexer's relative offset to the provided base time.
func (l *dateLexer) resolveOffset() {
	l.rel = l.rel.AddDate(l.offsets[O_YEAR], l.offsets[O_MONTH], l.offsets[O_DAY])
	l.rel = l.rel.Add(time.Hour*time.Duration(l.offsets[O_HOUR]) +
		time.Minute*time.Duration(l.offsets[O_MIN]) +
		time.Second*time.Duration(l.offsets[O_SEC]))
}

func (l *dateLexer) resolve() (time.Time, error) {
	DPrintf("Lexer state: %s\n", l.states)
	if (l.state(HAVE_DATE) && l.state(HAVE_ABSYEAR)) ||
		(l.state(HAVE_DYEAR) && l.state(HAVE_MYEAR)) {
		// DATE is absolute, another absolute DAYS or MONTHS is ambiguous
		// Providing an absolute monthyear and an absolute dayyear is also bad
		return time.Time{}, errors.New("multiple conflicting absolute dates")
	}

	// Resolve any absolute time and date first
	if l.state(HAVE_TIME) {
		DPrintf("HAVE_TIME before: %s %s\n", l.rel.Weekday(), l.rel)
		DPrintf("Lexer time: %s\n", l.time)
		l.resolveTime()
		DPrintf("HAVE_TIME after: %s %s\n", l.rel.Weekday(), l.rel)
	}
	if l.state(HAVE_DATE) {
		DPrintf("HAVE_DATE before: %s %s\n", l.rel.Weekday(), l.rel)
		DPrintf("Lexer date: %s\n", l.date)
		l.resolveDate()
		DPrintf("HAVE_DATE after: %s %s\n", l.rel.Weekday(), l.rel)
	}
	if l.state(HAVE_DAY) {
		DPrintf("HAVE_DAY before: %s %s\n", l.rel.Weekday(), l.rel)
		DPrintf("Lexer day: %d\n", l.day)
		l.resolveDay()
		DPrintf("HAVE_DAY after: %s %s\n", l.rel.Weekday(), l.rel)
	}
	// Apply any offsets (so that l.relative days are from the offset time)
	if l.state(HAVE_OFFSET) {
		DPrintf("HAVE_OFFSET before: %s %s\n", l.rel.Weekday(), l.rel)
		DPrintf("Lexer offsets: %s\n", l.offsets)
		l.resolveOffset()
		DPrintf("HAVE_OFFSET after: %s %s\n", l.rel.Weekday(), l.rel)
	}
	// Resolve l.relative/absolute day/month/years
	if l.state(HAVE_DMY) {
		DPrintf("HAVE_DMY before: %s %s\n", l.rel.Weekday(), l.rel)
		DPrintf("Lexer months: %s\n", l.months)
		DPrintf("Lexer days: %s\n", l.days)
		l.resolveDMY()
		DPrintf("HAVE_DMY after: %s %s\n", l.rel.Weekday(), l.rel)
	}
	return l.rel, nil
}

func parse(input string, rel time.Time) (time.Time, error) {
	yyDebug = 0
	yyErrorVerbose = true

	lexer := &dateLexer{Lexer: &util.Lexer{Input: input}, rel: rel}
	if ret := yyParse(lexer); ret != 0 {
		return time.Time{}, errors.New(strings.Join(lexer.errors, "; "))
	}
	if lexer.states == 0 {
		return time.Time{}, errors.New("no dates parsed from input")
	}
	return lexer.resolve()
}

func Parse(input string) (time.Time, error) {
	return parse(input, time.Now().In(local))
}

func ParseZ(input string, zone *time.Location) (time.Time, error) {
	return parse(input, time.Now().In(zone))
}
