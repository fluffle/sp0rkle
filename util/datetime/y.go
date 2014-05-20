//line datetime.y:2
package datetime

import __yyfmt__ "fmt"

//line datetime.y:2
// Based upon parse-datetime.y in GNU coreutils.
// also an exercise in learning goyacc in particular.
// This file contains the yacc grammar only.
// See lexer.go for the lexer and parse functions,
// and tokenmaps.go for the token maps.

import (
	"time"
)

type textint struct {
	i, l int
	s    string
}

//line datetime.y:21
type yySymType struct {
	yys     int
	tval    textint
	intval  int
	zoneval *time.Location
}

const T_OF = 57346
const T_THE = 57347
const T_IGNORE = 57348
const T_DAYQUAL = 57349
const T_INTEGER = 57350
const T_PLUS = 57351
const T_MINUS = 57352
const T_MONTHNAME = 57353
const T_DAYNAME = 57354
const T_DAYS = 57355
const T_DAYSHIFT = 57356
const T_OFFSET = 57357
const T_ISOYD = 57358
const T_ISOHS = 57359
const T_RELATIVE = 57360
const T_AGO = 57361
const T_ZONE = 57362

var yyToknames = []string{
	"T_OF",
	"T_THE",
	"T_IGNORE",
	"T_DAYQUAL",
	"T_INTEGER",
	"T_PLUS",
	"T_MINUS",
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

//line datetime.y:502

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 16,
	4, 10,
	11, 10,
	12, 22,
	13, 22,
	15, 22,
	16, 22,
	17, 22,
	23, 22,
	-2, 105,
	-1, 18,
	8, 8,
	-2, 67,
	-1, 115,
	8, 4,
	-2, 54,
	-1, 144,
	8, 4,
	-2, 52,
}

const yyNprod = 106
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 198

var yyAct = []int{

	113, 51, 32, 33, 52, 108, 36, 69, 81, 67,
	105, 7, 31, 40, 114, 47, 37, 106, 88, 4,
	142, 87, 103, 104, 123, 34, 45, 141, 41, 22,
	44, 42, 43, 80, 35, 38, 39, 58, 60, 89,
	59, 61, 62, 78, 70, 122, 17, 15, 63, 16,
	24, 25, 18, 19, 64, 28, 73, 21, 68, 20,
	95, 96, 29, 27, 53, 60, 23, 59, 61, 62,
	93, 94, 121, 65, 101, 63, 102, 92, 47, 82,
	70, 64, 152, 110, 74, 144, 111, 112, 34, 45,
	128, 41, 119, 44, 42, 43, 129, 105, 115, 47,
	82, 57, 137, 56, 106, 117, 125, 47, 82, 77,
	45, 76, 131, 127, 44, 116, 43, 134, 45, 139,
	41, 159, 138, 42, 158, 155, 30, 24, 25, 140,
	145, 154, 28, 148, 146, 147, 75, 50, 47, 82,
	27, 55, 54, 57, 151, 56, 153, 150, 156, 45,
	30, 24, 25, 149, 143, 136, 148, 135, 133, 132,
	130, 109, 118, 99, 97, 91, 90, 84, 83, 79,
	72, 71, 48, 126, 98, 124, 120, 86, 107, 100,
	66, 26, 14, 13, 12, 11, 10, 9, 8, 6,
	5, 49, 85, 3, 2, 1, 157, 46,
}
var yyPact = []int{

	-8, -1000, -1000, 41, 142, -1000, -1000, -1000, -1000, -18,
	-1000, -1000, -1000, -1000, -1000, -1000, 6, 164, 132, 43,
	130, 25, 54, 50, 163, 162, 118, 96, -1000, -1000,
	-1000, 161, 129, 160, -1000, 159, 173, 10, 158, 157,
	65, 47, 37, -1000, -1000, -1000, 156, -1000, 167, 155,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, 14, -1000, -6, -1000,
	153, -1000, -1000, -1000, 52, 88, -1000, -1000, -1000, 90,
	-1000, -1000, -1000, 69, -14, 87, -1000, 105, 95, 154,
	-1000, 129, 172, -1000, 49, -1000, 22, -2, 171, 166,
	-1000, -1000, -1000, 81, -1000, -1000, -1000, 153, -1000, 73,
	152, 129, 151, -1000, 150, 43, 149, 147, 92, -1000,
	111, 3, -4, 146, 74, 43, -1000, -1000, -1000, -1000,
	90, -1000, 98, -1000, 145, -1000, -1000, 139, 136, -1000,
	71, -1000, -1000, -1000, 43, 123, 117, 129, 116, -1000,
	-1000, -1000, -1000, 113, -1000, 129, -1000, -1000, -1000, -1000,
}
var yyPgo = []int{

	0, 197, 2, 57, 196, 8, 0, 195, 194, 193,
	4, 1, 192, 191, 6, 3, 190, 189, 11, 188,
	187, 186, 185, 184, 183, 182, 29, 181, 180, 179,
	9, 7, 178, 5,
}
var yyR1 = []int{

	0, 7, 7, 10, 11, 11, 12, 12, 13, 13,
	14, 14, 4, 4, 2, 2, 2, 2, 15, 15,
	1, 1, 3, 3, 3, 5, 5, 5, 6, 6,
	8, 9, 9, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 17, 17, 17, 18, 18, 18, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 20,
	20, 20, 20, 20, 21, 21, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 23, 23, 26, 26,
	27, 27, 27, 27, 27, 27, 27, 27, 27, 27,
	27, 24, 24, 24, 28, 28, 31, 31, 32, 32,
	33, 33, 30, 29, 29, 25,
}
var yyR2 = []int{

	0, 1, 1, 1, 0, 1, 0, 1, 0, 1,
	0, 1, 0, 1, 2, 4, 2, 4, 1, 1,
	1, 1, 1, 2, 2, 1, 2, 4, 0, 1,
	2, 0, 2, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 5, 7, 2, 4, 7, 3,
	5, 3, 5, 7, 4, 4, 6, 5, 6, 3,
	5, 3, 4, 6, 3, 4, 2, 1, 2, 2,
	2, 3, 5, 5, 6, 6, 1, 2, 1, 2,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	1, 3, 2, 3, 1, 2, 2, 2, 1, 2,
	2, 2, 2, 0, 1, 1,
}
var yyChk = []int{

	-1000, -7, -8, -9, 27, -16, -17, -18, -19, -20,
	-21, -22, -23, -24, -25, 6, 8, 5, 11, 12,
	18, -3, -26, 25, 9, 10, -27, 22, 14, -3,
	8, 30, -2, -15, -5, 28, -14, 10, 29, 30,
	7, 22, 25, 26, 24, 20, -1, 9, 8, -13,
	5, -11, -10, 21, 12, 11, 15, 13, 12, 15,
	13, 16, 17, 23, 29, 19, -28, -30, 8, -31,
	30, 8, 8, -26, -3, 18, 15, 13, -18, 8,
	-6, -5, 10, 8, 8, -12, 4, 11, 8, 29,
	8, 8, 12, 23, 24, 23, 24, 8, 7, 8,
	-29, -31, -30, 8, 29, 16, 23, -32, -33, 8,
	-15, -2, -15, -6, 28, 11, 10, 10, 8, -6,
	4, 23, 23, 26, 4, -14, 7, -33, 17, 23,
	8, -6, 8, 8, -11, 8, 8, 10, 11, 8,
	18, 24, 24, 8, 11, -10, -15, -2, -6, 8,
	8, 8, 11, -11, 8, 8, -6, -4, 8, 8,
}
var yyDef = []int{

	31, -2, 1, 2, 0, 32, 33, 34, 35, 36,
	37, 38, 39, 40, 41, 42, -2, 0, -2, 4,
	0, 0, 76, 0, 0, 0, 78, 0, 90, 30,
	22, 0, 28, 0, 46, 0, 6, 0, 0, 0,
	11, 0, 0, 18, 19, 25, 0, 20, 0, 0,
	9, 66, 5, 3, 68, 69, 81, 84, 70, 80,
	83, 86, 87, 88, 89, 77, 103, 92, 0, 94,
	0, 23, 24, 79, 0, 0, 82, 85, 64, 0,
	43, 29, 21, 28, 49, 0, 7, 0, 59, 0,
	61, 28, 71, 14, 0, 16, 0, 26, 51, 10,
	91, 95, 104, 0, 93, 96, 97, 102, 98, 0,
	0, 28, 0, 47, 0, -2, 0, 0, 62, 65,
	0, 0, 0, 0, 0, 55, 11, 99, 100, 101,
	28, 44, 28, 50, 0, 57, 60, 0, 72, 73,
	0, 15, 17, 27, -2, 0, 0, 28, 12, 56,
	63, 74, 75, 0, 58, 28, 45, 48, 13, 53,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 21, 3, 24, 28, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 26, 3,
	3, 3, 3, 3, 27, 22, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 23, 3, 3,
	25, 3, 3, 3, 30, 3, 3, 29,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20,
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
	// 4 is TOKSTART above
	if c >= 4 && c-4 < len(yyToknames) {
		if yyToknames[c-4] != "" {
			return yyToknames[c-4]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
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
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(c), uint(char))
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
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yychar), yyStatname(yystate))
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
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yychar))
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
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yychar))
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
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
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

	case 12:
		//line datetime.y:54
		{
			yyVAL.tval = textint{}
		}
	case 13:
		yyVAL.tval = yyS[yypt-0].tval
	case 14:
		//line datetime.y:57
		{
			yyVAL.intval = 0
		}
	case 15:
		//line datetime.y:60
		{
			yyVAL.intval = 0
		}
	case 16:
		//line datetime.y:63
		{
			yyVAL.intval = 12
		}
	case 17:
		//line datetime.y:66
		{
			yyVAL.intval = 12
		}
	case 20:
		yyVAL.intval = yyS[yypt-0].intval
	case 21:
		yyVAL.intval = yyS[yypt-0].intval
	case 22:
		yyVAL.tval = yyS[yypt-0].tval
	case 23:
		//line datetime.y:77
		{
			yyS[yypt-0].tval.s = "+" + yyS[yypt-0].tval.s
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 24:
		//line datetime.y:81
		{
			yyS[yypt-0].tval.s = "-" + yyS[yypt-0].tval.s
			yyS[yypt-0].tval.i *= -1
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 25:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 26:
		//line datetime.y:89
		{
			hrs, mins := yyS[yypt-0].tval.i, 0
			if yyS[yypt-0].tval.l == 4 {
				hrs, mins = (yyS[yypt-0].tval.i / 100), (yyS[yypt-0].tval.i % 100)
			} else if yyS[yypt-0].tval.l == 2 {
				hrs *= 100
			} else {
				yylex.Error("Invalid timezone offset " + yyS[yypt-0].tval.s)
			}
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-1].intval*(3600*hrs+60*mins))
		}
	case 27:
		//line datetime.y:100
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-3].intval*(3600*yyS[yypt-2].tval.i+60*yyS[yypt-0].tval.i))
		}
	case 28:
		//line datetime.y:105
		{
			yyVAL.zoneval = nil
		}
	case 29:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 30:
		//line datetime.y:109
		{
			l := yylex.(*dateLexer)
			if !l.state(HAVE_TIME, true) {
				l.time = time.Unix(int64(yyS[yypt-0].tval.i), 0)
			}
		}
	case 43:
		//line datetime.y:135
		{
			l := yylex.(*dateLexer)
			// Hack to allow HHMMam to parse correctly, cos adie is a mong.
			if yyS[yypt-2].tval.l == 3 || yyS[yypt-2].tval.l == 4 {
				l.setTime(ampm(yyS[yypt-2].tval.i/100, yyS[yypt-1].intval), yyS[yypt-2].tval.i%100, 0, yyS[yypt-0].zoneval)
			} else {
				l.setTime(ampm(yyS[yypt-2].tval.i, yyS[yypt-1].intval), 0, 0, yyS[yypt-0].zoneval)
			}
		}
	case 44:
		//line datetime.y:144
		{
			yylex.(*dateLexer).setTime(yyS[yypt-4].tval.i+yyS[yypt-1].intval, yyS[yypt-2].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 45:
		//line datetime.y:147
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i+yyS[yypt-1].intval, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].zoneval)
		}
	case 46:
		//line datetime.y:154
		{
			yylex.(*dateLexer).setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
		}
	case 47:
		//line datetime.y:157
		{
			yylex.(*dateLexer).setTime(yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 48:
		//line datetime.y:160
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-1].zoneval)
			// Hack to make time.ANSIC, time.UnixDate and time.RubyDate parse
			if yyS[yypt-0].tval.l == 4 {
				yylex.(*dateLexer).setYear(yyS[yypt-0].tval.i)
			}
		}
	case 49:
		//line datetime.y:172
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
	case 50:
		//line datetime.y:182
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume we have DD/MM/YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-2].tval.i, yyS[yypt-4].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume we have DD/MM/YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-2].tval.i, yyS[yypt-4].tval.i)
			} else {
				// assume we have DD/MM/YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-2].tval.i, yyS[yypt-4].tval.i)
			}
		}
	case 51:
		//line datetime.y:195
		{
			// the DDth
			yylex.(*dateLexer).setDay(yyS[yypt-1].tval.i)
		}
	case 52:
		//line datetime.y:199
		{
			// the DDth of Month
			yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 53:
		//line datetime.y:203
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// the DDth of Month[,] YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-2].intval, yyS[yypt-5].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// the DDth of Month[,] YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-2].intval, yyS[yypt-5].tval.i)
			} else {
				// the DDth of Month[,] YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-2].intval, yyS[yypt-5].tval.i)
			}
		}
	case 54:
		//line datetime.y:216
		{
			// DD[th] [of] Month
			yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 55:
		//line datetime.y:220
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-1].tval.l == 4 {
				// assume Month YYYY
				l.setDate(yyS[yypt-1].tval.i, yyS[yypt-3].intval, 1)
			} else {
				// assume Month [the] DD[th]
				l.setDate(0, yyS[yypt-3].intval, yyS[yypt-1].tval.i)
			}
		}
	case 56:
		//line datetime.y:230
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume DD[th] [of] Month[,] YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-2].intval, yyS[yypt-5].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume DD[th] [of] Month[,] YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-2].intval, yyS[yypt-5].tval.i)
			} else {
				// assume DD[th] [of] Month[,] YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-2].intval, yyS[yypt-5].tval.i)
			}
		}
	case 57:
		//line datetime.y:243
		{
			// RFC 850, srsly :(
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume DD-Mon-YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume DD-Mon-YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			} else {
				// assume DD-Mon-YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			}
		}
	case 58:
		//line datetime.y:257
		{
			// comma cannot be optional here; T_MONTHNAME T_INTEGER T_INTEGER
			// can easily match [March 02 10]:30:00 and break parsing.
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume Month [the] DD[th], YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-5].intval, yyS[yypt-3].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume Month [the] DD[th], YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-5].intval, yyS[yypt-3].tval.i)
			} else {
				// assume Month [the] DD[th], YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-5].intval, yyS[yypt-3].tval.i)
			}
		}
	case 59:
		//line datetime.y:275
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
				l.setDate(0, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
			}
		}
	case 60:
		//line datetime.y:289
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-4].tval.l == 4 {
				// assume we have YYYY-MM-DD
				l.setDate(yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
			} else if yyS[yypt-4].tval.i > 68 {
				// assume we have YY-MM-DD, add 1900 if YY > 68
				l.setDate(yyS[yypt-4].tval.i+1900, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
			} else {
				// assume we have YY-MM-DD, add 2000 otherwise
				l.setDate(yyS[yypt-4].tval.i+2000, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
			}
		}
	case 61:
		//line datetime.y:302
		{
			l := yylex.(*dateLexer)
			wday, week := 1, yyS[yypt-0].tval.i
			if yyS[yypt-0].tval.l == 3 {
				// assume YYYY'W'WWD
				wday = week % 10
				week = week / 10
			}
			l.setWeek(yyS[yypt-2].tval.i, week, wday)
		}
	case 62:
		//line datetime.y:312
		{
			// assume YYYY-'W'WW
			yylex.(*dateLexer).setWeek(yyS[yypt-3].tval.i, yyS[yypt-0].tval.i, 1)
		}
	case 63:
		//line datetime.y:316
		{
			// assume YYYY-'W'WW-D
			yylex.(*dateLexer).setWeek(yyS[yypt-5].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
		}
	case 65:
		//line datetime.y:324
		{
			// this goes here because the YYYYMMDD and HHMMSS forms of the
			// ISO 8601 format date and time are handled by 'integer' below.
			l := yylex.(*dateLexer)
			if yyS[yypt-3].tval.l == 8 {
				// assume ISO 8601 YYYYMMDD
				l.setYMD(yyS[yypt-3].tval.i, yyS[yypt-3].tval.l)
			} else if yyS[yypt-3].tval.l == 7 {
				// assume ISO 8601 ordinal YYYYDDD
				l.setDate(yyS[yypt-3].tval.i/1000, 1, yyS[yypt-3].tval.i%1000)
			}
			l.setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
		}
	case 66:
		//line datetime.y:339
		{
			// Tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-1].intval, 0)
		}
	case 67:
		//line datetime.y:343
		{
			// March
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, 0)
		}
	case 68:
		//line datetime.y:347
		{
			// Next tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 69:
		//line datetime.y:351
		{
			// Next march
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 70:
		//line datetime.y:355
		{
			// +-N Tuesdays
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].tval.i)
		}
	case 71:
		//line datetime.y:359
		{
			// 3rd Tuesday (of implicit this month)
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-0].intval, yyS[yypt-2].tval.i)
			l.setMonths(0, 0)
		}
	case 72:
		//line datetime.y:365
		{
			// 3rd Tuesday of (implicit this) March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			l.setMonths(yyS[yypt-0].intval, 0)
		}
	case 73:
		//line datetime.y:371
		{
			// 3rd Tuesday of 2012
			yylex.(*dateLexer).setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i, yyS[yypt-0].tval.i)
		}
	case 74:
		//line datetime.y:375
		{
			// 3rd Tuesday of March 2012
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-1].intval, 0, yyS[yypt-0].tval.i)
		}
	case 75:
		//line datetime.y:381
		{
			// 3rd Tuesday of next March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 77:
		//line datetime.y:390
		{
			yylex.(*dateLexer).setAgo()
		}
	case 80:
		//line datetime.y:399
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 81:
		//line datetime.y:402
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 82:
		//line datetime.y:405
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), 1)
		}
	case 83:
		//line datetime.y:408
		{
			// Special-case to handle "week" and "fortnight"
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].tval.i*yyS[yypt-0].intval)
		}
	case 84:
		//line datetime.y:412
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval*yyS[yypt-0].intval)
		}
	case 85:
		//line datetime.y:415
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-0].intval)
		}
	case 86:
		//line datetime.y:418
		{
			// As we need to be able to separate out YD from HS in ISO durations
			// this becomes a fair bit messier than if Y D H S were just T_OFFSET
			// Because writing "next y" or "two h" would be odd, disallow
			// T_RELATIVE tokens from being used with ISO single-letter notation
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 87:
		//line datetime.y:425
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 88:
		//line datetime.y:428
		{
			// Resolve 'm' ambiguity in favour of minutes outside ISO duration
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 89:
		//line datetime.y:432
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].tval.i*7)
		}
	case 90:
		//line datetime.y:435
		{
			// yesterday or tomorrow
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-0].intval)
		}
	case 93:
		//line datetime.y:444
		{
			yylex.(*dateLexer).addOffset(O_DAY, 7*yyS[yypt-1].tval.i)
		}
	case 96:
		//line datetime.y:454
		{
			// takes care of Y and D
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 97:
		//line datetime.y:458
		{
			yylex.(*dateLexer).addOffset(O_MONTH, yyS[yypt-1].tval.i)
		}
	case 100:
		//line datetime.y:467
		{
			// takes care of H and S
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 101:
		//line datetime.y:471
		{
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 105:
		//line datetime.y:483
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 8 {
				// assume ISO 8601 YYYYMMDD
				l.setYMD(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l)
			} else if yyS[yypt-0].tval.l == 7 {
				// assume ISO 8601 ordinal YYYYDDD
				l.setDate(yyS[yypt-0].tval.i/1000, 1, yyS[yypt-0].tval.i%1000)
			} else if yyS[yypt-0].tval.l == 6 {
				// assume ISO 8601 HHMMSS with no zone
				l.setHMS(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l, nil)
			} else if yyS[yypt-0].tval.l == 4 {
				// Assuming HHMM because that's more useful on IRC.
				l.setHMS(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l, nil)
			} else if yyS[yypt-0].tval.l == 2 {
				// assume HH with no zone
				l.setHMS(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l, nil)
			}
		}
	}
	goto yystack /* stack new state and value */
}
