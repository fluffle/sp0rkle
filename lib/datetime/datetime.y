%{
package datetime

// Based upon parse-datetime.y in GNU coreutils.
// also an exercise in learning goyacc in particular.

import (
	"fmt"
	"github.com/fluffle/sp0rkle/lib/util"
//	"math"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type textint struct {
    i, l int
    s string
}

%}

%union
{
	tval    textint
	intval  int
	zoneval *time.Location
}

%token           T_DAYQUAL
%token <tval>    T_INTEGER
%token <intval>  T_PLUS T_MINUS
%token <intval>  T_AMPM T_MONTHNAME T_DAYNAME T_DAYS T_DAYSHIFT
%token <intval>  T_OFFSET T_ISOYD T_ISOHS T_RELATIVE T_AGO
%token <zoneval> T_ZONE

%type <intval> sign
%type <tval> o_sign_integer
%type <zoneval> zone o_zone

%%

spec:
	unixtime | items;

o_comma:
	/* empty */ | ',';

of: 'O' 'F';

o_of: /* empty */ | of;

o_dayqual: /* empty */ | T_DAYQUAL;

sign:
	T_PLUS | T_MINUS; 

o_sign_integer:
	T_INTEGER
	| T_PLUS T_INTEGER {
        $2.s = "+" + $2.s
        $$ = $2
    }
	| T_MINUS T_INTEGER {
        $2.s = "-" + $2.s
        $2.i *= -1
        $$ = $2
    };

zone:
	T_ZONE
	| sign T_INTEGER {
        hrs, mins := $2.i, 0
        if ($2.l == 4) {
            hrs, mins = ($2.i / 100), ($2.i % 100)
        } else if ($2.l == 2) {
            hrs *= 100
        } else {
            yylex.Error("Invalid timezone offset " +$2.s)
        }
		$$ = time.FixedZone("WTF", $1 * (3600 * hrs + 60 * mins))
	}
	| sign T_INTEGER ':' T_INTEGER {
		$$ = time.FixedZone("WTF", $1 * (3600 * $2.i + 60 * $4.i))
	};

o_zone:
	/* empty */ { $$ = nil }
    | zone;
    
unixtime:
	'@' o_sign_integer {
		yylex.(*dateLexer).time = time.Unix(int64($2.i), 0)
	};

items:
	/* empty */
	| items item;

item:
	time
    | iso_8601_time
	| date
    | iso_8601_date
    | iso_8601_date_time
	| day_or_month
	| relative
    | iso_8601_duration
    | integer;

// ISO 8601 takes care of 24h time formats, so this deals with
// 12-hour HH, HH:MM or HH:MM:SS with am/pm and optional timezone
time:
	T_INTEGER T_AMPM o_zone {
		yylex.(*dateLexer).setTime($1.i + $2, 0, 0, $3)
	}
	| T_INTEGER ':' T_INTEGER T_AMPM o_zone {
		yylex.(*dateLexer).setTime($1.i + $4, $3.i, 0, $5)
	}
	| T_INTEGER ':' T_INTEGER ':' T_INTEGER T_AMPM o_zone {
		yylex.(*dateLexer).setTime($1.i + $6, $3.i, $5.i, $7)
	};

// The "basic" ISO 8601 format (without a timezone) is lexed as
// an integer and handled in 'integer' below
iso_8601_time:
    T_INTEGER zone {
        yylex.(*dateLexer).setHMS($1.i, $1.l, $2)
    }
    | T_INTEGER ':' T_INTEGER o_zone {
        yylex.(*dateLexer).setTime($1.i, $3.i, 0, $4)
    }
    | T_INTEGER ':' T_INTEGER ':' T_INTEGER o_zone {
        yylex.(*dateLexer).setTime($1.i, $3.i, $5.i, $6)
    };

// ISO 8601 takes care of dash-separated big-endian date formats,
// so this deals with /-separated little-endian formats (dd/mm/yyyy)
// and more "english" ones like "20th of March 2012"
date:
	T_INTEGER '/' T_INTEGER {
		l := yylex.(*dateLexer)
		if $3.l == 4 {
			// assume we have MM/YYYY
			l.setDate($3.i, $1.i, 1)
		} else {
            // assume we have DD/MM (too bad, americans)
            l.setDate(0, $3.i, $1.i)
		}
	}
	| T_INTEGER '/' T_INTEGER '/' T_INTEGER {
		l := yylex.(*dateLexer)
		if $5.l == 4 {
			// assume we have DD/MM/YYYY
			l.setDate($5.i, $3.i, $1.i)
		} else if $5.i > 68 {
			// assume we have DD/MM/YY, add 1900 if YY > 68
			l.setDate($5.i + 1900, $3.i, $1.i)
		} else {
			// assume we have DD/MM/YY, add 2000 otherwise
			l.setDate($5.i + 2000, $3.i, $1.i)
		}
	}
	| T_INTEGER o_dayqual o_of T_MONTHNAME {
		// DDth of Mon
		yylex.(*dateLexer).setDate(0, $4, $1.i)
	}
	| T_MONTHNAME T_INTEGER o_dayqual {
		l := yylex.(*dateLexer)
		if $2.l == 4 {
			// assume Mon YYYY
			l.setDate($2.i, $1, 1)
		} else {
		    // assume Mon DDth
			l.setDate(0, $1, $2.i)
		}
    }
	| T_INTEGER o_dayqual o_of T_MONTHNAME T_INTEGER {
		l := yylex.(*dateLexer)
		if $5.l == 4 {
			// assume DDth of Mon YYYY
			l.setDate($5.i, $4, $1.i)
		} else if $5.i > 68 {
			// assume DDth of Mon YY, add 1900 if YY > 68
			l.setDate($5.i + 1900, $4, $1.i)
		} else {
			// assume DDth of Mon YY, add 2000 otherwise
			l.setDate($5.i + 2000, $4, $1.i)
		}
	}
	| T_MONTHNAME T_INTEGER o_dayqual o_comma T_INTEGER {
		l := yylex.(*dateLexer)
		if $5.l == 4 {
			// assume Mon DDth, YYYY
			l.setDate($5.i, $1, $2.i)
		} else if $5.i > 68 {
			// assume Mon DDth, YY, add 1900 if YY > 68
			l.setDate($5.i + 1900, $1, $2.i)
		} else {
			// assume Mon DDth YY, add 2000 otherwise
			l.setDate($5.i + 2000, $1, $2.i)
		}
	};

// The "basic" ISO 8601 format is lexed as an integer and handled in "integer"
iso_8601_date:
	T_INTEGER T_MINUS T_INTEGER {
		l := yylex.(*dateLexer)
		if $1.l == 4 && $3.l == 3 {
            // assume we have YYYY-DDD
            l.setDate($1.i, 1, $3.i)
        } else if $1.l == 4 {
			// assume we have YYYY-MM
			l.setDate($1.i, $3.i, 1)
		} else {
            // assume we have MM-DD (not strictly ISO compliant)
            // this is for americans, because of DD/MM above ;-)
            l.setDate(0, $3.i, $1.i)
		}
	}
	| T_INTEGER T_MINUS T_INTEGER T_MINUS T_INTEGER {
		l := yylex.(*dateLexer)
		if $1.l == 4 {
			// assume we have YYYY-MM-DD
			l.setDate($1.i, $3.i, $5.i)
		} else if $1.i > 68 {
			// assume we have YY-MM-DD, add 1900 if YY > 68
			l.setDate($1.i + 1900, $3.i, $5.i)
		} else {
			// assume we have YY-MM-DD, add 2000 otherwise
			l.setDate($1.i + 2000, $3.i, $5.i)
		}
	}
    | T_INTEGER 'W' T_INTEGER {
        l := yylex.(*dateLexer)
        wday, week := 1, $3.i
        if $3.l == 3 {
            // assume YYYY'W'WWD
            week = week / 10
            wday = week % 10
        }
        l.setWeek($1.i, week, wday)
    }
    | T_INTEGER T_MINUS 'W' T_INTEGER {
        // assume YYYY-'W'WW
        yylex.(*dateLexer).setWeek($1.i, $4.i, 1)
    }
    | T_INTEGER T_MINUS 'W' T_INTEGER T_MINUS T_INTEGER {
        // assume YYYY-'W'WW-D
        yylex.(*dateLexer).setWeek($1.i, $4.i, $6.i)
    };

// NOTE: this doesn't enforce that the date is complete.
iso_8601_date_time:
    iso_8601_date 'T' iso_8601_time
    | T_INTEGER 'T' T_INTEGER o_zone {
        // this goes here because the YYYYMMDD and HHMMSS forms of the
        // ISO 8601 format date and time are handled by 'integer' below.
        l := yylex.(*dateLexer)
        l.setYMD($1.i, $1.l)
        l.setHMS($3.i, $3.l, $4)
    };

day_or_month:
	T_DAYNAME o_comma {
		// Tuesday,
		yylex.(*dateLexer).setDay($1, 1)
	}
	| T_MONTHNAME {
		// March
		yylex.(*dateLexer).setMonth($1, 1)
	}
	| T_RELATIVE T_DAYNAME {
		// Next tuesday
		yylex.(*dateLexer).setDay($2, $1)
	}
	| T_RELATIVE T_MONTHNAME {
		// Next march
		yylex.(*dateLexer).setMonth($2, $1)
	}
	| o_sign_integer T_DAYNAME {
		// +-N Tuesdays
		yylex.(*dateLexer).setDay($2, $1.i)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME {
		// 3rd Tuesday 
		yylex.(*dateLexer).setDay($3, $1.i)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME of T_MONTHNAME {
		// 3rd Tuesday of (implicit this) March
		l := yylex.(*dateLexer)
		l.setDay($3, $1.i)
		l.setMonth($5, 1)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME of T_INTEGER {
		// 3rd Tuesday of 2012
		yylex.(*dateLexer).setDay($3, $1.i, $5.i)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME of T_MONTHNAME T_INTEGER {
		// 3rd Tuesday of March 2012
		l := yylex.(*dateLexer)
		l.setDay($3, $1.i)
		l.setMonth($5, 1, $6.i)
	}
	| T_INTEGER T_DAYQUAL T_DAYNAME of T_RELATIVE T_MONTHNAME {
		// 3rd Tuesday of next March
		l := yylex.(*dateLexer)
		l.setDay($3, $1.i)
		l.setMonth($6, $5)
	}
	| T_DAYSHIFT {
		// yesterday or tomorrow
		d := time.Now().Weekday()
		yylex.(*dateLexer).setDay((7+int(d)+$1)%7, $1)
	};

relative:
	relunits 
	| relunits T_AGO {
		yylex.(*dateLexer).setAgo()
	};

relunits:
	relunit
	| relunit relunits;

relunit:
	o_sign_integer T_OFFSET {
		yylex.(*dateLexer).addOffset(offset($2), $1.i)
	}
	| T_RELATIVE T_OFFSET {
		yylex.(*dateLexer).addOffset(offset($2), $1)
	} 
	| o_sign_integer T_DAYS {
		// Special-case to handle "week" and "fortnight"
		yylex.(*dateLexer).addOffset(O_DAY, $1.i * $2)
	}
	| T_RELATIVE T_DAYS {
		yylex.(*dateLexer).addOffset(O_DAY, $1 * $2)
	}
    | o_sign_integer T_ISOYD {
        // As we need to be able to separate out YD from HS in ISO durations
        // this becomes a fair bit messier than if Y D H S were just T_OFFSET
        // Because writing "next y" or "two h" would be odd, disallow
        // T_RELATIVE tokens from being used with ISO single-letter notation
        yylex.(*dateLexer).addOffset(offset($2), $1.i)
    }
    | o_sign_integer T_ISOHS {
        yylex.(*dateLexer).addOffset(offset($2), $1.i)
    }
    | o_sign_integer 'M' {
        // Resolve 'm' ambiguity in favour of minutes outside ISO duration
        yylex.(*dateLexer).addOffset(O_MIN, $1.i)
    };

/* date/time based durations not yet supported */
iso_8601_duration:
    'P' ymd_units o_t_hms_units
    | 'P' t_hms_units
    | 'P' T_INTEGER 'W' {
        yylex.(*dateLexer).addOffset(O_DAY, 7 * $2.i)
    };

/* This is a bit lazy compared to specifying the combinations of nYnMnS */
ymd_units:
    ymd_unit
    | ymd_units ymd_unit;

ymd_unit:
    T_INTEGER T_ISOYD {
        // takes care of Y and D
        yylex.(*dateLexer).addOffset(offset($2), $1.i)
    }
    | T_INTEGER 'M' {
        yylex.(*dateLexer).addOffset(O_MONTH, $1.i)
    };

hms_units:
    hms_unit
    | hms_units hms_unit;

hms_unit:
    T_INTEGER T_ISOHS {
        // takes care of H and S
        yylex.(*dateLexer).addOffset(offset($2), $1.i)
    }
    | T_INTEGER 'M' {
        yylex.(*dateLexer).addOffset(O_MIN, $1.i)
    };

t_hms_units:
    'T' hms_units;

o_t_hms_units:
    /* empty */
    | t_hms_units;

integer:
	T_INTEGER {
        l := yylex.(*dateLexer)
        if $1.l == 8 {
            // assume ISO 8601 YYYYMMDD
            l.setYMD($1.i, $1.l)
        } else {
            // assume ISO 8601 HHMMSS with no zone
            l.setHMS($1.i, $1.l, nil)
        }   
    };

/*
hybrid:
	tUNUMBER relunit_snumber
	  {
		// Hybrid all-digit and relative offset, so that we accept e.g.,
		// "YYYYMMDD +N days" as well as "YYYYMMDD N days".
		digits_to_date_time (pc, $1);
		apply_relative_time (pc, $2, 1);
	  }
  ;

*/
%%

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
type relTime struct {
	offsets [6]int
	seen bool
}
func (rt relTime) String() string {
	if !rt.seen {
		return "No time offsets seen"
	}
	s := make([]string, 0, 6)
	for off, val := range rt.offsets {
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
	seen bool
}
func (rd relDays) String() string {
	if !rd.seen {
		return "No relative days seen"
	}
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
	seen bool
}
func (rm relMonths) String() string {
	if !rm.seen {
		return "No relative months seen"
	}
	s := fmt.Sprintf("%d %s", rm.num, rm.month)
	if rm.year != 0 {
		s += fmt.Sprintf(" of %d", rm.year)
	}
	return s
}

type dateLexer struct {
	*util.Lexer
	hourfmt, ampmfmt, zonefmt string
	time, date time.Time
	offsets relTime       // takes care of +- ymd hms
    days    relDays       // takes care of specific days into future
	months  relMonths     // takes care of specific months into future
    ago     bool          // more than one "ago" is probably bad
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
	if ! l.time.IsZero() {
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
	if ! l.date.IsZero() {
		l.Error("Parsed two dates")
        return
	}
	l.date = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}

func (l *dateLexer) setDay(d, n int, year ...int) {
	fmt.Printf("Setting day to %d %s\n", n, time.Weekday(d))
	if l.days.seen {
		l.Error("Parsed two days")
        return
	}
	l.days = relDays{time.Weekday(d), n, 0, true}
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
	if l.months.seen {
		l.Error("Parsed two months")
        return
	}
	l.months = relMonths{time.Month(m), n, 0, true}
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
	l.offsets.seen = true
	l.offsets.offsets[off] += rel
}

func (l *dateLexer) setAgo() {
    if l.ago {
        l.Error("Parsed two agos")
        return
    }
	for i := range l.offsets.offsets {
		l.offsets.offsets[i] *= -1
	}
    l.ago = true
}

func Parse(input string) (time.Time, bool) {
    lexer, ret := lexAndParse(input)
    if lexer == nil {
        fmt.Println("Parse error: ", ret)
    	return time.Time{}, false
    }
    // return time.Time{}, false
    return resolve(lexer, time.Now())
}

func lexAndParse(input string) (*dateLexer, int) {
	lexer := &dateLexer{Lexer: &util.Lexer{Input: input}}
	yyDebug = 5
	if ret := yyParse(lexer); ret != 0 {
        return nil, ret
	}
    fmt.Println(lexer.time)
    fmt.Println(lexer.date)
    fmt.Println(lexer.days)
    fmt.Println(lexer.months)
    fmt.Println(lexer.offsets)
    return lexer, 0
}

const (
    HAVE_TIME = 1 << iota
    HAVE_DATE
    HAVE_DAYS
    HAVE_MONTHS
    HAVE_OFFSET
)

func resolve(l *dateLexer, now time.Time) (time.Time, bool) {
    state := 0
    if !l.time.IsZero() {
        state |= HAVE_TIME
    }
    if !l.date.IsZero() {
        state |= HAVE_DATE
    }
    if l.days.seen {
        state |= HAVE_DAYS
    }
    if l.months.seen {
        state |= HAVE_MONTHS
    }
    if l.offsets.seen {
        state |= HAVE_OFFSET
    }
    switch state {
    case HAVE_TIME:
        y, m, d := now.Date()
        h, n, s := l.time.Clock()
        t := time.Date(y, m, d, h, n, s, 0, l.time.Location())
        fmt.Printf("Parsed time as %s %s\n", t.Weekday(), t)
        // check if >24h has been given. Results of this may be *very* sketchy.
        // We can:
        //   a) drop >24h info completely, raise an error/warning
        //   b) save the integer number of hours as "days" and add that
        // Currently, do (a), but (b) would be nice.
        if y, m, d = l.time.Date(); y != 1 || m != time.January || d != 1 {
            // TODO(fluffle): better error reporting!
            fmt.Printf("Time >24h specified, ignoring it")
        }
        return t, true
    case HAVE_DATE:
        y, m, d := l.date.Date()
        if y == 0 {
            y = now.Year()
        }
        h, n, s := now.Clock()
        t := time.Date(y, m, d, h, n, s, 0, now.Location())
        fmt.Printf("Parsed time as %s %s\n", t.Weekday(), t)
        return t, true
    case HAVE_DAYS:
        var t time.Time
        if l.days.year != 0 {
            // this is num'th weekday of year, so start by finding jan 1
            h, n, s := now.Clock()
            t = time.Date(l.days.year, 1, 1, h, n, s, 0, now.Location())
            diff := int(l.days.day - t.Weekday())
            if diff < 0 {
                l.days.num -= 1
            }
            t = t.AddDate(0, 0, l.days.num * 7 + diff)
        } else {
            diff := int(l.days.day - now.Weekday())
            if diff < 0 && l.days.num < 0 {
                l.days.num += 1
            } else if diff > 0 && l.days.num > 0 {
                l.days.num -= 1
            }
            t = now.AddDate(0, 0, l.days.num * 7 + diff)
        }
        fmt.Printf("Parsed time as %s %s\n", t.Weekday(), t)
        return t, true
    }
    return time.Time{}, false
}
