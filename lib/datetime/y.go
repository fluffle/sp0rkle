
//line datetime.y:2
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


//line datetime.y:20
type yySymType struct
{
	yys int
	strval  string
	intval  int
	zoneval *time.Location
}

const T_PLUS = 57346
const T_MINUS = 57347
const T_AMPM = 57348
const T_INTEGER = 57349
const T_MONTH = 57350
const T_DAY = 57351
const T_ZONE = 57352

var yyToknames = []string{
	"T_PLUS",
	"T_MINUS",
	"T_AMPM",
	"T_INTEGER",
	"T_MONTH",
	"T_DAY",
	"T_ZONE",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line datetime.y:329


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

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 38,
	7, 9,
	-2, 35,
}

const yyNprod = 38
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 67

var yyAct = []int{

	39, 41, 15, 17, 14, 16, 19, 22, 20, 37,
	21, 13, 18, 19, 22, 20, 26, 21, 36, 35,
	34, 30, 31, 42, 17, 50, 47, 29, 54, 38,
	40, 4, 33, 18, 52, 44, 30, 31, 42, 30,
	31, 51, 29, 48, 49, 29, 8, 9, 45, 53,
	43, 32, 25, 24, 23, 11, 12, 27, 7, 6,
	5, 46, 3, 2, 1, 10, 28,
}
var yyPact = []int{

	20, -1000, -1000, 39, 51, -1000, -1000, -1000, -2, 47,
	46, -1000, -1000, 45, 35, 44, 24, -1000, -1000, 4,
	1, 0, -11, -9, -1000, 17, -1000, -1000, 43, -1000,
	-1000, -1000, 19, 41, -1000, -1000, -1000, -1000, 14, -1000,
	36, -1000, 35, 12, 34, -1000, 27, -1000, 32, -1000,
	21, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 66, 65, 5, 0, 64, 63, 62, 61, 60,
	59, 58, 1, 57, 2,
}
var yyR1 = []int{

	0, 5, 5, 6, 2, 2, 2, 1, 1, 8,
	8, 7, 7, 9, 9, 10, 10, 10, 4, 4,
	12, 12, 13, 13, 13, 14, 14, 3, 3, 3,
	3, 3, 11, 11, 11, 11, 11, 11,
}
var yyR2 = []int{

	0, 1, 1, 3, 0, 1, 1, 1, 1, 0,
	1, 0, 2, 1, 1, 4, 6, 3, 1, 2,
	0, 1, 2, 4, 1, 1, 1, 0, 2, 2,
	2, 2, 3, 5, 3, 3, 4, 5,
}
var yyChk = []int{

	-1000, -5, -6, -7, 11, -9, -10, -11, 7, 8,
	-2, 4, 5, 13, 6, -14, -3, 5, 14, 15,
	17, 19, 16, 7, 7, 7, -12, -13, -1, 10,
	4, 5, 7, 8, 16, 18, 18, 20, -3, -4,
	13, -12, 6, 7, -14, 7, -8, 12, 7, -12,
	13, 7, 7, -4, 7,
}
var yyDef = []int{

	11, -2, 1, 2, 4, 12, 13, 14, 27, 0,
	0, 5, 6, 0, 20, 0, 0, 25, 26, 0,
	0, 0, 0, 27, 3, 20, 17, 21, 0, 24,
	7, 8, 32, 34, 28, 29, 30, 31, -2, 15,
	0, 18, 20, 22, 0, 36, 0, 10, 20, 19,
	0, 33, 37, 16, 23,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 12, 3, 3, 14, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 13, 3,
	3, 3, 3, 3, 11, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	18, 3, 3, 3, 20, 3, 3, 3, 3, 3,
	17, 3, 3, 3, 19, 15, 16,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10,
}
var yyTok3 = []int{
	0,
}

//line yaccpar:1

/*	parser for yacc output	*/

var yyDebug = 0

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c > 0 && c <= len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return fmt.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return fmt.Sprintf("state-%v", s)
}

func yylex1(lex yyLexer, lval *yySymType) int {
	c := 0
	char := lex.Lex(lval)
	if char <= 0 {
		c = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		c = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			c = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		c = yyTok3[i+0]
		if c == char {
			c = yyTok3[i+1]
			goto out
		}
	}

out:
	if c == 0 {
		c = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		fmt.Printf("lex %U %s\n", uint(char), yyTokname(c))
	}
	return c
}

func yyParse(yylex yyLexer) int {
	var yyn int
	var yylval yySymType
	var yyVAL yySymType
	yyS := make([]yySymType, yyMaxDepth)

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yychar := -1
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		fmt.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yychar < 0 {
		yychar = yylex1(yylex, &yylval)
	}
	yyn += yychar
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yychar { /* valid shift */
		yychar = -1
		yyVAL = yylval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yychar < 0 {
			yychar = yylex1(yylex, &yylval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yychar {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error("syntax error")
			Nerrs++
			if yyDebug >= 1 {
				fmt.Printf("%s", yyStatname(yystate))
				fmt.Printf("saw %s\n", yyTokname(yychar))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				fmt.Printf("error recovery discards %s\n", yyTokname(yychar))
			}
			if yychar == yyEofCode {
				goto ret1
			}
			yychar = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		fmt.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 3:
		//line datetime.y:54
		{
			if yyS[yypt-1].strval == "-" {
				yylex.(*dateLexer).time = time.Unix(int64(-yyS[yypt-0].intval), 0)
			} else {
				yylex.(*dateLexer).time = time.Unix(int64(yyS[yypt-0].intval), 0)
			}
		}
	case 4:
		//line datetime.y:63
		{ yyVAL.strval = "" }
	case 5:
		yyVAL.strval = yyS[yypt-0].strval
	case 6:
		yyVAL.strval = yyS[yypt-0].strval
	case 7:
		yyVAL.strval = yyS[yypt-0].strval
	case 8:
		yyVAL.strval = yyS[yypt-0].strval
	case 15:
		//line datetime.y:82
		{
			l := yylex.(*dateLexer)
			l.setTime(yyS[yypt-3].intval + yyS[yypt-0].intval, yyS[yypt-1].intval, 0, l.loc)
		}
	case 16:
		//line datetime.y:86
		{
			l := yylex.(*dateLexer)
			l.setTime(yyS[yypt-5].intval + yyS[yypt-0].intval, yyS[yypt-3].intval, yyS[yypt-1].intval, l.loc)
		}
	case 17:
		//line datetime.y:90
		{
			l := yylex.(*dateLexer)
			l.setTime(yyS[yypt-2].intval + yyS[yypt-1].intval, 0, 0, l.loc)
		}
	case 18:
		//line datetime.y:96
		{
			yyVAL.intval = 0
		}
	case 19:
		//line datetime.y:99
		{
			yyVAL.intval = yyS[yypt-1].intval
		}
	case 22:
		//line datetime.y:108
		{
			l := yylex.(*dateLexer)
			hrs, mins := (yyS[yypt-0].intval / 100), (yyS[yypt-0].intval % 100)
			if (yyS[yypt-1].strval == "-") {
				l.loc = time.FixedZone("WTF", -3600 * hrs - 60 * mins)
			} else {
				l.loc = time.FixedZone("WTF", 3600 * hrs + 60 * mins)
			}   
		}
	case 23:
		//line datetime.y:117
		{
			l := yylex.(*dateLexer)
			if (yyS[yypt-3].strval == "-") {
				l.loc = time.FixedZone("WTF", -3600 * yyS[yypt-2].intval - 60 * yyS[yypt-0].intval)
			} else {
				l.loc = time.FixedZone("WTF", 3600 * yyS[yypt-2].intval + 60 * yyS[yypt-0].intval)
			}   
		}
	case 24:
		//line datetime.y:125
		{
			l := yylex.(*dateLexer)
			l.loc = yyS[yypt-0].zoneval
		}
	case 27:
		//line datetime.y:160
		{ yyVAL.strval = "" }
	case 28:
		//line datetime.y:161
		{ yyVAL.strval = "st" }
	case 29:
		//line datetime.y:162
		{ yyVAL.strval = "nd" }
	case 30:
		//line datetime.y:163
		{ yyVAL.strval = "rd" }
	case 31:
		//line datetime.y:164
		{ yyVAL.strval = "th" }
	case 32:
		//line datetime.y:168
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].intval > 12 {
				// assume we have MM-YYYY
			l.setDate(yyS[yypt-0].intval, yyS[yypt-2].intval, 1)
			} else if yyS[yypt-2].intval > 31 {
				// assume we have YYYY-MM
			l.setDate(yyS[yypt-2].intval, yyS[yypt-0].intval, 1)
			} else {
				// assume we have DD-MM (too bad, americans)
			l.setDate(0, yyS[yypt-0].intval, yyS[yypt-2].intval)
			}
		}
	case 33:
		//line datetime.y:181
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-4].intval > 31 {
				// assume we have YYYY-MM-DD
			l.setDate(yyS[yypt-4].intval, yyS[yypt-2].intval, yyS[yypt-0].intval)
			} else if yyS[yypt-0].intval > 99 {
				// assume we have DD-MM-YYYY
			l.setDate(yyS[yypt-0].intval, yyS[yypt-2].intval, yyS[yypt-4].intval)
			} else if yyS[yypt-0].intval > 40 {
				// assume we have DD-MM-YY, add 1900 if YY > 40
			l.setDate(yyS[yypt-0].intval + 1900, yyS[yypt-2].intval, yyS[yypt-4].intval)
			} else {
				// assume we have DD-MM-YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].intval + 2000, yyS[yypt-2].intval, yyS[yypt-4].intval)
			}
		}
	case 34:
		//line datetime.y:197
		{
			// DDth Mon
		l := yylex.(*dateLexer)
			l.setDate(0, yyS[yypt-0].intval, yyS[yypt-2].intval)
		}
	case 35:
		//line datetime.y:202
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-1].intval > 31 && yyS[yypt-0].strval == "" {
				// assume Mon YYYY
			l.setDate(yyS[yypt-1].intval, yyS[yypt-2].intval, 1)
			} else {
			    // assume Mon DDth
			l.setDate(0, yyS[yypt-2].intval, yyS[yypt-1].intval)
			}
		}
	case 36:
		//line datetime.y:212
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].intval > 99 {
				// assume DDth Mon YYYY
			l.setDate(yyS[yypt-0].intval, yyS[yypt-1].intval, yyS[yypt-3].intval)
			} else if yyS[yypt-0].intval > 40 {
				// assume DDth Mon YY, add 1900 if YY > 40
			l.setDate(yyS[yypt-0].intval + 1900, yyS[yypt-1].intval, yyS[yypt-3].intval)
			} else {
				// assume DDth Mon YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].intval + 2000, yyS[yypt-1].intval, yyS[yypt-3].intval)
			}
		}
	case 37:
		//line datetime.y:225
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].intval > 99 {
				// assume Mon DDth, YYYY
			l.setDate(yyS[yypt-0].intval, yyS[yypt-4].intval, yyS[yypt-3].intval)
			} else if yyS[yypt-0].intval > 40 {
				// assume Mon DDth, YY, add 1900 if YY > 40
			l.setDate(yyS[yypt-0].intval + 1900, yyS[yypt-4].intval, yyS[yypt-3].intval)
			} else {
				// assume Mon DDth YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].intval + 2000, yyS[yypt-4].intval, yyS[yypt-3].intval)
			}
		}
	}
	goto yystack /* stack new state and value */
}
