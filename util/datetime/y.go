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
	s    string
}

//line datetime.y:22
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

//line datetime.y:501

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
	-1, 145,
	8, 4,
	-2, 52,
}

const yyNprod = 106
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 198

var yyAct = []int{

	113, 51, 32, 33, 52, 108, 21, 68, 36, 80,
	66, 29, 31, 147, 105, 87, 7, 114, 86, 4,
	103, 106, 40, 123, 47, 37, 34, 104, 23, 94,
	95, 122, 67, 79, 73, 45, 88, 41, 143, 44,
	42, 43, 69, 35, 38, 39, 17, 15, 77, 16,
	25, 26, 18, 19, 69, 22, 72, 47, 81, 20,
	92, 93, 142, 28, 47, 81, 24, 121, 45, 64,
	53, 91, 44, 101, 43, 45, 102, 41, 154, 44,
	42, 43, 110, 47, 81, 111, 112, 129, 34, 145,
	57, 119, 56, 130, 45, 115, 41, 58, 60, 42,
	59, 61, 62, 105, 47, 81, 127, 125, 63, 76,
	106, 75, 132, 128, 138, 45, 60, 135, 59, 61,
	62, 117, 116, 30, 25, 26, 63, 161, 140, 160,
	146, 139, 157, 74, 150, 148, 149, 28, 141, 55,
	54, 57, 156, 56, 30, 25, 26, 155, 153, 152,
	158, 151, 144, 137, 136, 134, 133, 131, 150, 109,
	118, 98, 96, 90, 89, 83, 82, 78, 71, 70,
	48, 126, 97, 50, 124, 120, 85, 107, 100, 65,
	99, 27, 14, 13, 12, 11, 10, 9, 8, 6,
	5, 49, 84, 3, 2, 1, 159, 46,
}
var yyPact = []int{

	-8, -1000, -1000, 41, 136, -1000, -1000, -1000, -1000, -18,
	-1000, -1000, -1000, -1000, -1000, -1000, 15, 162, 168, 49,
	128, 85, -1000, 50, 24, 161, 160, 115, 96, -1000,
	-1000, 159, 95, 158, -1000, 157, 172, 7, 156, 155,
	59, 37, 6, -1000, -1000, -1000, 154, -1000, 165, 153,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, 12, -1000, -2, -1000, 151,
	-1000, -1000, -1000, 103, 77, -1000, -1000, -1000, 48, -1000,
	-1000, -1000, 55, -11, 84, -1000, 112, 111, 152, -1000,
	95, 171, -1000, 44, -1000, 8, -3, 170, 164, 136,
	-1000, -1000, -1000, 87, -1000, -1000, -1000, 151, -1000, 70,
	149, 95, 148, -1000, 147, 49, 146, 145, 104, -1000,
	120, 38, 14, 144, 78, 49, -1000, -16, -1000, -1000,
	-1000, 48, -1000, 74, -1000, 143, -1000, -1000, 141, 140,
	-1000, 67, -1000, -1000, -1000, 49, 134, -1000, 124, 95,
	121, -1000, -1000, -1000, -1000, 119, -1000, 95, -1000, -1000,
	-1000, -1000,
}
var yyPgo = []int{

	0, 197, 2, 6, 196, 9, 0, 195, 194, 193,
	4, 1, 192, 191, 8, 3, 190, 189, 16, 188,
	187, 186, 185, 184, 183, 182, 28, 181, 180, 179,
	178, 10, 7, 177, 5,
}
var yyR1 = []int{

	0, 7, 7, 10, 11, 11, 12, 12, 13, 13,
	14, 14, 4, 4, 2, 2, 2, 2, 15, 15,
	1, 1, 3, 3, 3, 5, 5, 5, 6, 6,
	8, 9, 9, 16, 16, 16, 16, 16, 16, 16,
	16, 16, 16, 17, 17, 17, 18, 18, 18, 19,
	19, 19, 19, 19, 19, 19, 19, 19, 19, 20,
	20, 20, 20, 20, 21, 21, 22, 22, 22, 22,
	22, 22, 22, 22, 22, 22, 22, 23, 23, 26,
	26, 27, 27, 27, 27, 27, 27, 27, 27, 28,
	27, 24, 24, 24, 29, 29, 32, 32, 33, 33,
	34, 34, 31, 30, 30, 25,
}
var yyR2 = []int{

	0, 1, 1, 1, 0, 1, 0, 1, 0, 1,
	0, 1, 0, 1, 2, 4, 2, 4, 1, 1,
	1, 1, 1, 2, 2, 1, 2, 4, 0, 1,
	2, 0, 2, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 3, 5, 7, 2, 4, 7, 3,
	5, 3, 5, 7, 4, 4, 6, 5, 6, 3,
	5, 3, 4, 6, 3, 4, 2, 1, 2, 2,
	2, 3, 5, 5, 6, 6, 1, 1, 2, 1,
	2, 2, 2, 2, 2, 2, 2, 2, 2, 0,
	5, 3, 2, 3, 1, 2, 2, 2, 1, 2,
	2, 2, 2, 0, 1, 1,
}
var yyChk = []int{

	-1000, -7, -8, -9, 27, -16, -17, -18, -19, -20,
	-21, -22, -23, -24, -25, 6, 8, 5, 11, 12,
	18, -3, 14, -26, 25, 9, 10, -27, 22, -3,
	8, 30, -2, -15, -5, 28, -14, 10, 29, 30,
	7, 22, 25, 26, 24, 20, -1, 9, 8, -13,
	5, -11, -10, 21, 12, 11, 15, 13, 12, 15,
	13, 16, 17, 23, 19, -29, -31, 8, -32, 30,
	8, 8, -26, -3, 18, 15, 13, -18, 8, -6,
	-5, 10, 8, 8, -12, 4, 11, 8, 29, 8,
	8, 12, 23, 24, 23, 24, 8, 7, 8, -28,
	-30, -32, -31, 8, 29, 16, 23, -33, -34, 8,
	-15, -2, -15, -6, 28, 11, 10, 10, 8, -6,
	4, 23, 23, 26, 4, -14, 7, -3, -34, 17,
	23, 8, -6, 8, 8, -11, 8, 8, 10, 11,
	8, 18, 24, 24, 8, 11, -10, 29, -15, -2,
	-6, 8, 8, 8, 11, -11, 8, 8, -6, -4,
	8, 8,
}
var yyDef = []int{

	31, -2, 1, 2, 0, 32, 33, 34, 35, 36,
	37, 38, 39, 40, 41, 42, -2, 0, -2, 4,
	0, 0, 76, 77, 0, 0, 0, 79, 0, 30,
	22, 0, 28, 0, 46, 0, 6, 0, 0, 0,
	11, 0, 0, 18, 19, 25, 0, 20, 0, 0,
	9, 66, 5, 3, 68, 69, 82, 85, 70, 81,
	84, 87, 88, 89, 78, 103, 92, 0, 94, 0,
	23, 24, 80, 0, 0, 83, 86, 64, 0, 43,
	29, 21, 28, 49, 0, 7, 0, 59, 0, 61,
	28, 71, 14, 0, 16, 0, 26, 51, 10, 0,
	91, 95, 104, 0, 93, 96, 97, 102, 98, 0,
	0, 28, 0, 47, 0, -2, 0, 0, 62, 65,
	0, 0, 0, 0, 0, 55, 11, 0, 99, 100,
	101, 28, 44, 28, 50, 0, 57, 60, 0, 72,
	73, 0, 15, 17, 27, -2, 0, 90, 0, 28,
	12, 56, 63, 74, 75, 0, 58, 28, 45, 48,
	13, 53,
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

	case 12:
		//line datetime.y:55
		{
			yyVAL.tval = textint{}
		}
	case 13:
		yyVAL.tval = yyS[yypt-0].tval
	case 14:
		//line datetime.y:58
		{
			yyVAL.intval = 0
		}
	case 15:
		//line datetime.y:61
		{
			yyVAL.intval = 0
		}
	case 16:
		//line datetime.y:64
		{
			yyVAL.intval = 12
		}
	case 17:
		//line datetime.y:67
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
		//line datetime.y:78
		{
			yyS[yypt-0].tval.s = "+" + yyS[yypt-0].tval.s
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 24:
		//line datetime.y:82
		{
			yyS[yypt-0].tval.s = "-" + yyS[yypt-0].tval.s
			yyS[yypt-0].tval.i *= -1
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 25:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 26:
		//line datetime.y:90
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
		//line datetime.y:101
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-3].intval*(3600*yyS[yypt-2].tval.i+60*yyS[yypt-0].tval.i))
		}
	case 28:
		//line datetime.y:106
		{
			yyVAL.zoneval = nil
		}
	case 29:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 30:
		//line datetime.y:110
		{
			l := yylex.(*dateLexer)
			if !l.state(HAVE_TIME, true) {
				l.time = time.Unix(int64(yyS[yypt-0].tval.i), 0)
			}
		}
	case 43:
		//line datetime.y:136
		{
			l := yylex.(*dateLexer)
			// Hack to allow HHMMam to parse correctly, cos adie is a mong.
			if yyS[yypt-2].tval.l == 3 || yyS[yypt-2].tval.l == 4 {
				l.setTime(yyS[yypt-2].tval.i/100+yyS[yypt-1].intval, yyS[yypt-2].tval.i%100, 0, yyS[yypt-0].zoneval)
			} else {
				l.setTime(yyS[yypt-2].tval.i+yyS[yypt-1].intval, 0, 0, yyS[yypt-0].zoneval)
			}
		}
	case 44:
		//line datetime.y:145
		{
			yylex.(*dateLexer).setTime(yyS[yypt-4].tval.i+yyS[yypt-1].intval, yyS[yypt-2].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 45:
		//line datetime.y:148
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i+yyS[yypt-1].intval, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].zoneval)
		}
	case 46:
		//line datetime.y:155
		{
			yylex.(*dateLexer).setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
		}
	case 47:
		//line datetime.y:158
		{
			yylex.(*dateLexer).setTime(yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 48:
		//line datetime.y:161
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-1].zoneval)
			// Hack to make time.ANSIC, time.UnixDate and time.RubyDate parse
			if yyS[yypt-0].tval.l == 4 {
				yylex.(*dateLexer).setYear(yyS[yypt-0].tval.i)
			}
		}
	case 49:
		//line datetime.y:173
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
		//line datetime.y:183
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
		//line datetime.y:196
		{
			// the DDth
			yylex.(*dateLexer).setDay(yyS[yypt-1].tval.i)
		}
	case 52:
		//line datetime.y:200
		{
			// the DDth of Month
			yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 53:
		//line datetime.y:204
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
		//line datetime.y:217
		{
			// DD[th] [of] Month
			yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 55:
		//line datetime.y:221
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
		//line datetime.y:231
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
		//line datetime.y:244
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
		//line datetime.y:258
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
		//line datetime.y:276
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
		//line datetime.y:290
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
		//line datetime.y:303
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
		//line datetime.y:313
		{
			// assume YYYY-'W'WW
			yylex.(*dateLexer).setWeek(yyS[yypt-3].tval.i, yyS[yypt-0].tval.i, 1)
		}
	case 63:
		//line datetime.y:317
		{
			// assume YYYY-'W'WW-D
			yylex.(*dateLexer).setWeek(yyS[yypt-5].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
		}
	case 65:
		//line datetime.y:325
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
		//line datetime.y:340
		{
			// Tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-1].intval, 0)
		}
	case 67:
		//line datetime.y:344
		{
			// March
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, 0)
		}
	case 68:
		//line datetime.y:348
		{
			// Next tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 69:
		//line datetime.y:352
		{
			// Next march
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 70:
		//line datetime.y:356
		{
			// +-N Tuesdays
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].tval.i)
		}
	case 71:
		//line datetime.y:360
		{
			// 3rd Tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-2].tval.i)
		}
	case 72:
		//line datetime.y:364
		{
			// 3rd Tuesday of (implicit this) March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			l.setMonths(yyS[yypt-0].intval, 1)
		}
	case 73:
		//line datetime.y:370
		{
			// 3rd Tuesday of 2012
			yylex.(*dateLexer).setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i, yyS[yypt-0].tval.i)
		}
	case 74:
		//line datetime.y:374
		{
			// 3rd Tuesday of March 2012
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-1].intval, 1, yyS[yypt-0].tval.i)
		}
	case 75:
		//line datetime.y:380
		{
			// 3rd Tuesday of next March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 76:
		//line datetime.y:386
		{
			// yesterday or tomorrow
			d := time.Now().Weekday()
			yylex.(*dateLexer).setDays((7+int(d)+yyS[yypt-0].intval)%7, yyS[yypt-0].intval)
		}
	case 78:
		//line datetime.y:394
		{
			yylex.(*dateLexer).setAgo()
		}
	case 81:
		//line datetime.y:403
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 82:
		//line datetime.y:406
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 83:
		//line datetime.y:409
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), 1)
		}
	case 84:
		//line datetime.y:412
		{
			// Special-case to handle "week" and "fortnight"
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].tval.i*yyS[yypt-0].intval)
		}
	case 85:
		//line datetime.y:416
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval*yyS[yypt-0].intval)
		}
	case 86:
		//line datetime.y:419
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-0].intval)
		}
	case 87:
		//line datetime.y:422
		{
			// As we need to be able to separate out YD from HS in ISO durations
			// this becomes a fair bit messier than if Y D H S were just T_OFFSET
			// Because writing "next y" or "two h" would be odd, disallow
			// T_RELATIVE tokens from being used with ISO single-letter notation
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 88:
		//line datetime.y:429
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 89:
		//line datetime.y:432
		{
			// Resolve 'm' ambiguity in favour of minutes outside ISO duration
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 90:
		//line datetime.y:435
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-4].tval.i*7)
		}
	case 93:
		//line datetime.y:443
		{
			yylex.(*dateLexer).addOffset(O_DAY, 7*yyS[yypt-1].tval.i)
		}
	case 96:
		//line datetime.y:453
		{
			// takes care of Y and D
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 97:
		//line datetime.y:457
		{
			yylex.(*dateLexer).addOffset(O_MONTH, yyS[yypt-1].tval.i)
		}
	case 100:
		//line datetime.y:466
		{
			// takes care of H and S
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 101:
		//line datetime.y:470
		{
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 105:
		//line datetime.y:482
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
