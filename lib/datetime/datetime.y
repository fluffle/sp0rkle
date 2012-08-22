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

%}

%union
{
	strval  string
	intval  int
	zoneval *time.Location
}

%token T_PLUS T_MINUS
%token <intval> T_AMPM T_INTEGER T_MONTHNAME T_DAYNAME
%token <intval> T_OFFSET T_DAY T_RELATIVE T_DAYSHIFT T_AGO
%token <zoneval> T_ZONE

%type <intval> sign o_sign_integer o_ampm
%type <zoneval> o_zone

%%

spec:
	unixtime | items;

sign:
	T_PLUS { $$ = 1 }
	| T_MINUS { $$ = -1 }; 

o_sign_integer:
	T_INTEGER
	| T_PLUS T_INTEGER { $$ = $2 }
	| T_MINUS T_INTEGER { $$ = -$2 };

o_comma:
	/* empty */ | ',';

datesep:
	T_MINUS | '/';

dayqual:
	's' 't'
	| 'n' 'd'
	| 'r' 'd'
	| 't' 'h';

o_dayqual: /* empty */ | dayqual;

of: 'o' 'f';

o_of: /* empty */ | of;

unixtime:
	'@' o_sign_integer {
		yylex.(*dateLexer).time = time.Unix(int64($2), 0)
	};

items:
	/* empty */
	| items item;

item:
	time
	| date
	| day_or_month
	| relative;

// HH:MM or HH:MM:SS with optional am/pm and/or timezone
// *or* HH with non-optional am/pm
time:
	T_INTEGER ':' T_INTEGER o_ampm o_zone {
		yylex.(*dateLexer).setTime($1 + $4, $3, 0, $5)
	}
	| T_INTEGER ':' T_INTEGER ':' T_INTEGER o_ampm o_zone {
		yylex.(*dateLexer).setTime($1 + $6, $3, $5, $7)
	}
	| T_INTEGER T_AMPM o_zone {
		yylex.(*dateLexer).setTime($1 + $2, 0, 0, $3)
	};

o_ampm:
	/* empty */ { $$ = 0 }
	| T_AMPM;

o_zone:
	/* empty */ { $$ = nil }
	| T_ZONE
	| sign T_INTEGER {
        hrs, mins := $2, 0
        if (hrs > 100) {
            hrs, mins = ($2 / 100), ($2 % 100)
        } else {
            hrs *= 100
        }
		$$ = time.FixedZone("WTF", $1 * (3600 * hrs + 60 * mins))
	}
	| sign T_INTEGER ':' T_INTEGER {
		$$ = time.FixedZone("WTF", $1 * (3600 * $2 + 60 * $4))
	};

date:
	T_INTEGER datesep T_INTEGER {
		l := yylex.(*dateLexer)
		if $3 > 12 {
			// assume we have MM-YYYY
			l.setDate($3, $1, 1)
		} else if $1 > 31 {
			// assume we have YYYY-MM
			l.setDate($1, $3, 1)
		} else {
			// assume we have DD-MM (too bad, americans)
			l.setDate(0, $3, $1)
		}
	}
	| T_INTEGER datesep T_INTEGER datesep T_INTEGER {
		l := yylex.(*dateLexer)
		if $1 > 31 {
			// assume we have YYYY-MM-DD
			l.setDate($1, $3, $5)
		} else if $5 > 99 {
			// assume we have DD-MM-YYYY
			l.setDate($5, $3, $1)
		} else if $5 > 68 {
			// assume we have DD-MM-YY, add 1900 if YY > 68
			l.setDate($5 + 1900, $3, $1)
		} else {
			// assume we have DD-MM-YY, add 2000 otherwise
			l.setDate($5 + 2000, $3, $1)
		}
	}
	| T_INTEGER o_dayqual o_of T_MONTHNAME {
		// DDth of Mon
		yylex.(*dateLexer).setDate(0, $4, $1)
	}
	| T_MONTHNAME T_INTEGER o_dayqual {
		l := yylex.(*dateLexer)
		if $2 > 999 {
			// assume Mon YYYY
			l.setDate($2, $1, 1)
		} else if $2 <= 31 {
		    // assume Mon DDth
			l.setDate(0, $1, $2)
		} else {
			l.Error("Ambiguous T_MONTHNAME T_INTEGER")
		}
	}
	| T_INTEGER o_dayqual o_of T_MONTHNAME T_INTEGER {
		l := yylex.(*dateLexer)
		if $4 > 999 {
			// assume DDth of Mon YYYY
			l.setDate($5, $4, $1)
		} else if $4 > 68 {
			// assume DDth of Mon YY, add 1900 if YY > 68
			l.setDate($5 + 1900, $4, $1)
		} else {
			// assume DDth of Mon YY, add 2000 otherwise
			l.setDate($5 + 2000, $4, $1)
		}
	}
	| T_MONTHNAME T_INTEGER o_dayqual o_comma T_INTEGER {
		l := yylex.(*dateLexer)
		if $5 > 999 {
			// assume Mon DDth, YYYY
			l.setDate($5, $1, $2)
		} else if $5 > 68 {
			// assume Mon DDth, YY, add 1900 if YY > 68
			l.setDate($5 + 1900, $1, $2)
		} else {
			// assume Mon DDth YY, add 2000 otherwise
			l.setDate($5 + 2000, $1, $2)
		}
	};

day_or_month:
	T_DAYNAME o_comma {
		// Tuesday,
		yylex.(*dateLexer).setDay($1, 1)
	}
	T_MONTHNAME {
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
		yylex.(*dateLexer).setDay($2, $1)
	}
	| T_INTEGER dayqual T_DAYNAME {
		// 3rd Tuesday 
		yylex.(*dateLexer).setDay($3, $1)
	}
	| T_INTEGER dayqual T_DAYNAME of T_MONTHNAME {
		// 3rd Tuesday of (implicit this) March
		l := yylex.(*dateLexer)
		l.setDay($3, $1)
		l.setMonth($5, 1)
	}
	| T_INTEGER dayqual T_DAYNAME of T_INTEGER {
		// 3rd Tuesday of 2012
		yylex.(*dateLexer).setDay($3, $1, $5)
	}
	| T_INTEGER dayqual T_DAYNAME of T_MONTHNAME T_INTEGER {
		// 3rd Tuesday of March 2012
		l := yylex.(*dateLexer)
		l.setDay($3, $1)
		l.setMonth($5, 1, $6)
	}
	| T_INTEGER dayqual T_DAYNAME of T_RELATIVE T_MONTHNAME {
		// 3rd Tuesday of next March
		l := yylex.(*dateLexer)
		l.setDay($3, $1)
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
		yylex.(*dateLexer).addOffset(offset($2), $1)
	}
	| T_RELATIVE T_OFFSET {
		yylex.(*dateLexer).addOffset(offset($2), $1)
	} 
	| o_sign_integer T_DAY {
		// Special-case to handle "week" and "fortnight"
		yylex.(*dateLexer).addOffset(O_DAY, $1 * $2)
	}
	| T_RELATIVE T_DAY {
		yylex.(*dateLexer).addOffset(O_DAY, $1 * $2)
	};



/*
number:
	tUNUMBER
	  { digits_to_date_time (pc, $1); }
  ;

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
}


func (l *dateLexer) Lex(lval *yySymType) int {
	l.Scan(unicode.IsSpace)
	c := l.Peek()
	
	switch {
	case c == '+':
		lval.strval = "+"
		l.Next()
		return T_PLUS
	case c == '-':
		lval.strval = "-"
		l.Next()
		return T_MINUS
	case unicode.IsDigit(c):
		lval.intval, _ = strconv.Atoi(l.Scan(unicode.IsDigit))
		return T_INTEGER
	case unicode.IsLetter(c):
		input := strings.ToUpper(l.Scan(unicode.IsLetter))
		fmt.Printf("Map lookup: %s\n", input)
		// These maps are defined in tokenmaps.go
		for _, m := range tokenMaps {
			if tok, ok := m.Lookup(input, lval); ok {
				fmt.Printf("Map got: %d %d\n", lval.intval, tok)
				return tok
			}
		}
		// If we've not returned yet, no token recognised, so rewind.
		fmt.Printf("Map lookup failed\n")
		l.Rewind()
	}
	l.Next()
	return int(c)
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

func (l *dateLexer) setDate(y, m, d int) {
	fmt.Printf("Setting date to %d-%d-%d\n", y, m, d)
	if ! l.date.IsZero() {
		l.Error("Parsed two dates")
	}
	l.date = time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
}

func (l *dateLexer) setDay(d, n int, year ...int) {
	fmt.Printf("Setting day to %d %s\n", n, time.Weekday(d))
	if l.days.seen {
		l.Error("Parsed two days")
	}
	l.days = relDays{time.Weekday(d), n, 0, true}
	if len(year) > 0 {
		l.days.year = year[0]
	}
}

func (l *dateLexer) setMonth(m, n int, year ...int) {
	fmt.Printf("Setting month to %d %s\n", n, time.Month(m))
	if l.months.seen {
		l.Error("Parsed two months")
	}
	l.months = relMonths{time.Month(m), n, 0, true}
	if len(year) > 0 {
		l.months.year = year[0]
	}
}

func (l *dateLexer) addOffset(off offset, rel int) {
	fmt.Printf("Adding relative offset of %d %s\n", rel, off)
	l.offsets.seen = true
	l.offsets.offsets[off] += rel
}

func (l *dateLexer) setAgo() {
	for i := range l.offsets.offsets {
		l.offsets.offsets[i] *= -1
	}
}

func Parse(input string) time.Time {
	lexer := &dateLexer{Lexer: &util.Lexer{Input: input}}
	yyDebug = 5
	if ret := yyParse(lexer); ret == 0 {
		fmt.Println(lexer.time)
		fmt.Println(lexer.date)
		fmt.Println(lexer.days)
		fmt.Println(lexer.months)
		fmt.Println(lexer.offsets)
		return lexer.time
	}
	return time.Time{}
}
