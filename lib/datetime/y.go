
//line datetime.y:2
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


//line datetime.y:22
type	yySymType	struct
{
	yys	int;
  strval    string
  intval	int
}
const	T_AMPM	= 57346
const	T_PLUS	= 57347
const	T_MINUS	= 57348
const	T_MONTH	= 57349
const	T_INTEGER	= 57350
var	yyToknames	 =[]string {
	"T_AMPM",
	"T_PLUS",
	"T_MINUS",
	"T_MONTH",
	"T_INTEGER",
}
var	yyStatenames	 =[]string {
}
																										const	yyEofCode	= 1
const	yyErrCode	= 2
const	yyMaxDepth	= 200

//line datetime.y:350


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

//line yacctab:1
var	yyExca = []int {
-1, 1,
	1, -1,
	-2, 0,
-1, 31,
	8, 9,
	-2, 33,
}
const	yyNprod	= 36
const	yyPrivate	= 57344
var	yyTokenNames []string
var	yyStates []string
const	yyLast	= 61
var	yyAct	= []int {

  32,  35,  14,  16,  15,  30,  29,  28,  13,  17,
  18,  21,  19,  27,  20,  18,  21,  19,  49,  20,
  16,  42,  34,  37,  38,   4,  17,  31,  39,  33,
   9,   8,  50,  34,  37,  38,  44,   7,  47,  46,
  45,  43,  40,  25,  48,  24,  23,  22,  37,  38,
   6,  26,  11,  12,   5,  41,   3,   2,   1,  10,
  36,
};
var	yyPact	= []int {

  16,-1000,-1000,  23,  47,-1000,-1000,-1000,  -3,  39,
  38,-1000,-1000,  37,  35,  44,-1000,-1000,  -1,  -9,
 -10, -13,   2,-1000,  18,  14,  34,-1000,-1000,-1000,
-1000,  11,-1000,  33,  43,-1000,  32,-1000,-1000,  31,
-1000,  30,-1000,  29,-1000,   7,-1000,-1000,-1000,  24,
-1000,
};
var	yyPgo	= []int {

   0,  60,   1,  59,   0,  58,  57,  56,  55,  54,
  50,  37,   2,   4,
};
var	yyR1	= []int {

   0,   5,   5,   6,   3,   3,   3,   1,   1,   8,
   8,   7,   7,   9,   9,  10,  10,   4,   4,   4,
   4,   2,   2,  12,  12,  13,  13,  13,  13,  13,
  11,  11,  11,  11,  11,  11,
};
var	yyR2	= []int {

   0,   1,   1,   3,   0,   1,   1,   1,   1,   0,
   1,   0,   2,   1,   1,   4,   6,   0,   1,   2,
   1,   2,   4,   1,   1,   0,   2,   2,   2,   2,
   3,   5,   3,   3,   4,   5,
};
var	yyChk	= []int {

-1000,  -5,  -6,  -7,   9,  -9, -10, -11,   8,   7,
  -3,   5,   6,  11, -12, -13,   6,  12,  13,  15,
  17,  14,   8,   8,   8,   8,   7,  14,  16,  16,
  18, -13,  -4,  11,   4,  -2,  -1,   5,   6, -12,
   8,  -8,  10,   8,  -2,   8,   8,   8,  -4,  11,
   8,
};
var	yyDef	= []int {

  11,  -2,   1,   2,   4,  12,  13,  14,  25,   0,
   0,   5,   6,   0,   0,   0,  23,  24,   0,   0,
   0,   0,  25,   3,  17,  30,  32,  26,  27,  28,
  29,  -2,  15,   0,  18,  20,   0,   7,   8,   0,
  34,   0,  10,  17,  19,  21,  31,  35,  16,   0,
  22,
};
var	yyTok1	= []int {

   1,   3,   3,   3,   3,   3,   3,   3,   3,   3,
   3,   3,   3,   3,   3,   3,   3,   3,   3,   3,
   3,   3,   3,   3,   3,   3,   3,   3,   3,   3,
   3,   3,   3,   3,   3,   3,   3,   3,   3,   3,
   3,   3,   3,   3,  10,   3,   3,  12,   3,   3,
   3,   3,   3,   3,   3,   3,   3,   3,  11,   3,
   3,   3,   3,   3,   9,   3,   3,   3,   3,   3,
   3,   3,   3,   3,   3,   3,   3,   3,   3,   3,
   3,   3,   3,   3,   3,   3,   3,   3,   3,   3,
   3,   3,   3,   3,   3,   3,   3,   3,   3,   3,
  16,   3,   3,   3,  18,   3,   3,   3,   3,   3,
  15,   3,   3,   3,  17,  13,  14,
};
var	yyTok2	= []int {

   2,   3,   4,   5,   6,   7,   8,
};
var	yyTok3	= []int {
   0,
 };

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

				/* the current p has no shift onn "error", pop stack */
				if yyDebug >= 2 {
					fmt.Printf("error recovery pops state %d, uncovers %d\n",
						yyS[yyp].yys, yyS[yyp-1].yys)
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
	_ = yypt		// guard against "declared and not used"

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
//line datetime.y:53
{
		if yyS[yypt-1].strval == "-" {
			yylex.(*dateLexer).time = time.SecondsToUTC(int64(-yyS[yypt-0].intval))
		} else {
			yylex.(*dateLexer).time = time.SecondsToUTC(int64(yyS[yypt-0].intval))
		}
	}
case 4:
//line datetime.y:62
{ yyVAL.strval = "" }
case 5:
	yyVAL.strval = yyS[yypt-0].strval;
case 6:
	yyVAL.strval = yyS[yypt-0].strval;
case 7:
	yyVAL.strval = yyS[yypt-0].strval;
case 8:
	yyVAL.strval = yyS[yypt-0].strval;
case 15:
//line datetime.y:95
{
		l := yylex.(*dateLexer)
		l.parseTime(
			fmt.Sprintf("%s:04%s%s", l.hourfmt, l.ampmfmt, l.zonefmt),
			fmt.Sprintf("%02d:%02d%s", yyS[yypt-3].intval, yyS[yypt-1].intval, yyS[yypt-0].strval)) 
	}
case 16:
//line datetime.y:101
{
		l := yylex.(*dateLexer)
		l.parseTime(
			fmt.Sprintf("%s:04:05%s%s", l.hourfmt, l.ampmfmt, l.zonefmt),
			fmt.Sprintf("%d:%02d:%02d%s", yyS[yypt-5].intval, yyS[yypt-3].intval, yyS[yypt-1].intval, yyS[yypt-0].strval))
	}
case 17:
//line datetime.y:109
{
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt, l.zonefmt = "15", "", ""
		yyVAL.strval = ""
	}
case 18:
//line datetime.y:114
{
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt, l.zonefmt = "3", yyS[yypt-0].strval, ""
	}
case 19:
//line datetime.y:118
{
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt = "3", yyS[yypt-1].strval
		yyVAL.strval = fmt.Sprintf("%s%s", yyS[yypt-1].strval, yyS[yypt-0].strval)
	}
case 20:
//line datetime.y:123
{
		l := yylex.(*dateLexer)
		l.hourfmt, l.ampmfmt = "15", ""
	}
case 21:
//line datetime.y:129
{
		l := yylex.(*dateLexer)
		l.zonefmt = "-0700"
		yyVAL.strval = fmt.Sprintf("%s%04d", yyS[yypt-1].strval, yyS[yypt-0].intval)
	}
case 22:
//line datetime.y:134
{
		l := yylex.(*dateLexer)
		l.zonefmt = "-07:00"
		yyVAL.strval = fmt.Sprintf("%s%02d:%02d", yyS[yypt-3].strval, yyS[yypt-2].intval, yyS[yypt-0].intval)
	}
case 30:
//line datetime.y:210
{
		// DD-MM or MM-YYYY or YYYY-MM
		l := yylex.(*dateLexer)
		if yyS[yypt-0].intval > 12 {
			l.parseDate("1 2006", fmt.Sprintf("%d %04d", yyS[yypt-2].intval, yyS[yypt-0].intval))
		} else if yyS[yypt-2].intval > 31 {
			l.parseDate("2006 1", fmt.Sprintf("%04d %d", yyS[yypt-2].intval, yyS[yypt-0].intval))
		} else {
			l.parseDate("2 1", fmt.Sprintf("%d %d", yyS[yypt-2].intval, yyS[yypt-0].intval))
		}
	}
case 31:
//line datetime.y:221
{
		// YYYY-MM-DD or DD-MM-YY(YY?).
		l := yylex.(*dateLexer)
		if yyS[yypt-4].intval > 31 {
			l.parseDate("2006 1 2", fmt.Sprintf("%04d %d %d", yyS[yypt-4].intval, yyS[yypt-2].intval, yyS[yypt-0].intval))
		} else if yyS[yypt-2].intval > 99 {
			l.parseDate("2 1 2006", fmt.Sprintf("%d %d %04d", yyS[yypt-4].intval, yyS[yypt-2].intval, yyS[yypt-0].intval))
		} else {
			l.parseDate("2 1 06", fmt.Sprintf("%d %d %02d", yyS[yypt-4].intval, yyS[yypt-2].intval, yyS[yypt-0].intval))
		}
	}
case 32:
//line datetime.y:232
{
		// 15th feb
		l := yylex.(*dateLexer)
		l.parseDate("2 Jan", fmt.Sprintf("%d %s", yyS[yypt-2].intval, yyS[yypt-0].strval))
	}
case 33:
//line datetime.y:237
{
		// feb 15th or feb 2010
		l := yylex.(*dateLexer)
		if yyS[yypt-1].intval > 31 {
			l.parseDate("Jan 2006", fmt.Sprintf("%s %04d", yyS[yypt-2].strval, yyS[yypt-1].intval))
		} else {
			l.parseDate("2 Jan", fmt.Sprintf("%d %s", yyS[yypt-1].intval, yyS[yypt-2].strval))
		}
	}
case 34:
//line datetime.y:246
{
		// 15th feb 2010
		l := yylex.(*dateLexer)
		l.parseDate("2 Jan 2006", fmt.Sprintf("%d %s %04d", yyS[yypt-3].intval, yyS[yypt-1].strval, yyS[yypt-0].intval))
	}
case 35:
//line datetime.y:251
{
		// feb 15th 2010
		l := yylex.(*dateLexer)
		l.parseDate("Jan 2 2006", fmt.Sprintf("%s %d %04d", yyS[yypt-4].strval, yyS[yypt-3].intval, yyS[yypt-0].intval))
	}
	}
	goto yystack /* stack new state and value */
}
