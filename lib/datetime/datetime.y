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

%token <strval> T_PLUS T_MINUS
%token <intval> T_AMPM T_INTEGER T_MONTHNAME T_DAYNAME
%token <intval> T_OFFSET T_DAY T_RELATIVE T_DAYSHIFT T_AGO
%token <zoneval> T_ZONE

%type <strval> sign o_sign dayqual
%type <intval> o_ampm_or_zone o_ago o_relspec

/*
%token tAGO tDST

%token tYEAR_UNIT tMONTH_UNIT tHOUR_UNIT tMINUTE_UNIT tSEC_UNIT
%token <intval> tDAY_UNIT tDAY_SHIFT

//%token <intval> tDAY tDAYZONE tLOCAL_ZONE tMERIDIAN
/%token <intval> tMONTH tORDINAL tZONE


%type <intval> o_colon_minutes
%type <intval> relunit relunit_snumber dayshift
*/

%%

spec:
	unixtime | items;

unixtime:
	'@' o_sign T_INTEGER {
		if $2 == "-" {
			yylex.(*dateLexer).time = time.Unix(int64(-$3), 0)
		} else {
			yylex.(*dateLexer).time = time.Unix(int64($3), 0)
		}
	};

o_sign:
	/* empty */	{ $$ = "" }
	| T_PLUS | T_MINUS;

sign:
	T_PLUS | T_MINUS; 

o_comma:
	/* empty */ | ',';

items:
	/* empty */
	| items item;

item:
/*	time
	| 
	date
	| 
	day
	| 
	month
	| */
	relative;

// HH:MM or HH:MM:SS with optional am/pm and/or timezone
// *or* HH with non-optional am/pm
time:
	T_INTEGER ':' T_INTEGER o_ampm_or_zone {
		l := yylex.(*dateLexer)
		l.setTime($1 + $4, $3, 0, l.loc)
	}
	| T_INTEGER ':' T_INTEGER ':' T_INTEGER o_ampm_or_zone {
		l := yylex.(*dateLexer)
		l.setTime($1 + $6, $3, $5, l.loc)
	}
	| T_INTEGER T_AMPM o_zone {
		l := yylex.(*dateLexer)
		l.setTime($1 + $2, 0, 0, l.loc)
	};

o_ampm_or_zone:
	o_zone {
		$$ = 0
	}
	| T_AMPM o_zone {
		$$ = $1
	};

o_zone:
	/* empty */
	| zone;

zone:
	sign T_INTEGER {
		l := yylex.(*dateLexer)
        hrs, mins := $2, 0
        if (hrs > 100) {
            hrs, mins = ($2 / 100), ($2 % 100)
        } else {
            hrs *= 100
        }
		if ($1 == "-") {
			l.loc = time.FixedZone("WTF", -3600 * hrs - 60 * mins)
		} else {
			l.loc = time.FixedZone("WTF", 3600 * hrs + 60 * mins)
		}   
	}
	| sign T_INTEGER ':' T_INTEGER {
		l := yylex.(*dateLexer)
		if ($1 == "-") {
			l.loc = time.FixedZone("WTF", -3600 * $2 - 60 * $4)
		} else {
			l.loc = time.FixedZone("WTF", 3600 * $2 + 60 * $4)
		}   
	}
	| T_ZONE {
		l := yylex.(*dateLexer)
		l.loc = $1
	};

datesep:
	T_MINUS | '/';

dayqual:
	/* empty */ { $$ = "" }
	| 's' 't' { $$ = "st" }
	| 'n' 'd' { $$ = "nd" }
	| 'r' 'd' { $$ = "rd" }
	| 't' 'h' { $$ = "th" };

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
	| T_INTEGER dayqual T_MONTHNAME {
		// DDth Mon
		l := yylex.(*dateLexer)
		l.setDate(0, $3, $1)
	}
	| T_MONTHNAME T_INTEGER dayqual {
		l := yylex.(*dateLexer)
		if $2 > 31 && $3 == "" {
			// assume Mon YYYY
			l.setDate($2, $1, 1)
		} else {
		    // assume Mon DDth
			l.setDate(0, $1, $2)
		}
	}
	| T_INTEGER dayqual T_MONTHNAME T_INTEGER {
		l := yylex.(*dateLexer)
		if $4 > 99 {
			// assume DDth Mon YYYY
			l.setDate($4, $3, $1)
		} else if $4 > 68 {
			// assume DDth Mon YY, add 1900 if YY > 68
			l.setDate($4 + 1900, $3, $1)
		} else {
			// assume DDth Mon YY, add 2000 otherwise
			l.setDate($4 + 2000, $3, $1)
		}
	}
	| T_MONTHNAME T_INTEGER dayqual o_comma T_INTEGER {
		l := yylex.(*dateLexer)
		if $5 > 99 {
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

day:
	o_relspec T_DAYNAME o_comma {
		l := yylex.(*dateLexer)
		l.setDay($2, $1)
	}
	| T_DAYSHIFT {
		l := yylex.(*dateLexer)
		// translate "tomorrow" and "yesterday" to weekdays
		d := time.Now().Weekday()
		l.setDay((7+int(d)+$1)%7, $1)
	};

month:
	o_relspec T_MONTHNAME {
		l := yylex.(*dateLexer)
		l.setMonth($2, $1)
	};

o_relspec:
	/* empty */ { $$ = 1 }
	| o_sign T_INTEGER {
		if $1 == "-" {
			$$ = -$2
		} else {
			$$ = $2
		}
	}
	| T_RELATIVE;

o_ago:
	/* empty */ { $$ = 1 }
	| T_AGO;

relative:
	relunits o_ago;

relunits:
	/* empty */
	| relunits relunit;

relunit:
	o_relspec T_OFFSET {
		l := yylex.(*dateLexer)
		l.addOffset(offset($2), $1)
	}
	| o_relspec T_DAY {
		// Special-case to handle "week" and "fortnight"
		l := yylex.(*dateLexer)
		l.addOffset(O_DAY, $1 * $2)
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
type relTime [6]int
func (rt relTime) String() string {
	s := ""
	for off, val := range rt {
		if val != 0 {
			s += fmt.Sprintf("%d %s ", val, offsets[off])
		}
	}
	return s[:len(s)-1]
}

type relDays struct {
	day time.Weekday
	num int
	seen bool
}
func (rd relDays) String() string {
	if !rd.seen {
		return "No relative days seen"
	}
	return fmt.Sprintf("%d %s", rd.num, rd.day)
}

type relMonths struct {
	month time.Month
	num int
	seen bool
}
func (rm relMonths) String() string {
	if !rm.seen {
		return "No relative months seen"
	}
	return fmt.Sprintf("%d %s", rm.num, rm.month)
}

type dateLexer struct {
	*util.Lexer
	hourfmt, ampmfmt, zonefmt string
	time, date time.Time
	offset  relTime       // takes care of +- ymd hms
    days    relDays       // takes care of specific days into future
	months  relMonths     // takes care of specific months into future
	loc     *time.Location
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

func (l *dateLexer) setDay(d, n int) {
	fmt.Printf("Setting day to %d %s\n", n, time.Weekday(d))
	if l.days.seen {
		l.Error("Parsed two days")
	}
	l.days = relDays{time.Weekday(d), n, true}
}

func (l *dateLexer) setMonth(m, n int) {
	fmt.Printf("Setting month to %d %s\n", n, time.Month(m))
	if l.months.seen {
		l.Error("Parsed two months")
	}
	l.months = relMonths{time.Month(m), n, true}
}

func (l *dateLexer) addOffset(off offset, rel int) {
	fmt.Printf("Adding relative offset of %d %s\n", rel, off)
	l.offset[off] += rel;
}

func (l *dateLexer) setAgo(ago int) {
	for i := range l.offset {
		l.offset[i] *= ago
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
		fmt.Println(lexer.offset)
		return lexer.time
	}
	return time.Time{}
}
