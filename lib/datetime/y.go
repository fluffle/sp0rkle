
//line datetime.y:2
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


//line datetime.y:19
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
const T_MONTHNAME = 57350
const T_DAYNAME = 57351
const T_OFFSET = 57352
const T_DAY = 57353
const T_RELATIVE = 57354
const T_DAYSHIFT = 57355
const T_AGO = 57356
const T_ZONE = 57357

var yyToknames = []string{
	"T_PLUS",
	"T_MINUS",
	"T_AMPM",
	"T_INTEGER",
	"T_MONTHNAME",
	"T_DAYNAME",
	"T_OFFSET",
	"T_DAY",
	"T_RELATIVE",
	"T_DAYSHIFT",
	"T_AGO",
	"T_ZONE",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line datetime.y:284


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

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 10,
	8, 32,
	26, 32,
	-2, 5,
	-1, 64,
	7, 8,
	-2, 40,
}

const yyNprod = 62
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 106

var yyAct = []int{

	50, 67, 34, 57, 24, 25, 26, 27, 23, 73,
	29, 32, 30, 58, 31, 62, 63, 61, 60, 27,
	22, 28, 29, 32, 30, 14, 31, 16, 79, 69,
	20, 53, 54, 28, 4, 35, 43, 59, 89, 64,
	65, 68, 51, 38, 39, 47, 76, 46, 17, 18,
	72, 10, 11, 12, 88, 87, 13, 15, 17, 18,
	71, 21, 85, 74, 83, 82, 48, 75, 77, 84,
	37, 36, 38, 39, 40, 41, 42, 41, 42, 81,
	86, 17, 18, 80, 21, 78, 70, 90, 55, 49,
	45, 44, 33, 69, 19, 66, 56, 9, 8, 7,
	6, 5, 3, 2, 1, 52,
}
var yyPact = []int{

	17, -1000, -1000, 44, 77, -1000, -1000, -1000, -1000, -1000,
	2, 85, 19, 62, 65, -1000, 22, 84, 83, 54,
	-1000, -1000, 82, 27, 81, -13, 28, -1000, -1000, -3,
	-6, -8, -9, -10, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 67, 33, 23,
	-1000, -1000, 79, -1000, -1000, 14, 42, -1000, -18, -13,
	-1000, -1000, -1000, -1000, 19, -1000, 38, 27, 78, -1000,
	10, 76, 72, -1000, 57, 55, -1000, -1000, 87, 48,
	-1000, -1000, 47, -1000, 30, -1000, 27, -1000, -1000, -1000,
	-1000,
}
var yyPgo = []int{

	0, 105, 25, 1, 0, 104, 103, 102, 2, 101,
	100, 99, 98, 97, 4, 6, 5, 3, 96, 95,
	27, 94,
}
var yyR1 = []int{

	0, 5, 5, 1, 1, 2, 2, 2, 8, 8,
	6, 7, 7, 9, 9, 9, 9, 10, 10, 10,
	3, 3, 4, 4, 4, 4, 14, 14, 15, 15,
	15, 15, 16, 16, 17, 18, 18, 11, 11, 11,
	11, 11, 11, 19, 12, 12, 12, 12, 12, 12,
	12, 12, 12, 12, 13, 13, 20, 20, 21, 21,
	21, 21,
}
var yyR2 = []int{

	0, 1, 1, 1, 1, 1, 2, 2, 0, 1,
	2, 0, 2, 1, 1, 1, 1, 5, 7, 3,
	0, 1, 0, 1, 2, 4, 1, 1, 2, 2,
	2, 2, 0, 1, 2, 0, 1, 3, 5, 4,
	3, 5, 5, 0, 4, 2, 2, 2, 3, 5,
	5, 6, 6, 1, 1, 2, 1, 2, 2, 2,
	2, 2,
}
var yyChk = []int{

	-1000, -5, -6, -7, 17, -9, -10, -11, -12, -13,
	7, 8, 9, 12, -2, 13, -20, 4, 5, -21,
	-2, 7, 18, 6, -14, -16, -15, 5, 19, 20,
	22, 24, 21, 7, -8, 16, 9, 8, 10, 11,
	9, 10, 11, 14, 7, 7, -20, -2, 12, 7,
	-4, 15, -1, 4, 5, 7, -18, -17, 26, 9,
	21, 23, 23, 25, -16, -15, -19, -3, 18, 6,
	7, -14, 8, 27, -17, -8, 8, -4, 7, 18,
	7, 7, 8, 7, 12, 7, -3, 7, 7, 8,
	-4,
}
var yyDef = []int{

	11, -2, 1, 2, 0, 12, 13, 14, 15, 16,
	-2, 0, 8, 0, 0, 53, 54, 0, 0, 56,
	10, 5, 0, 22, 0, 35, 33, 26, 27, 0,
	0, 0, 0, 32, 43, 9, 45, 46, 59, 61,
	47, 58, 60, 55, 6, 7, 57, 0, 0, 20,
	19, 23, 0, 3, 4, 37, 0, 36, 0, 48,
	28, 29, 30, 31, -2, 33, 0, 22, 0, 21,
	24, 0, 39, 34, 0, 0, 44, 17, 20, 0,
	38, 41, 49, 50, 0, 42, 22, 25, 51, 52,
	18,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 16, 3, 3, 19, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 18, 3,
	3, 3, 3, 3, 17, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	23, 3, 27, 3, 25, 3, 3, 3, 3, 3,
	22, 26, 3, 3, 24, 20, 21,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15,
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
		//line datetime.y:40
		{ yyVAL.intval = 1 }
	case 4:
		//line datetime.y:41
		{ yyVAL.intval = -1 }
	case 5:
		yyVAL.intval = yyS[yypt-0].intval
	case 6:
		//line datetime.y:45
		{ yyVAL.intval = yyS[yypt-0].intval }
	case 7:
		//line datetime.y:46
		{ yyVAL.intval = -yyS[yypt-0].intval }
	case 10:
		//line datetime.y:52
		{
			yylex.(*dateLexer).time = time.Unix(int64(yyS[yypt-0].intval), 0)
		}
	case 17:
		//line datetime.y:69
		{
			yylex.(*dateLexer).setTime(yyS[yypt-4].intval + yyS[yypt-1].intval, yyS[yypt-2].intval, 0, yyS[yypt-0].zoneval)
		}
	case 18:
		//line datetime.y:72
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].intval + yyS[yypt-1].intval, yyS[yypt-4].intval, yyS[yypt-2].intval, yyS[yypt-0].zoneval)
		}
	case 19:
		//line datetime.y:75
		{
			yylex.(*dateLexer).setTime(yyS[yypt-2].intval + yyS[yypt-1].intval, 0, 0, yyS[yypt-0].zoneval)
		}
	case 20:
		//line datetime.y:80
		{ yyVAL.intval = 0 }
	case 21:
		yyVAL.intval = yyS[yypt-0].intval
	case 22:
		//line datetime.y:84
		{ yyVAL.zoneval = nil }
	case 23:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 24:
		//line datetime.y:86
		{
	        hrs, mins := yyS[yypt-0].intval, 0
	        if (hrs > 100) {
	            hrs, mins = (yyS[yypt-0].intval / 100), (yyS[yypt-0].intval % 100)
	        } else {
	            hrs *= 100
	        }
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-1].intval * (3600 * hrs + 60 * mins))
		}
	case 25:
		//line datetime.y:95
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-3].intval * (3600 * yyS[yypt-2].intval + 60 * yyS[yypt-0].intval))
		}
	case 37:
		//line datetime.y:115
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
	case 38:
		//line datetime.y:128
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-4].intval > 31 {
				// assume we have YYYY-MM-DD
			l.setDate(yyS[yypt-4].intval, yyS[yypt-2].intval, yyS[yypt-0].intval)
			} else if yyS[yypt-0].intval > 99 {
				// assume we have DD-MM-YYYY
			l.setDate(yyS[yypt-0].intval, yyS[yypt-2].intval, yyS[yypt-4].intval)
			} else if yyS[yypt-0].intval > 68 {
				// assume we have DD-MM-YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].intval + 1900, yyS[yypt-2].intval, yyS[yypt-4].intval)
			} else {
				// assume we have DD-MM-YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].intval + 2000, yyS[yypt-2].intval, yyS[yypt-4].intval)
			}
		}
	case 39:
		//line datetime.y:144
		{
			// DDth of Mon
		yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].intval)
		}
	case 40:
		//line datetime.y:148
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-1].intval > 999 {
				// assume Mon YYYY
			l.setDate(yyS[yypt-1].intval, yyS[yypt-2].intval, 1)
			} else if yyS[yypt-1].intval <= 31 {
			    // assume Mon DDth
			l.setDate(0, yyS[yypt-2].intval, yyS[yypt-1].intval)
			} else {
				l.Error("Ambiguous T_MONTHNAME T_INTEGER")
			}
		}
	case 41:
		//line datetime.y:160
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-1].intval > 999 {
				// assume DDth of Mon YYYY
			l.setDate(yyS[yypt-0].intval, yyS[yypt-1].intval, yyS[yypt-4].intval)
			} else if yyS[yypt-1].intval > 68 {
				// assume DDth of Mon YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].intval + 1900, yyS[yypt-1].intval, yyS[yypt-4].intval)
			} else {
				// assume DDth of Mon YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].intval + 2000, yyS[yypt-1].intval, yyS[yypt-4].intval)
			}
		}
	case 42:
		//line datetime.y:173
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].intval > 999 {
				// assume Mon DDth, YYYY
			l.setDate(yyS[yypt-0].intval, yyS[yypt-4].intval, yyS[yypt-3].intval)
			} else if yyS[yypt-0].intval > 68 {
				// assume Mon DDth, YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].intval + 1900, yyS[yypt-4].intval, yyS[yypt-3].intval)
			} else {
				// assume Mon DDth YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].intval + 2000, yyS[yypt-4].intval, yyS[yypt-3].intval)
			}
		}
	case 43:
		//line datetime.y:188
		{
			// Tuesday,
		yylex.(*dateLexer).setDay(yyS[yypt-1].intval, 1)
		}
	case 44:
		//line datetime.y:192
		{
			// March
		yylex.(*dateLexer).setMonth(yyS[yypt-3].intval, 1)
		}
	case 45:
		//line datetime.y:196
		{
			// Next tuesday
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 46:
		//line datetime.y:200
		{
			// Next march
		yylex.(*dateLexer).setMonth(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 47:
		//line datetime.y:204
		{
			// +-N Tuesdays
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 48:
		//line datetime.y:208
		{
			// 3rd Tuesday 
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-2].intval)
		}
	case 49:
		//line datetime.y:212
		{
			// 3rd Tuesday of (implicit this) March
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-2].intval, yyS[yypt-4].intval)
			l.setMonth(yyS[yypt-0].intval, 1)
		}
	case 50:
		//line datetime.y:218
		{
			// 3rd Tuesday of 2012
		yylex.(*dateLexer).setDay(yyS[yypt-2].intval, yyS[yypt-4].intval, yyS[yypt-0].intval)
		}
	case 51:
		//line datetime.y:222
		{
			// 3rd Tuesday of March 2012
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-3].intval, yyS[yypt-5].intval)
			l.setMonth(yyS[yypt-1].intval, 1, yyS[yypt-0].intval)
		}
	case 52:
		//line datetime.y:228
		{
			// 3rd Tuesday of next March
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-3].intval, yyS[yypt-5].intval)
			l.setMonth(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 53:
		//line datetime.y:234
		{
			// yesterday or tomorrow
		d := time.Now().Weekday()
			yylex.(*dateLexer).setDay((7+int(d)+yyS[yypt-0].intval)%7, yyS[yypt-0].intval)
		}
	case 55:
		//line datetime.y:242
		{
			yylex.(*dateLexer).setAgo()
		}
	case 58:
		//line datetime.y:251
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 59:
		//line datetime.y:254
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 60:
		//line datetime.y:257
		{
			// Special-case to handle "week" and "fortnight"
		yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval * yyS[yypt-0].intval)
		}
	case 61:
		//line datetime.y:261
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval * yyS[yypt-0].intval)
		}
	}
	goto yystack /* stack new state and value */
}
