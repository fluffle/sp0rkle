
//line datetime.y:2
package datetime

// Based upon parse-datetime.y in GNU coreutils.
// also an exercise in learning goyacc in particular.
// This file contains the yacc grammar only.
// See lexer.go for the lexer and parse functions,
// and tokenmaps.go for the token maps.

import (
	"fmt"
	"time"
)

type textint struct {
	i, l int
	s string
}


//line datetime.y:22
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
const T_MONTHNAME = 57350
const T_DAYNAME = 57351
const T_DAYS = 57352
const T_DAYSHIFT = 57353
const T_OFFSET = 57354
const T_ISOYD = 57355
const T_ISOHS = 57356
const T_RELATIVE = 57357
const T_AGO = 57358
const T_ZONE = 57359

var yyToknames = []string{
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

//line datetime.y:466


//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 15,
	1, 96,
	5, 96,
	8, 9,
	11, 96,
	15, 96,
	19, 9,
	-2, 17,
}

const yyNprod = 97
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 184

var yyAct = []int{

	108, 30, 103, 19, 46, 79, 74, 62, 27, 60,
	34, 98, 100, 38, 29, 43, 35, 140, 109, 4,
	119, 101, 32, 7, 137, 21, 41, 99, 136, 67,
	39, 73, 118, 40, 31, 63, 33, 36, 37, 61,
	15, 23, 24, 16, 17, 100, 20, 43, 75, 82,
	18, 66, 81, 71, 101, 92, 26, 111, 41, 22,
	89, 90, 39, 63, 117, 40, 107, 96, 80, 97,
	43, 75, 83, 87, 88, 47, 43, 75, 106, 32,
	86, 41, 43, 75, 58, 39, 115, 41, 40, 51,
	146, 50, 116, 41, 70, 141, 69, 120, 121, 52,
	54, 105, 53, 55, 56, 122, 54, 126, 53, 55,
	56, 123, 57, 110, 28, 23, 24, 132, 57, 124,
	43, 75, 134, 113, 68, 133, 112, 147, 143, 142,
	26, 41, 135, 49, 48, 51, 145, 50, 28, 23,
	24, 144, 139, 148, 138, 131, 130, 129, 143, 128,
	127, 125, 104, 114, 91, 85, 84, 77, 76, 72,
	65, 64, 44, 93, 102, 95, 59, 94, 25, 14,
	13, 12, 11, 10, 9, 8, 6, 5, 78, 45,
	3, 2, 1, 42,
}
var yyPact = []int{

	-7, -1000, -1000, 35, 133, -1000, -1000, -1000, -1000, -15,
	-1000, -1000, -1000, -1000, -1000, 9, 157, 57, 125, 90,
	-1000, 68, 34, 156, 155, 109, 84, -1000, -1000, 154,
	114, 153, -1000, 152, 49, 44, 151, 150, 71, 51,
	38, -1000, 149, -1000, 159, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, 6,
	-1000, -1, -1000, 147, -1000, -1000, -1000, 96, 79, -1000,
	-1000, -1000, 76, -1000, -1000, -1000, 41, -9, 105, -1000,
	37, 119, 116, 148, -1000, 114, 49, -1000, 42, -1000,
	10, -5, 57, -1000, 133, -1000, -1000, -1000, 32, -1000,
	-1000, -1000, 147, -1000, 97, 146, 114, 145, -1000, 144,
	142, -1000, 141, 140, 110, -1000, 117, 5, 1, 139,
	137, -11, -1000, -1000, -1000, 70, -1000, 64, -1000, -1000,
	-1000, -1000, 136, 131, -1000, 82, -1000, -1000, -1000, -1000,
	-1000, 122, 114, -1000, -1000, -1000, -1000, 114, -1000,
}
var yyPgo = []int{

	0, 183, 1, 3, 6, 0, 182, 181, 180, 4,
	179, 5, 178, 10, 177, 176, 23, 175, 174, 173,
	172, 171, 170, 169, 25, 168, 167, 166, 165, 9,
	7, 164, 2,
}
var yyR1 = []int{

	0, 6, 6, 9, 10, 10, 11, 12, 12, 13,
	13, 2, 2, 2, 2, 1, 1, 3, 3, 3,
	4, 4, 4, 5, 5, 7, 8, 8, 14, 14,
	14, 14, 14, 14, 14, 14, 14, 15, 15, 15,
	16, 16, 16, 17, 17, 17, 17, 17, 17, 17,
	18, 18, 18, 18, 18, 19, 19, 20, 20, 20,
	20, 20, 20, 20, 20, 20, 20, 20, 21, 21,
	24, 24, 25, 25, 25, 25, 25, 25, 25, 25,
	26, 25, 22, 22, 22, 27, 27, 30, 30, 31,
	31, 32, 32, 29, 28, 28, 23,
}
var yyR2 = []int{

	0, 1, 1, 1, 0, 1, 2, 0, 1, 0,
	1, 2, 4, 2, 4, 1, 1, 1, 2, 2,
	1, 2, 4, 0, 1, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 3, 5, 7,
	2, 4, 6, 3, 5, 4, 3, 5, 5, 5,
	3, 5, 3, 4, 6, 3, 4, 2, 1, 2,
	2, 2, 3, 5, 5, 6, 6, 1, 1, 2,
	1, 2, 2, 2, 2, 2, 2, 2, 2, 2,
	0, 5, 3, 2, 3, 1, 2, 2, 2, 1,
	2, 2, 2, 2, 0, 1, 1,
}
var yyChk = []int{

	-1000, -6, -7, -8, 26, -14, -15, -16, -17, -18,
	-19, -20, -21, -22, -23, 5, 8, 9, 15, -3,
	11, -24, 24, 6, 7, -25, 21, -3, 5, 29,
	-2, 25, -4, 27, -13, 7, 28, 29, 4, 21,
	24, 17, -1, 6, 5, -10, -9, 18, 9, 8,
	12, 10, 9, 12, 10, 13, 14, 22, 16, -27,
	-29, 5, -30, 29, 5, 5, -24, -3, 15, 12,
	10, -16, 5, -5, -4, 7, 5, 5, -12, -11,
	19, 8, 5, 28, 5, 5, 9, 22, 23, 22,
	23, 5, -13, 4, -26, -28, -30, -29, 5, 28,
	13, 22, -31, -32, 5, 25, -2, 25, -5, 27,
	8, 20, 7, 7, 5, -5, -11, 22, 22, 25,
	-9, -3, -32, 14, 22, 5, -5, 5, 5, 5,
	5, 5, 7, 8, 5, 15, 23, 23, 5, 5,
	28, 25, -2, -5, 5, 5, 8, 5, -5,
}
var yyDef = []int{

	26, -2, 1, 2, 0, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, -2, 58, 4, 0, 0,
	67, 68, 0, 0, 0, 70, 0, 25, 17, 0,
	23, 0, 40, 0, 7, 0, 0, 0, 10, 0,
	0, 20, 0, 15, 9, 57, 5, 3, 59, 60,
	73, 76, 61, 72, 75, 78, 79, 80, 69, 94,
	83, 0, 85, 0, 18, 19, 71, 0, 0, 74,
	77, 55, 0, 37, 24, 16, 23, 43, 0, 8,
	0, 0, 50, 0, 52, 23, 62, 11, 0, 13,
	0, 21, 46, 10, 0, 82, 86, 95, 0, 84,
	87, 88, 93, 89, 0, 0, 23, 0, 41, 0,
	45, 6, 0, 0, 53, 56, 0, 0, 0, 0,
	0, 0, 90, 91, 92, 23, 38, 23, 44, 47,
	48, 51, 0, 63, 64, 0, 12, 14, 22, 49,
	81, 0, 23, 42, 54, 65, 66, 23, 39,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 18, 3, 23, 27, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 25, 3,
	3, 3, 3, 3, 26, 21, 3, 3, 3, 3,
	20, 3, 3, 3, 3, 3, 3, 22, 3, 19,
	24, 3, 3, 3, 29, 3, 3, 28,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17,
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

	case 11:
		//line datetime.y:57
		{
		    yyVAL.intval = 0
		}
	case 12:
		//line datetime.y:60
		{
		    yyVAL.intval = 0
		}
	case 13:
		//line datetime.y:63
		{
		    yyVAL.intval = 12
		}
	case 14:
		//line datetime.y:66
		{
		    yyVAL.intval = 12
		}
	case 15:
		yyVAL.intval = yyS[yypt-0].intval
	case 16:
		yyVAL.intval = yyS[yypt-0].intval
	case 17:
		yyVAL.tval = yyS[yypt-0].tval
	case 18:
		//line datetime.y:75
		{
			yyS[yypt-0].tval.s = "+" + yyS[yypt-0].tval.s
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 19:
		//line datetime.y:79
		{
			yyS[yypt-0].tval.s = "-" + yyS[yypt-0].tval.s
			yyS[yypt-0].tval.i *= -1
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 20:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 21:
		//line datetime.y:87
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
	case 22:
		//line datetime.y:98
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-3].intval * (3600 * yyS[yypt-2].tval.i + 60 * yyS[yypt-0].tval.i))
		}
	case 23:
		//line datetime.y:103
		{ yyVAL.zoneval = nil }
	case 24:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 25:
		//line datetime.y:107
		{
			l := yylex.(*dateLexer)
			if ! l.state(HAVE_TIME, true) {
				l.time = time.Unix(int64(yyS[yypt-0].tval.i), 0)
			}
		}
	case 37:
		//line datetime.y:132
		{
			yylex.(*dateLexer).setTime(yyS[yypt-2].tval.i + yyS[yypt-1].intval, 0, 0, yyS[yypt-0].zoneval)
		}
	case 38:
		//line datetime.y:135
		{
			yylex.(*dateLexer).setTime(yyS[yypt-4].tval.i + yyS[yypt-1].intval, yyS[yypt-2].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 39:
		//line datetime.y:138
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i + yyS[yypt-1].intval, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].zoneval)
		}
	case 40:
		//line datetime.y:145
		{
			yylex.(*dateLexer).setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
		}
	case 41:
		//line datetime.y:148
		{
			yylex.(*dateLexer).setTime(yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 42:
		//line datetime.y:151
		{
			yylex.(*dateLexer).setTime(yyS[yypt-5].tval.i, yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, yyS[yypt-0].zoneval)
		}
	case 43:
		//line datetime.y:159
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
	case 44:
		//line datetime.y:169
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
	case 45:
		//line datetime.y:182
		{
			// DDth of Mon
		yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 46:
		//line datetime.y:186
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
	case 47:
		//line datetime.y:196
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
	case 48:
		//line datetime.y:209
		{
		    // RFC 850, srsly :(
		l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume DD-Mon-YYYY
			l.setDate(yyS[yypt-0].tval.i, yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume DD-Mon-YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].tval.i + 1900, yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			} else {
				// assume DD-Mon-YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].tval.i + 2000, yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			}
		}
	case 49:
		//line datetime.y:223
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume Mon DDth, YYYY
			l.setDate(yyS[yypt-0].tval.i, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume Mon DDth, YY, add 1900 if YY > 68
			l.setDate(yyS[yypt-0].tval.i + 1900, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			} else {
				// assume Mon DDth, YY, add 2000 otherwise
			l.setDate(yyS[yypt-0].tval.i + 2000, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			}
		}
	case 50:
		//line datetime.y:239
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
	case 51:
		//line datetime.y:253
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
	case 52:
		//line datetime.y:266
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
	case 53:
		//line datetime.y:276
		{
			// assume YYYY-'W'WW
		yylex.(*dateLexer).setWeek(yyS[yypt-3].tval.i, yyS[yypt-0].tval.i, 1)
		}
	case 54:
		//line datetime.y:280
		{
			// assume YYYY-'W'WW-D
		yylex.(*dateLexer).setWeek(yyS[yypt-5].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
		}
	case 56:
		//line datetime.y:288
		{
			// this goes here because the YYYYMMDD and HHMMSS forms of the
		// ISO 8601 format date and time are handled by 'integer' below.
		l := yylex.(*dateLexer)
			if yyS[yypt-3].tval.l == 8 {
				// assume ISO 8601 YYYYMMDD
			l.setYMD(yyS[yypt-3].tval.i, yyS[yypt-3].tval.l)
	        } else if yyS[yypt-3].tval.l == 7 {
	            // assume ISO 8601 ordinal YYYYDDD
			l.setDate(yyS[yypt-3].tval.i / 1000, 1, yyS[yypt-3].tval.i % 1000)
	        }
			l.setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
		}
	case 57:
		//line datetime.y:303
		{
			// Tuesday
		yylex.(*dateLexer).setDay(yyS[yypt-1].intval, 0)
		}
	case 58:
		//line datetime.y:307
		{
			// March
		yylex.(*dateLexer).setMonth(yyS[yypt-0].intval, 0)
		}
	case 59:
		//line datetime.y:311
		{
			// Next tuesday
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 60:
		//line datetime.y:315
		{
			// Next march
		yylex.(*dateLexer).setMonth(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 61:
		//line datetime.y:319
		{
			// +-N Tuesdays
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-1].tval.i)
		}
	case 62:
		//line datetime.y:323
		{
			// 3rd Tuesday 
		yylex.(*dateLexer).setDay(yyS[yypt-0].intval, yyS[yypt-2].tval.i)
		}
	case 63:
		//line datetime.y:327
		{
			// 3rd Tuesday of (implicit this) March
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			l.setMonth(yyS[yypt-0].intval, 1)
		}
	case 64:
		//line datetime.y:333
		{
			// 3rd Tuesday of 2012
		yylex.(*dateLexer).setDay(yyS[yypt-2].intval, yyS[yypt-4].tval.i, yyS[yypt-0].tval.i)
		}
	case 65:
		//line datetime.y:337
		{
			// 3rd Tuesday of March 2012
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonth(yyS[yypt-1].intval, 1, yyS[yypt-0].tval.i)
		}
	case 66:
		//line datetime.y:343
		{
			// 3rd Tuesday of next March
		l := yylex.(*dateLexer)
			l.setDay(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonth(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 67:
		//line datetime.y:349
		{
			// yesterday or tomorrow
		d := time.Now().Weekday()
			yylex.(*dateLexer).setDay((7+int(d)+yyS[yypt-0].intval)%7, yyS[yypt-0].intval)
		}
	case 69:
		//line datetime.y:357
		{
			yylex.(*dateLexer).setAgo()
		}
	case 72:
		//line datetime.y:366
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 73:
		//line datetime.y:369
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 74:
		//line datetime.y:372
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), 1)
		}
	case 75:
		//line datetime.y:375
		{
			// Special-case to handle "week" and "fortnight"
		yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].tval.i * yyS[yypt-0].intval)
		}
	case 76:
		//line datetime.y:379
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval * yyS[yypt-0].intval)
		}
	case 77:
		//line datetime.y:382
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-0].intval)
		}
	case 78:
		//line datetime.y:385
		{
			// As we need to be able to separate out YD from HS in ISO durations
		// this becomes a fair bit messier than if Y D H S were just T_OFFSET
		// Because writing "next y" or "two h" would be odd, disallow
		// T_RELATIVE tokens from being used with ISO single-letter notation
		yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 79:
		//line datetime.y:392
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 80:
		//line datetime.y:395
		{
			// Resolve 'm' ambiguity in favour of minutes outside ISO duration
		yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 81:
		//line datetime.y:398
		{
		    yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-4].tval.i * 7)
		}
	case 84:
		//line datetime.y:406
		{
			yylex.(*dateLexer).addOffset(O_DAY, 7 * yyS[yypt-1].tval.i)
		}
	case 87:
		//line datetime.y:416
		{
			// takes care of Y and D
		yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 88:
		//line datetime.y:420
		{
			yylex.(*dateLexer).addOffset(O_MONTH, yyS[yypt-1].tval.i)
		}
	case 91:
		//line datetime.y:429
		{
			// takes care of H and S
		yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 92:
		//line datetime.y:433
		{
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 96:
		//line datetime.y:445
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 8 {
				// assume ISO 8601 YYYYMMDD
			l.setYMD(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l)
	        } else if yyS[yypt-0].tval.l == 7 {
	            // assume ISO 8601 ordinal YYYYDDD
			l.setDate(yyS[yypt-0].tval.i / 1000, 1, yyS[yypt-0].tval.i % 1000)
			} else if yyS[yypt-0].tval.l == 6 {
				// assume ISO 8601 HHMMSS with no zone
			l.setHMS(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l, nil)
			} else if yyS[yypt-0].tval.l == 4 {
				// assume setting YYYY, because otherwise parsing ANSIC, UnixTime
			// and RubyTime formats fails as the year is after the time
			// Probably should be HHMM instead...
			l.setYear(yyS[yypt-0].tval.i)
			} else if yyS[yypt-0].tval.l == 2 {
	            // assume HH with no zone
            l.setHMS(yyS[yypt-0].tval.i, yyS[yypt-0].tval.l, nil)
	        }
		}
	}
	goto yystack /* stack new state and value */
}
