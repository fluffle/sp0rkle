%{
package datetime

// A frontend to time.Parse() to restructure arbitrary dates.
// based upon parse-datetime.y in GNU coreutils.
// also an exercise in learning goyacc in particular.

import (
	"fmt"
//	"math"
	"strconv"
	"strings"
	"time"
	"utf8"
	"unicode"
)

%}

// We essentially pull out substrings from the text and reformat them
// for time.Parse to understand. Makes life easier :-)
%union
{
  strval    string
  intval	int
}

%token <strval> T_AMPM T_PLUS T_MINUS T_MONTH
%token <intval> T_INTEGER

%type <strval> sign zone o_sign o_ampm_or_zone

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
			yylex.(*dateLexer).time = time.SecondsToUTC(int64(-$3))
		} else {
			yylex.(*dateLexer).time = time.SecondsToUTC(int64($3))
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
	time | date;

/*
  | local_zone
	  { pc->local_zones_seen++; }
  | zone
	  { pc->zones_seen++; }
  | date
	  { pc->dates_seen++; }
  | day
	  { pc->days_seen++; }
  | rel
  | number
  | hybrid
  ;
*/

// HH:MM or HH:MM:SS with optional am/pm and/or timezone
time:
	T_INTEGER ':' T_INTEGER o_ampm_or_zone {
		l := yylex.(*dateLexer)
		l.parseTime(
			fmt.Sprintf("%s:04%s%s", l.hourfmt, l.ampmfmt, l.zonefmt),
			fmt.Sprintf("%02d:%02d%s", $1, $3, $4)) 
	}
	| T_INTEGER ':' T_INTEGER ':' T_INTEGER o_ampm_or_zone {
		l := yylex.(*dateLexer)
		l.parseTime(
			fmt.Sprintf("%s:04:05%s%s", l.hourfmt, l.ampmfmt, l.zonefmt),
			fmt.Sprintf("%d:%02d:%02d%s", $1, $3, $5, $6))
	};

o_ampm_or_zone:
	/* empty */ {
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt, l.zonefmt = "15", "", ""
		$$ = ""
	}
	| T_AMPM {
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt, l.zonefmt = "3", $1, ""
	}
	| T_AMPM zone {
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt = "3", $1
		$$ = fmt.Sprintf("%s%s", $1, $2)
	}
	| zone {
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt = "15", ""
	};

zone:
	sign T_INTEGER {
		l := yylex.(*dateLexer)
		l.zonefmt = "-0700"
		$$ = fmt.Sprintf("%s%04d", $1, $2)
	}
	| sign T_INTEGER ':' T_INTEGER {
		l := yylex.(*dateLexer)
		l.zonefmt = "-07:00"
		$$ = fmt.Sprintf("%s%02d:%02d", $1, $2, $4)
	};

/*
local_zone:
	tLOCAL_ZONE
	  {
		pc->local_isdst = $1;
		pc->dsts_seen += (0 < $1);
	  }
  | tLOCAL_ZONE tDST
	  {
		pc->local_isdst = 1;
		pc->dsts_seen += (0 < $1) + 1;
	  }
  ;

// Note 'T' is a special case, as it is used as the separator in ISO
// 8601 date and time of day representation.
zone:
	tZONE
	  { pc->time_zone = $1; }
  | 'T'
	  { pc->time_zone = HOUR(7); }
  | tZONE relunit_snumber
	  { pc->time_zone = $1;
		apply_relative_time (pc, $2, 1); }
  | 'T' relunit_snumber
	  { pc->time_zone = HOUR(7);
		apply_relative_time (pc, $2, 1); }
  | tZONE tSNUMBER o_colon_minutes
	  { pc->time_zone = $1 + time_zone_hhmm (pc, $2, $3); }
  | tDAYZONE
	  { pc->time_zone = $1 + 60; }
  | tZONE tDST
	  { pc->time_zone = $1 + 60; }
  ;

day:
	tDAY
	  {
		pc->day_ordinal = 0;
		pc->day_number = $1;
	  }
  | tDAY ','
	  {
		pc->day_ordinal = 0;
		pc->day_number = $1;
	  }
  | tORDINAL tDAY
	  {
		pc->day_ordinal = $1;
		pc->day_number = $2;
	  }
  | tUNUMBER tDAY
	  {
		pc->day_ordinal = $1.value;
		pc->day_number = $2;
	  }
  ;
*/

datesep:
    T_MINUS | '/';

dayqual:
	/* empty */
	| 's' 't'
	| 'n' 'd'
	| 'r' 'd'
	| 't' 'h';

date:
	T_INTEGER datesep T_INTEGER {
		// DD-MM or MM-YYYY or YYYY-MM
		l := yylex.(*dateLexer)
		if $3 > 12 {
			l.parseDate("1 2006", fmt.Sprintf("%d %04d", $1, $3))
		} else if $1 > 31 {
			l.parseDate("2006 1", fmt.Sprintf("%04d %d", $1, $3))
		} else {
			l.parseDate("2 1", fmt.Sprintf("%d %d", $1, $3))
		}
	}
	| T_INTEGER datesep T_INTEGER datesep T_INTEGER {
		// YYYY-MM-DD or DD-MM-YY(YY?).
		l := yylex.(*dateLexer)
		if $1 > 31 {
			l.parseDate("2006 1 2", fmt.Sprintf("%04d %d %d", $1, $3, $5))
		} else if $3 > 99 {
			l.parseDate("2 1 2006", fmt.Sprintf("%d %d %04d", $1, $3, $5))
		} else {
			l.parseDate("2 1 06", fmt.Sprintf("%d %d %02d", $1, $3, $5))
		}
	}
	| T_INTEGER dayqual T_MONTH {
		// 15th feb
		l := yylex.(*dateLexer)
		l.parseDate("2 Jan", fmt.Sprintf("%d %s", $1, $3))
	}
	| T_MONTH T_INTEGER dayqual {
		// feb 15th or feb 2010
		l := yylex.(*dateLexer)
		if $2 > 31 {
			l.parseDate("Jan 2006", fmt.Sprintf("%s %04d", $1, $2))
		} else {
			l.parseDate("2 Jan", fmt.Sprintf("%d %s", $2, $1))
		}
	}
	| T_INTEGER dayqual T_MONTH T_INTEGER {
		// 15th feb 2010
		l := yylex.(*dateLexer)
		l.parseDate("2 Jan 2006", fmt.Sprintf("%d %s %04d", $1, $3, $4))
	}
	| T_MONTH T_INTEGER dayqual o_comma T_INTEGER {
		// feb 15th 2010
		l := yylex.(*dateLexer)
		l.parseDate("Jan 2 2006", fmt.Sprintf("%s %d %04d", $1, $2, $5))
	};
/*
  | iso_8601_date
  ;

/*
rel:
	relunit tAGO
	  { apply_relative_time (pc, $1, -1); }
  | relunit
	  { apply_relative_time (pc, $1, 1); }
  | dayshift
	  { apply_relative_time (pc, $1, 1); }
  ;

relunit:
	tORDINAL tYEAR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.year = $1; }
  | tUNUMBER tYEAR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.year = $1.value; }
  | tYEAR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.year = 1; }
  | tORDINAL tMONTH_UNIT
	  { $$ = RELATIVE_TIME_0; $$.month = $1; }
  | tUNUMBER tMONTH_UNIT
	  { $$ = RELATIVE_TIME_0; $$.month = $1.value; }
  | tMONTH_UNIT
	  { $$ = RELATIVE_TIME_0; $$.month = 1; }
  | tORDINAL tDAY_UNIT
	  { $$ = RELATIVE_TIME_0; $$.day = $1 * $2; }
  | tUNUMBER tDAY_UNIT
	  { $$ = RELATIVE_TIME_0; $$.day = $1.value * $2; }
  | tDAY_UNIT
	  { $$ = RELATIVE_TIME_0; $$.day = $1; }
  | tORDINAL tHOUR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.hour = $1; }
  | tUNUMBER tHOUR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.hour = $1.value; }
  | tHOUR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.hour = 1; }
  | tORDINAL tMINUTE_UNIT
	  { $$ = RELATIVE_TIME_0; $$.minutes = $1; }
  | tUNUMBER tMINUTE_UNIT
	  { $$ = RELATIVE_TIME_0; $$.minutes = $1.value; }
  | tMINUTE_UNIT
	  { $$ = RELATIVE_TIME_0; $$.minutes = 1; }
  | tORDINAL tSEC_UNIT
	  { $$ = RELATIVE_TIME_0; $$.seconds = $1; }
  | tUNUMBER tSEC_UNIT
	  { $$ = RELATIVE_TIME_0; $$.seconds = $1.value; }
  | tSDECIMAL_NUMBER tSEC_UNIT
	  { $$ = RELATIVE_TIME_0; $$.seconds = $1.tv_sec; $$.ns = $1.tv_nsec; }
  | tUDECIMAL_NUMBER tSEC_UNIT
	  { $$ = RELATIVE_TIME_0; $$.seconds = $1.tv_sec; $$.ns = $1.tv_nsec; }
  | tSEC_UNIT
	  { $$ = RELATIVE_TIME_0; $$.seconds = 1; }
  | relunit_snumber
  ;

relunit_snumber:
	tSNUMBER tYEAR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.year = $1.value; }
  | tSNUMBER tMONTH_UNIT
	  { $$ = RELATIVE_TIME_0; $$.month = $1.value; }
  | tSNUMBER tDAY_UNIT
	  { $$ = RELATIVE_TIME_0; $$.day = $1.value * $2; }
  | tSNUMBER tHOUR_UNIT
	  { $$ = RELATIVE_TIME_0; $$.hour = $1.value; }
  | tSNUMBER tMINUTE_UNIT
	  { $$ = RELATIVE_TIME_0; $$.minutes = $1.value; }
  | tSNUMBER tSEC_UNIT
	  { $$ = RELATIVE_TIME_0; $$.seconds = $1.value; }
  ;

dayshift:
	tDAY_SHIFT
	  { $$ = RELATIVE_TIME_0; $$.day = $1; }
  ;

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

const EPOCH_YEAR = 1970
const eof = 0

type token struct {
	canonical string
	tokentype int
}

var simpleTokenMap = map[string]int{
	"AM": T_AMPM,
	"PM": T_AMPM,
}

var irritatingTokenMap = map[string]token{
	"JAN": token{"Jan", T_MONTH},
	"FEB": token{"Feb", T_MONTH},
	"MAR": token{"Mar", T_MONTH},
	"APR": token{"Apr", T_MONTH},
	"MAY": token{"May", T_MONTH},
	"JUN": token{"Jun", T_MONTH},
	"JUL": token{"Jul", T_MONTH},
	"AUG": token{"Aug", T_MONTH},
	"SEP": token{"Sep", T_MONTH},
	"OCT": token{"Oct", T_MONTH},
	"NOV": token{"Nov", T_MONTH},
	"DEC": token{"Dec", T_MONTH},
}

type dateLexer struct {
	input string
	start, pos, width int
	hourfmt, ampmfmt, zonefmt string
	time, date *time.Time
}

func (l *dateLexer) Lex(lval *yySymType) int {
	l.scan(unicode.IsSpace)
	c := l.peek()
	
	switch {
	case c == '+':
		lval.strval = "+"
		l.next()
		return T_PLUS
	case c == '-':
		lval.strval = "-"
		l.next()
		return T_MINUS
	case unicode.IsDigit(c):
		lval.intval, _ = strconv.Atoi(l.scan(unicode.IsDigit))
		return T_INTEGER
	case unicode.IsLetter(c):
		input := l.scan(unicode.IsLetter)
		if tok, ok := l.lookup(input); ok {
			lval.strval = tok.canonical
			return tok.tokentype
		} else {
			l.rewind()
		}
	}
	return l.next()
}

func (l *dateLexer) Error(e string) {
	fmt.Println(e)
}

func (l *dateLexer) peek() (rune int) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	return rune
}

func (l *dateLexer) next() int {
	l.start = l.pos
	rune := l.peek()
	l.pos += l.width
	return rune
}

func (l *dateLexer) scan(f func(int) bool) string {
	l.start = l.pos
	for f(l.peek()) {
		l.pos += l.width
	}
	str := l.input[l.start:l.pos]
	return str
}

func (l *dateLexer) rewind() {
	l.pos = l.start
}

func (l *dateLexer) lookup(input string) (token, bool) {
	fmt.Printf("Looking up '%s'\n", input)
	input = strings.ToUpper(input)
	// try a simple lookup -- these tokens only ever appear as themselves
	if typ, ok := simpleTokenMap[input]; ok {
		return token{input,typ}, ok
	}
	// Otherwise, it's a more irritating token...
	if tok, ok := irritatingTokenMap[input]; ok {
		return tok, ok
	}
	// strip off a plural?
	if input[len(input)-1] == 'S' {
		if tok, ok := irritatingTokenMap[input[:len(input)-1]]; ok {
			return tok, ok
		}
	}
	// Look up first three letters.
	if len(input) > 3 {
		if tok, ok := irritatingTokenMap[input[:3]]; ok {
			return tok, ok
		}
	}
	return token{}, false
}

func (l *dateLexer) parseTime(fmt, timestr string) {
	if t, err := time.Parse(fmt, timestr); err == nil {
		l.time = t
	} else {
		l.Error(err.String())
	}
}

func (l *dateLexer) parseDate(fmt, timestr string) {
	if t, err := time.Parse(fmt, timestr); err == nil {
		l.date = t
	} else {
		l.Error(err.String())
	}
}

func Parse(input string) *time.Time {
	lexer := &dateLexer{input: input}
	yyDebug = 5
	if ret := yyParse(lexer); ret == 0 {
		fmt.Printf("%#v\n", lexer.time)
		fmt.Printf("%#v\n", lexer.date)
		return lexer.time
	}
	return nil
}
