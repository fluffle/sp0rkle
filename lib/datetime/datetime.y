%{
package datetime

// A frontend to time.Parse() to restructure arbitrary dates.
// based upon parse-datetime.y in GNU coreutils.
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
%token <intval> T_AMPM T_INTEGER T_MONTH T_DAY
%token <zoneval> T_ZONE

%type <strval> sign o_sign dayqual
%type <intval> o_ampm_or_zone

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
	time | date;

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
		hrs, mins := ($2 / 100), ($2 % 100)
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

/*

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
	/* empty */ { $$ = "" }
	| 's' 't' { $$ = "st" }
	| 'n' 'd' { $$ = "nd" }
	| 'r' 'd' { $$ = "rd" }
	| 't' 'h' { $$ = "th" }
	;

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
		} else if $5 > 40 {
			// assume we have DD-MM-YY, add 1900 if YY > 40
			l.setDate($5 + 1900, $3, $1)
		} else {
			// assume we have DD-MM-YY, add 2000 otherwise
			l.setDate($5 + 2000, $3, $1)
		}
	}
	| T_INTEGER dayqual T_MONTH {
		// DDth Mon
		l := yylex.(*dateLexer)
		l.setDate(0, $3, $1)
	}
	| T_MONTH T_INTEGER dayqual {
		l := yylex.(*dateLexer)
		if $2 > 31 && $3 == "" {
			// assume Mon YYYY
			l.setDate($2, $1, 1)
		} else {
		    // assume Mon DDth
			l.setDate(0, $1, $2)
		}
	}
	| T_INTEGER dayqual T_MONTH T_INTEGER {
		l := yylex.(*dateLexer)
		if $4 > 99 {
			// assume DDth Mon YYYY
			l.setDate($4, $3, $1)
		} else if $4 > 40 {
			// assume DDth Mon YY, add 1900 if YY > 40
			l.setDate($4 + 1900, $3, $1)
		} else {
			// assume DDth Mon YY, add 2000 otherwise
			l.setDate($4 + 2000, $3, $1)
		}
	}
	| T_MONTH T_INTEGER dayqual o_comma T_INTEGER {
		l := yylex.(*dateLexer)
		if $5 > 99 {
			// assume Mon DDth, YYYY
			l.setDate($5, $1, $2)
		} else if $5 > 40 {
			// assume Mon DDth, YY, add 1900 if YY > 40
			l.setDate($5 + 1900, $1, $2)
		} else {
			// assume Mon DDth YY, add 2000 otherwise
			l.setDate($5 + 2000, $1, $2)
		}
	};

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

type tokenMap interface {
	Lookup(input string, lval *yySymType) (tokenType int, ok bool)
}

type strMap map[string]struct{
	tokenType int
	tokenVal  string
}

var strTokenMap = strMap{
}

func (stm strMap) Lookup(input string, lval *yySymType) (int, bool) {
	if len(input) > 3 {
		input = input[:3]
	}
	if tok, ok := stm[input]; ok {
		lval.strval = tok.tokenVal
		return tok.tokenType, ok
	}
	return -1, false
}

type numMap map[string]struct{
	tokenType int
	tokenVal int
}

var numTokenMap = numMap{
	"AM": {T_AMPM, 0},
	"PM": {T_AMPM, 12},
	"JAN": {T_MONTH, 1},
	"FEB": {T_MONTH, 2},
	"MAR": {T_MONTH, 3},
	"APR": {T_MONTH, 4},
	"MAY": {T_MONTH, 5},
	"JUN": {T_MONTH, 6},
	"JUL": {T_MONTH, 7},
	"AUG": {T_MONTH, 8},
	"SEP": {T_MONTH, 9},
	"OCT": {T_MONTH, 10},
	"NOV": {T_MONTH, 11},
	"DEC": {T_MONTH, 12},
	"MON": {T_DAY, 1},
	"TUE": {T_DAY, 2},
	"WED": {T_DAY, 3},
	"THU": {T_DAY, 4},
	"FRI": {T_DAY, 5},
	"SAT": {T_DAY, 6},
	"SUN": {T_DAY, 0},
}

func (ntm numMap) Lookup(input string, lval *yySymType) (int, bool) {
	if len(input) > 3 {
		input = input[:3]
	}
	if tok, ok := ntm[input]; ok {
		lval.intval = tok.tokenVal
		return tok.tokenType, ok
	}
	return -1, false
}

type zoneMap map[string]string

var zoneCache = make(map[string]*time.Location)
func zone(loc string) *time.Location {
	if l, ok := zoneCache[loc]; ok {
		return l
	}
	l, _ := time.LoadLocation(loc)
	zoneCache[loc] = l
	return l
}

var zoneTokenMap = zoneMap{
	"ADT": "America/Barbados",
	"AFT": "Asia/Kabul",
	"AKST": "US/Alaska",
	"AKDT": "US/Alaska",
	"AMT": "America/Boa_Vista",
	"ANAT": "Asia/Anadyr",
	"ART": "America/Argentina/Buenos_Aires",
	"AST": "Asia/Qatar",
	"AZOT": "Atlantic/Azores",
	"BNT": "Asia/Brunei",
	"BRT": "Brazil/East",
	"BRST": "Brazil/East",
	"BST": "GB",
	"CAT": "Africa/Harare",
	"CCT": "Indian/Cocos",
	"CDT": "US/Central",
	"CET": "Europe/Zurich",
	"CEST": "Europe/Zurich",
	"CLST": "Chile/Continental",
	"CST": "Asia/Shanghai",
	"EAT": "Africa/Nairobi",
	"EDT": "US/Eastern",
	"EET": "Europe/Athens",
	"EIT": "Asia/Jayapura",
	"EEST": "Europe/Athens",
	"EST": "Australia/Melbourne",
	"FET": "Europe/Kaliningrad",
	"FJT": "Pacific/Fiji",
	"FJST": "Pacific/Fiji",
	"GET": "Asia/Tbilisi",
	"GMT": "GMT",
	"GST": "Asia/Dubai",
	"HADT": "US/Aleutian",
	"HAST": "US/Aleutian",
	"HKT": "Hongkong",
	"HST": "US/Hawaii",
	"ICT": "Asia/Bangkok",
	"IDT": "Asia/Tel_Aviv",
	"IDDT": "Asia/Tel_Aviv",
	"IRDT": "Iran",
	"IRST": "Iran",
	"IOT": "Indian/Chagos",
	"IST": "Asia/Kolkata",
	"JST": "Asia/Tokyo",
	"KGT": "Asia/Bishkek",
	"KST": "Asia/Pyongyang",
	"MDT": "US/Mountain",
	"MART": "Pacific/Marquesas",
	"MET": "MET",
	"MEST": "MET",
	"MMT": "Asia/Rangoon",
	"MST": "US/Mountain",
	"MVT": "Indian/Maldives",
	"MYT": "Asia/Kuala_Lumpur",
	"NDT": "Canada/Newfoundland",
	"NPT": "Asia/Kathmandu",
	"NST": "Canada/Newfoundland",
	"NZDT": "Pacific/Auckland",
	"NZST": "Pacific/Auckland",
	"PDT": "US/Pacific",
	"PHT": "Asia/Manila",
	"PKT": "Asia/Karachi",
	"PST": "US/Pacific",
	"PWT": "Pacific/Palau",
	"RET": "Indian/Reunion",
	"SAST": "Africa/Johannesburg",
	"SCT": "Indian/Mahe",
	"SGT": "Asia/Singapore",
	"SST": "US/Samoa",
	"ULAT": "Asia/Ulaanbaatar",
	"UTC": "UTC",
	"UZT": "Asia/Tashkent",
	"WAT": "Africa/Lagos",
	"WAST": "Africa/Lagos",
	"WET": "WET",
	"WEST": "WET",
	"WIT": "Asia/Jakarta",
	"WST": "Australia/West",
	"VET": "America/Caracas",
	"VLAT": "Asia/Vladivostok",
}

func (ztm zoneMap) Lookup(input string, lval *yySymType) (int, bool) {
	if tok, ok := ztm[input]; ok {
		lval.zoneval = zone(tok)
		return T_ZONE, ok
	}
	return -1, false
}

var tokenMaps = []tokenMap{strTokenMap, numTokenMap, zoneTokenMap}

type dateLexer struct {
	*util.Lexer
	hourfmt, ampmfmt, zonefmt string
	time, date time.Time
	loc *time.Location
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
		for _, m := range tokenMaps {
			if tok, ok := m.Lookup(input, lval); ok {
				return tok
			}
		}
		// If we've not returned yet, no token recognised, so rewind.
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


func Parse(input string) time.Time {
	lexer := &dateLexer{Lexer: &util.Lexer{Input: input}}
	yyDebug = 5
	if ret := yyParse(lexer); ret == 0 {
		fmt.Printf("%s\n", lexer.time)
		fmt.Printf("%s\n", lexer.date)
		return lexer.time
	}
	return time.Time{}
}
