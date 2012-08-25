
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

type textint struct {
    i, l int
    s string
}


//line datetime.y:24
type yySymType struct
{
	yys int
	tval    textint
	intval  int
	zoneval *time.Location
}

const T_DAYQUAL = 57346
const T_INTEGER = 57347
const T_PLUS = 57348
const T_MINUS = 57349
const T_AMPM = 57350
const T_MONTHNAME = 57351
const T_DAYNAME = 57352
const T_DAYS = 57353
const T_DAYSHIFT = 57354
const T_OFFSET = 57355
const T_ISOYD = 57356
const T_ISOHS = 57357
const T_RELATIVE = 57358
const T_AGO = 57359
const T_ZONE = 57360

var yyToknames = []string{
	"T_DAYQUAL",
	"T_INTEGER",
	"T_PLUS",
	"T_MINUS",
	"T_AMPM",
	"T_MONTHNAME",
	"T_DAYNAME",
	"T_DAYS",
	"T_DAYSHIFT",
	"T_OFFSET",
	"T_ISOYD",
	"T_ISOHS",
	"T_RELATIVE",
	"T_AGO",
	"T_ZONE",
}
var yyStatenames = []string{}

const yyEofCode = 1
const yyErrCode = 2
const yyMaxDepth = 200

//line datetime.y:423


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
	-1, 15,
	1, 87,
	5, 87,
	9, 8,
	12, 87,
	16, 87,
	20, 8,
	28, 87,
	-2, 12,
	-1, 81,
	5, 3,
	-2, 41,
}

const yyNprod = 88
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 161

var yyAct = []int{

	97, 92, 42, 73, 68, 58, 56, 33, 15, 23,
	24, 28, 16, 17, 37, 20, 40, 34, 29, 18,
	31, 48, 50, 7, 49, 51, 52, 87, 38, 57,
	67, 22, 30, 89, 32, 35, 36, 50, 53, 49,
	51, 52, 109, 89, 88, 75, 90, 98, 59, 81,
	59, 21, 65, 53, 110, 19, 90, 40, 69, 95,
	26, 85, 86, 100, 4, 76, 40, 69, 105, 38,
	74, 31, 43, 96, 40, 69, 54, 62, 38, 103,
	79, 63, 123, 104, 106, 128, 38, 40, 69, 124,
	94, 40, 69, 108, 119, 47, 112, 46, 118, 38,
	27, 23, 24, 38, 129, 120, 45, 44, 47, 107,
	46, 64, 99, 117, 125, 27, 23, 24, 101, 127,
	126, 122, 121, 116, 115, 130, 114, 113, 111, 93,
	125, 102, 80, 78, 77, 71, 70, 66, 61, 60,
	41, 82, 91, 84, 55, 25, 83, 14, 13, 12,
	11, 10, 9, 8, 6, 5, 72, 3, 2, 1,
	39,
}
var yyPact = []int{

	41, -1000, -1000, 3, 110, -1000, -1000, -1000, -1000, -15,
	-1000, -1000, -1000, -1000, -1000, 10, 135, 53, 97, 11,
	-1000, 59, 24, 134, 133, 95, -1000, -1000, 132, 85,
	131, -1000, 130, 50, 40, 129, 128, 70, -1000, 127,
	-1000, 137, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 22, -1000, 19, -1000, 124,
	-1000, -1000, -1000, 26, 84, -1000, 68, -1000, -1000, -1000,
	51, 23, 103, -1000, 42, 111, 126, -1000, 85, 50,
	46, 53, -1000, 100, -1000, -1000, -1000, 29, -1000, -1000,
	-1000, 124, -1000, 27, 123, 85, 122, -1000, 121, 119,
	-1000, 118, 106, -1000, 89, 117, 116, -1000, -1000, -1000,
	-1000, 60, -1000, 81, -1000, -1000, -1000, 115, 114, -1000,
	76, -1000, -1000, 99, 85, -1000, -1000, -1000, -1000, 85,
	-1000,
}
var yyPgo = []int{

	0, 160, 55, 4, 0, 159, 158, 157, 2, 3,
	156, 7, 155, 154, 23, 153, 152, 151, 150, 149,
	148, 147, 146, 51, 145, 144, 143, 6, 5, 142,
	1,
}
var yyR1 = []int{

	0, 5, 5, 8, 8, 9, 10, 10, 11, 11,
	1, 1, 2, 2, 2, 3, 3, 3, 4, 4,
	6, 7, 7, 12, 12, 12, 12, 12, 12, 12,
	12, 12, 13, 13, 13, 14, 14, 14, 15, 15,
	15, 15, 15, 15, 16, 16, 16, 16, 16, 17,
	17, 22, 18, 18, 18, 18, 18, 18, 18, 18,
	18, 18, 19, 19, 23, 23, 24, 24, 24, 24,
	24, 24, 24, 20, 20, 20, 25, 25, 28, 28,
	29, 29, 30, 30, 27, 26, 26, 21,
}
var yyR2 = []int{

	0, 1, 1, 0, 1, 2, 0, 1, 0, 1,
	1, 1, 1, 2, 2, 1, 2, 4, 0, 1,
	2, 0, 2, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 3, 5, 7, 2, 4, 6, 3, 5,
	4, 3, 5, 5, 3, 5, 3, 4, 6, 3,
	4, 0, 4, 2, 2, 2, 3, 5, 5, 6,
	6, 1, 1, 2, 1, 2, 2, 2, 2, 2,
	2, 2, 2, 3, 2, 3, 1, 2, 2, 2,
	1, 2, 2, 2, 2, 0, 1, 1,
}
var yyChk = []int{

	-1000, -5, -6, -7, 23, -12, -13, -14, -15, -16,
	-17, -18, -19, -20, -21, 5, 9, 10, 16, -2,
	12, -23, 28, 6, 7, -24, -2, 5, 26, 8,
	22, -3, 24, -11, 7, 25, 26, 4, 18, -1,
	6, 5, -8, 19, 10, 9, 13, 11, 10, 13,
	11, 14, 15, 27, 17, -25, -27, 5, -28, 26,
	5, 5, -23, -2, 16, -14, 5, -4, -3, 7,
	5, 5, -10, -9, 20, 5, 25, 5, 5, 10,
	5, -11, 4, -22, -26, -28, -27, 5, 25, 14,
	27, -29, -30, 5, 22, 8, 22, -4, 24, 9,
	21, 7, 5, -4, -9, 22, -8, 9, -30, 15,
	27, 5, -4, 5, 5, 5, 5, 7, 9, 5,
	16, 5, 5, 22, 8, -4, 5, 5, 9, 5,
	-4,
}
var yyDef = []int{

	21, -2, 1, 2, 0, 22, 23, 24, 25, 26,
	27, 28, 29, 30, 31, -2, 0, 3, 0, 0,
	61, 62, 0, 0, 0, 64, 20, 12, 0, 18,
	0, 35, 0, 6, 0, 0, 0, 9, 15, 0,
	10, 8, 51, 4, 53, 54, 67, 69, 55, 66,
	68, 70, 71, 72, 63, 85, 74, 0, 76, 0,
	13, 14, 65, 0, 0, 49, 0, 32, 19, 11,
	18, 38, 0, 7, 0, 44, 0, 46, 18, 56,
	16, -2, 9, 0, 73, 77, 86, 0, 75, 78,
	79, 84, 80, 0, 0, 18, 0, 36, 0, 40,
	5, 0, 47, 50, 0, 0, 0, 52, 81, 82,
	83, 18, 33, 18, 39, 42, 45, 0, 57, 58,
	0, 17, 43, 0, 18, 37, 48, 59, 60, 18,
	34,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 19, 3, 3, 24, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 22, 3,
	3, 3, 3, 3, 23, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 27, 3, 3,
	28, 3, 3, 3, 26, 3, 3, 25, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 21, 3, 3, 3, 3, 3, 3, 3,
	3, 20,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18,
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

	case 10:
		yyVAL.intval = yyS[yypt-0].intval
	case 11:
		yyVAL.intval = yyS[yypt-0].intval
	case 12:
		yyVAL.tval = yyS[yypt-0].tval
	case 13:
		//line datetime.y:61
		{
	        yyS[yypt-0].tval.s = "+" + yyS[yypt-0].tval.s
	        yyVAL.tval = yyS[yypt-0].tval
	    }
	case 14:
		//line datetime.y:65
		{
	        yyS[yypt-0].tval.s = "-" + yyS[yypt-0].tval.s
	        yyS[yypt-0].tval.i *= -1
	        yyVAL.tval = yyS[yypt-0].tval
	    }
	case 15:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 16:
		//line datetime.y:73
		{
	        hrs, mins := yyS[yypt-0].tval.i, 0
	        if (yyS[yypt-0].tval.l == 4) {
	            hrs, mins = (yyS[yypt-0].tval.i / 100), (yyS[yypt-0].tval.i % 100)
	        } else if (yyS[yypt-0].tval.l == 2) {
	            hrs *= 100
	        } else {
	            yylex.Error("Invalid timezone offset " +yyS[yypt-0].tval.s)
	        }
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-1].intval * (3600 * hrs + 60 * mins))
		}
	case 17:
		//line datetime.y:84
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-3].intval * (3600 * yyS[yypt-2].tval.i + 60 * yyS[yypt-0].tval.i))
		}
	case 18:
		//line datetime.y:89
		{ yyVAL.zoneval = nil }
	case 19:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 20:
		//line datetime.y:93
		{
			yylex.(*dateLexer).time = time.Unix(int64(yyS[yypt-0].tval.i), 0)
		}
	case 32:
		//line datetime.y:115
		{
			yylex.(*dateLexer).setTime(yyS[yypt-2].tval.i + yyS[yypt-1].intval, 0, 0, yyS[yypt-0].zoneval)
		}
	case 33:
		//line datetime.y:118
		{
			yylex.(*dateLexer).setTime(yyS[yypt-4].tval.i + yyS[yypt-1].intval, yyS[yypt-2].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 34:
		//line datetime.y:121
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i + yyS[yypt-1].intval, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].zoneval)
		}
	case 35:
		//line datetime.y:128
		{
	        yylex.(*dateLexer).setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
	    }
	case 36:
		//line datetime.y:131
		{
	        yylex.(*dateLexer).setTime(yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, 0, yyS[yypt-0].zoneval)
	    }
	case 37:
		//line datetime.y:134
		{
	        yylex.(*dateLexer).setTime(yyS[yypt-5].tval.i, yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, yyS[yypt-0].zoneval)
	    }
	case 38:
		//line datetime.y:142
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume we have MM/YYYY
			l.setDate(yyS[yypt-0].tval.i, yyS[yypt-2].tval.i, 1)
			} else {
	            // assume we have DD/MM (too bad, americans)
            l.setDate(0, yyS[yypt-0].tval.i, yyS[yypt-2].tval.i)
			}
		}
	case 39:
		//line datetime.y:152
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume we have DD/MM/YYYY
			l.setDate(yyS[yypt-0].tval.i, yyS[yypt-2].tval.i, yyS[yypt-4].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume we have DD/MM/YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].tval.i + 1900, yyS[yypt-2].tval.i, yyS[yypt-4].tval.i)
			} else {
				// assume we have DD/MM/YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].tval.i + 2000, yyS[yypt-2].tval.i, yyS[yypt-4].tval.i)
			}
		}
	case 40:
		//line datetime.y:165
		{
			// DDth of Mon
		yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 41:
		//line datetime.y:169
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-1].tval.l == 4 {
				// assume Mon YYYY
			l.setDate(yyS[yypt-1].tval.i, yyS[yypt-2].intval, 1)
			} else {
			    // assume Mon DDth
			l.setDate(0, yyS[yypt-2].intval, yyS[yypt-1].tval.i)
			}
	    }
	case 42:
		//line datetime.y:179
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume DDth of Mon YYYY
			l.setDate(yyS[yypt-0].tval.i, yyS[yypt-1].intval, yyS[yypt-4].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume DDth of Mon YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].tval.i + 1900, yyS[yypt-1].intval, yyS[yypt-4].tval.i)
			} else {
				// assume DDth of Mon YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].tval.i + 2000, yyS[yypt-1].intval, yyS[yypt-4].tval.i)
			}
		}
	case 43:
		//line datetime.y:192
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume Mon DDth, YYYY
			l.setDate(yyS[yypt-0].tval.i, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume Mon DDth, YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].tval.i + 1900, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			} else {
				// assume Mon DDth YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].tval.i + 2000, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			}
		}
	case 44:
		//line datetime.y:208
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-2].tval.l == 4 && yyS[yypt-0].tval.l == 3 {
	            // assume we have YYYY-DDD
            l.setDate(yyS[yypt-2].tval.i, 1, yyS[yypt-0].tval.i)
	        } else if yyS[yypt-2].tval.l == 4 {
				// assume we have YYYY-MM
			l.setDate(yyS[yypt-2].tval.i, yyS[yypt-0].tval.i, 1)
			} else {
	            // assume we have MM-DD (not strictly ISO compliant)
            // this is for americans, because of DD/MM above ;-)
            l.setDate(0, yyS[yypt-0].tval.i, yyS[yypt-2].tval.i)
			}
		}
	case 45:
		//line datetime.y:222
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-4].tval.l == 4 {
				// assume we have YYYY-MM-DD
			l.setDate(yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
			} else if yyS[yypt-4].tval.i > 68 {
				// assume we have YY-MM-DD, add 1900 if YY > 68
			l.setDate(yyS[yypt-4].tval.i + 1900, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
			} else {
				// assume we have YY-MM-DD, add 2000 otherwise
			l.setDate(yyS[yypt-4].tval.i + 2000, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
			}
		}
	case 46:
		//line datetime.y:235
		{
	        l := yylex.(*dateLexer)
	        wday, week := 1, yyS[yypt-0].tval.i
	        if yyS[yypt-0].tval.l == 3 {
	            // assume YYYY'W'WWD
            week = week / 10
	            wday = week % 10
	        }
	        l.setWeek(yyS[yypt-2].tval.i, week, wday)
	    }
	case 47:
		//line datetime.y:245
		{
	        // assume YYYY-'W'WW
        yylex.(*dateLexer).setWeek(yyS[yypt-3].tval.i, yyS[yypt-0].tval.i, 1)
	    }
	case 48:
		//line datetime.y:249
		{
	        // assume YYYY-'W'WW-D
        yylex.(*dateLexer).setWeek(yyS[yypt-5].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
	    }
	case 50:
		//line datetime.y:257
		{
	        // this goes here because the YYYYMMDD and HHMMSS forms of the
        // ISO 8601 format date and time are handled by 'integer' below.
        l := yylex.(*dateLexer)
	        l.setYMD(yyS[yypt-3].tval.i, yyS[yypt-3].tval.l)
	        l.setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
	    }
	case 51:
		//line datetime.y:266
		{
			// Tuesday,
		yylex.(*dateLexer).setDay(yyS[yypt-1].intval, 1)
		}
	case 52:
		//line datetime.y:270
		{
			// March
		yylex.(*dateLexer).setMonth(yyS[yypt-3].intval, 1)
		}
	case 53:
		//line datetime.y:274
		{
			// Next tuesday
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 54:
		//line datetime.y:278
		{
			// Next march
		yylex.(*dateLexer).setMonth(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 55:
		//line datetime.y:282
		{
			// +-N Tuesdays
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-1].tval.i)
		}
	case 56:
		//line datetime.y:286
		{
			// 3rd Tuesday 
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-2].tval.i)
		}
	case 57:
		//line datetime.y:290
		{
			// 3rd Tuesday of (implicit this) March
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			l.setMonth(yyS[yypt-0].intval, 1)
		}
	case 58:
		//line datetime.y:296
		{
			// 3rd Tuesday of 2012
		yylex.(*dateLexer).setDay(yyS[yypt-2].intval, yyS[yypt-4].tval.i, yyS[yypt-0].tval.i)
		}
	case 59:
		//line datetime.y:300
		{
			// 3rd Tuesday of March 2012
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonth(yyS[yypt-1].intval, 1, yyS[yypt-0].tval.i)
		}
	case 60:
		//line datetime.y:306
		{
			// 3rd Tuesday of next March
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonth(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 61:
		//line datetime.y:312
		{
			// yesterday or tomorrow
		d := time.Now().Weekday()
			yylex.(*dateLexer).setDay((7+int(d)+yyS[yypt-0].intval)%7, yyS[yypt-0].intval)
		}
	case 63:
		//line datetime.y:320
		{
			yylex.(*dateLexer).setAgo()
		}
	case 66:
		//line datetime.y:329
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 67:
		//line datetime.y:332
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 68:
		//line datetime.y:335
		{
			// Special-case to handle "week" and "fortnight"
		yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].tval.i * yyS[yypt-0].intval)
		}
	case 69:
		//line datetime.y:339
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval * yyS[yypt-0].intval)
		}
	case 70:
		//line datetime.y:342
		{
	        // As we need to be able to separate out YD from HS in ISO durations
        // this becomes a fair bit messier than if Y D H S were just T_OFFSET
        // Because writing "next y" or "two h" would be odd, disallow
        // T_RELATIVE tokens from being used with ISO single-letter notation
        yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
	    }
	case 71:
		//line datetime.y:349
		{
	        yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
	    }
	case 72:
		//line datetime.y:352
		{
	        // Resolve 'm' ambiguity in favour of minutes outside ISO duration
        yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
	    }
	case 75:
		//line datetime.y:361
		{
	        yylex.(*dateLexer).addOffset(O_DAY, 7 * yyS[yypt-1].tval.i)
	    }
	case 78:
		//line datetime.y:371
		{
	        // takes care of Y and D
        yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
	    }
	case 79:
		//line datetime.y:375
		{
	        yylex.(*dateLexer).addOffset(O_MONTH, yyS[yypt-1].tval.i)
	    }
	case 82:
		//line datetime.y:384
		{
	        // takes care of H and S
        yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
	    }
	case 83:
		//line datetime.y:388
		{
	        yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
	    }
	case 87:
		//line datetime.y:400
		{
	        l := yylex.(*dateLexer)
	        if yyS[yypt-0].tval.l == 8 {
	            // assume ISO 8601 YYYYMMDD
            l.setYMD(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l)
	        } else {
	            // assume ISO 8601 HHMMSS with no zone
            l.setHMS(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l, nil)
	        }   
	    }
	}
	goto yystack /* stack new state and value */
}
