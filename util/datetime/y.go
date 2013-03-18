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

const T_DAYQUAL = 57346
const T_THE = 57347
const T_INTEGER = 57348
const T_PLUS = 57349
const T_MINUS = 57350
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
	"T_THE",
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

//line datetime.y:476

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 15,
	1, 99,
	5, 99,
	6, 99,
	9, 9,
	12, 99,
	16, 99,
	20, 9,
	-2, 19,
}

const yyNprod = 100
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 190

var yyAct = []int{

	113, 31, 108, 20, 32, 83, 50, 78, 28, 66,
	64, 35, 103, 39, 30, 145, 46, 36, 114, 4,
	7, 124, 142, 33, 93, 94, 65, 44, 91, 92,
	71, 40, 77, 43, 41, 42, 67, 34, 37, 38,
	105, 16, 15, 24, 25, 17, 18, 141, 21, 106,
	67, 75, 19, 46, 79, 104, 86, 128, 27, 85,
	97, 23, 22, 123, 44, 129, 122, 84, 40, 105,
	43, 41, 42, 101, 102, 46, 79, 116, 106, 87,
	51, 110, 111, 62, 33, 112, 44, 46, 79, 70,
	120, 90, 43, 55, 42, 54, 121, 58, 44, 57,
	59, 60, 40, 126, 125, 41, 74, 151, 73, 61,
	127, 115, 131, 56, 58, 137, 57, 59, 60, 118,
	29, 24, 25, 29, 24, 25, 61, 46, 79, 139,
	72, 152, 138, 148, 147, 146, 27, 117, 44, 140,
	53, 52, 55, 150, 54, 149, 144, 143, 153, 136,
	135, 134, 133, 148, 132, 130, 109, 119, 95, 89,
	88, 81, 80, 76, 69, 68, 48, 47, 98, 96,
	107, 100, 63, 99, 26, 14, 13, 12, 11, 10,
	9, 8, 6, 5, 82, 49, 3, 2, 1, 45,
}
var yyPact = []int{

	-8, -1000, -1000, 36, 117, -1000, -1000, -1000, -1000, -16,
	-1000, -1000, -1000, -1000, -1000, 9, 161, 160, 61, 131,
	103, -1000, 66, 20, 159, 158, 114, 95, -1000, -1000,
	157, 120, 156, -1000, 155, 47, 50, 154, 153, 81,
	5, 1, -1000, -1000, -1000, 152, -1000, 165, 164, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, 6, -1000, 26, -1000, 150, -1000, -1000,
	-1000, 86, 82, -1000, -1000, -1000, 68, -1000, -1000, -1000,
	46, -10, 102, -1000, 56, 129, 111, 151, -1000, 120,
	47, -1000, 43, -1000, 40, -5, -1000, 61, -1000, 117,
	-1000, -1000, -1000, 55, -1000, -1000, -1000, 150, -1000, 42,
	149, 120, 148, -1000, 146, 145, -1000, 144, 143, 107,
	-1000, 123, 23, -2, 141, 140, -14, -1000, -1000, -1000,
	68, -1000, 80, -1000, -1000, -1000, -1000, 139, 137, -1000,
	98, -1000, -1000, -1000, -1000, -1000, 125, 120, -1000, -1000,
	-1000, -1000, 120, -1000,
}
var yyPgo = []int{

	0, 189, 1, 3, 7, 0, 188, 187, 186, 6,
	185, 5, 184, 11, 4, 183, 182, 20, 181, 180,
	179, 178, 177, 176, 175, 62, 174, 173, 172, 171,
	10, 9, 170, 2,
}
var yyR1 = []int{

	0, 6, 6, 9, 10, 10, 11, 12, 12, 13,
	13, 2, 2, 2, 2, 14, 14, 1, 1, 3,
	3, 3, 4, 4, 4, 5, 5, 7, 8, 8,
	15, 15, 15, 15, 15, 15, 15, 15, 15, 16,
	16, 16, 17, 17, 17, 18, 18, 18, 18, 18,
	18, 18, 18, 19, 19, 19, 19, 19, 20, 20,
	21, 21, 21, 21, 21, 21, 21, 21, 21, 21,
	21, 22, 22, 25, 25, 26, 26, 26, 26, 26,
	26, 26, 26, 27, 26, 23, 23, 23, 28, 28,
	31, 31, 32, 32, 33, 33, 30, 29, 29, 24,
}
var yyR2 = []int{

	0, 1, 1, 1, 0, 1, 2, 0, 1, 0,
	1, 2, 4, 2, 4, 1, 1, 1, 1, 1,
	2, 2, 1, 2, 4, 0, 1, 2, 0, 2,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 3,
	5, 7, 2, 4, 6, 3, 5, 3, 4, 3,
	5, 5, 5, 3, 5, 3, 4, 6, 3, 4,
	2, 1, 2, 2, 2, 3, 5, 5, 6, 6,
	1, 1, 2, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 0, 5, 3, 2, 3, 1, 2,
	2, 2, 1, 2, 2, 2, 2, 0, 1, 1,
}
var yyChk = []int{

	-1000, -6, -7, -8, 27, -15, -16, -17, -18, -19,
	-20, -21, -22, -23, -24, 6, 5, 9, 10, 16,
	-3, 12, -25, 25, 7, 8, -26, 22, -3, 6,
	30, -2, -14, -4, 28, -13, 8, 29, 30, 4,
	22, 25, 26, 24, 18, -1, 7, 6, 6, -10,
	-9, 19, 10, 9, 13, 11, 10, 13, 11, 14,
	15, 23, 17, -28, -30, 6, -31, 30, 6, 6,
	-25, -3, 16, 13, 11, -17, 6, -5, -4, 8,
	6, 6, -12, -11, 20, 9, 6, 29, 6, 6,
	10, 23, 24, 23, 24, 6, 4, -13, 4, -27,
	-29, -31, -30, 6, 29, 14, 23, -32, -33, 6,
	-14, -2, -14, -5, 28, 9, 21, 8, 8, 6,
	-5, -11, 23, 23, 26, -9, -3, -33, 15, 23,
	6, -5, 6, 6, 6, 6, 6, 8, 9, 6,
	16, 24, 24, 6, 6, 29, -14, -2, -5, 6,
	6, 9, 6, -5,
}
var yyDef = []int{

	28, -2, 1, 2, 0, 29, 30, 31, 32, 33,
	34, 35, 36, 37, 38, -2, 0, 61, 4, 0,
	0, 70, 71, 0, 0, 0, 73, 0, 27, 19,
	0, 25, 0, 42, 0, 7, 0, 0, 0, 10,
	0, 0, 15, 16, 22, 0, 17, 0, 9, 60,
	5, 3, 62, 63, 76, 79, 64, 75, 78, 81,
	82, 83, 72, 97, 86, 0, 88, 0, 20, 21,
	74, 0, 0, 77, 80, 58, 0, 39, 26, 18,
	25, 45, 0, 8, 0, 0, 53, 0, 55, 25,
	65, 11, 0, 13, 0, 23, 47, 49, 10, 0,
	85, 89, 98, 0, 87, 90, 91, 96, 92, 0,
	0, 25, 0, 43, 0, 48, 6, 0, 0, 56,
	59, 0, 0, 0, 0, 0, 0, 93, 94, 95,
	25, 40, 25, 46, 50, 51, 54, 0, 66, 67,
	0, 12, 14, 24, 52, 84, 0, 25, 44, 57,
	68, 69, 25, 41,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 19, 3, 24, 28, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 26, 3,
	3, 3, 3, 3, 27, 22, 3, 3, 3, 3,
	21, 3, 3, 3, 3, 3, 3, 23, 3, 20,
	25, 3, 3, 3, 30, 3, 3, 29,
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
	case 17:
		yyVAL.intval = yyS[yypt-0].intval
	case 18:
		yyVAL.intval = yyS[yypt-0].intval
	case 19:
		yyVAL.tval = yyS[yypt-0].tval
	case 20:
		//line datetime.y:77
		{
			yyS[yypt-0].tval.s = "+" + yyS[yypt-0].tval.s
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 21:
		//line datetime.y:81
		{
			yyS[yypt-0].tval.s = "-" + yyS[yypt-0].tval.s
			yyS[yypt-0].tval.i *= -1
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 22:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 23:
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
	case 24:
		//line datetime.y:100
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-3].intval*(3600*yyS[yypt-2].tval.i+60*yyS[yypt-0].tval.i))
		}
	case 25:
		//line datetime.y:105
		{
			yyVAL.zoneval = nil
		}
	case 26:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 27:
		//line datetime.y:109
		{
			l := yylex.(*dateLexer)
			if !l.state(HAVE_TIME, true) {
				l.time = time.Unix(int64(yyS[yypt-0].tval.i), 0)
			}
		}
	case 39:
		//line datetime.y:134
		{
			l := yylex.(*dateLexer)
			// Hack to allow HHMMam to parse correctly, cos adie is a mong.
			if yyS[yypt-2].tval.l == 3 || yyS[yypt-2].tval.l == 4 {
				l.setTime(yyS[yypt-2].tval.i/100+yyS[yypt-1].intval, yyS[yypt-2].tval.i%100, 0, yyS[yypt-0].zoneval)
			} else {
				l.setTime(yyS[yypt-2].tval.i+yyS[yypt-1].intval, 0, 0, yyS[yypt-0].zoneval)
			}
		}
	case 40:
		//line datetime.y:143
		{
			yylex.(*dateLexer).setTime(yyS[yypt-4].tval.i+yyS[yypt-1].intval, yyS[yypt-2].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 41:
		//line datetime.y:146
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i+yyS[yypt-1].intval, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].zoneval)
		}
	case 42:
		//line datetime.y:153
		{
			yylex.(*dateLexer).setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
		}
	case 43:
		//line datetime.y:156
		{
			yylex.(*dateLexer).setTime(yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 44:
		//line datetime.y:159
		{
			yylex.(*dateLexer).setTime(yyS[yypt-5].tval.i, yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, yyS[yypt-0].zoneval)
		}
	case 45:
		//line datetime.y:167
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
	case 46:
		//line datetime.y:177
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
	case 47:
		//line datetime.y:190
		{
			// the DDth
			yylex.(*dateLexer).setDay(yyS[yypt-1].tval.i)
		}
	case 48:
		//line datetime.y:194
		{
			// DDth of Mon
			yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 49:
		//line datetime.y:198
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
	case 50:
		//line datetime.y:208
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume DDth of Mon YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-1].intval, yyS[yypt-4].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume DDth of Mon YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-1].intval, yyS[yypt-4].tval.i)
			} else {
				// assume DDth of Mon YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-1].intval, yyS[yypt-4].tval.i)
			}
		}
	case 51:
		//line datetime.y:221
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
	case 52:
		//line datetime.y:235
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume Mon DDth, YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume Mon DDth, YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			} else {
				// assume Mon DDth, YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-4].intval, yyS[yypt-3].tval.i)
			}
		}
	case 53:
		//line datetime.y:251
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
	case 54:
		//line datetime.y:265
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
	case 55:
		//line datetime.y:278
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
	case 56:
		//line datetime.y:288
		{
			// assume YYYY-'W'WW
			yylex.(*dateLexer).setWeek(yyS[yypt-3].tval.i, yyS[yypt-0].tval.i, 1)
		}
	case 57:
		//line datetime.y:292
		{
			// assume YYYY-'W'WW-D
			yylex.(*dateLexer).setWeek(yyS[yypt-5].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
		}
	case 59:
		//line datetime.y:300
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
	case 60:
		//line datetime.y:315
		{
			// Tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-1].intval, 0)
		}
	case 61:
		//line datetime.y:319
		{
			// March
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, 0)
		}
	case 62:
		//line datetime.y:323
		{
			// Next tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 63:
		//line datetime.y:327
		{
			// Next march
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 64:
		//line datetime.y:331
		{
			// +-N Tuesdays
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].tval.i)
		}
	case 65:
		//line datetime.y:335
		{
			// 3rd Tuesday 
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-2].tval.i)
		}
	case 66:
		//line datetime.y:339
		{
			// 3rd Tuesday of (implicit this) March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			l.setMonths(yyS[yypt-0].intval, 1)
		}
	case 67:
		//line datetime.y:345
		{
			// 3rd Tuesday of 2012
			yylex.(*dateLexer).setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i, yyS[yypt-0].tval.i)
		}
	case 68:
		//line datetime.y:349
		{
			// 3rd Tuesday of March 2012
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-1].intval, 1, yyS[yypt-0].tval.i)
		}
	case 69:
		//line datetime.y:355
		{
			// 3rd Tuesday of next March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 70:
		//line datetime.y:361
		{
			// yesterday or tomorrow
			d := time.Now().Weekday()
			yylex.(*dateLexer).setDays((7+int(d)+yyS[yypt-0].intval)%7, yyS[yypt-0].intval)
		}
	case 72:
		//line datetime.y:369
		{
			yylex.(*dateLexer).setAgo()
		}
	case 75:
		//line datetime.y:378
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 76:
		//line datetime.y:381
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 77:
		//line datetime.y:384
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), 1)
		}
	case 78:
		//line datetime.y:387
		{
			// Special-case to handle "week" and "fortnight"
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].tval.i*yyS[yypt-0].intval)
		}
	case 79:
		//line datetime.y:391
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval*yyS[yypt-0].intval)
		}
	case 80:
		//line datetime.y:394
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-0].intval)
		}
	case 81:
		//line datetime.y:397
		{
			// As we need to be able to separate out YD from HS in ISO durations
			// this becomes a fair bit messier than if Y D H S were just T_OFFSET
			// Because writing "next y" or "two h" would be odd, disallow
			// T_RELATIVE tokens from being used with ISO single-letter notation
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 82:
		//line datetime.y:404
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 83:
		//line datetime.y:407
		{
			// Resolve 'm' ambiguity in favour of minutes outside ISO duration
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 84:
		//line datetime.y:410
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-4].tval.i*7)
		}
	case 87:
		//line datetime.y:418
		{
			yylex.(*dateLexer).addOffset(O_DAY, 7*yyS[yypt-1].tval.i)
		}
	case 90:
		//line datetime.y:428
		{
			// takes care of Y and D
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 91:
		//line datetime.y:432
		{
			yylex.(*dateLexer).addOffset(O_MONTH, yyS[yypt-1].tval.i)
		}
	case 94:
		//line datetime.y:441
		{
			// takes care of H and S
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 95:
		//line datetime.y:445
		{
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 99:
		//line datetime.y:457
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
