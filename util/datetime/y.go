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
const T_DAYQUAL = 57348
const T_INTEGER = 57349
const T_PLUS = 57350
const T_MINUS = 57351
const T_MONTHNAME = 57352
const T_DAYNAME = 57353
const T_DAYS = 57354
const T_DAYSHIFT = 57355
const T_OFFSET = 57356
const T_ISOYD = 57357
const T_ISOHS = 57358
const T_RELATIVE = 57359
const T_AGO = 57360
const T_ZONE = 57361

var yyToknames = []string{
	"T_OF",
	"T_THE",
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

//line datetime.y:492

//line yacctab:1
var yyExca = []int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 15,
	1, 102,
	4, 10,
	5, 102,
	7, 102,
	10, 10,
	13, 102,
	17, 102,
	-2, 20,
	-1, 17,
	7, 8,
	-2, 64,
	-1, 114,
	7, 4,
	-2, 51,
	-1, 124,
	7, 4,
	-2, 52,
	-1, 144,
	7, 4,
	-2, 49,
}

const yyNprod = 103
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 195

var yyAct = []int{

	112, 50, 31, 32, 20, 107, 35, 67, 79, 28,
	65, 30, 104, 86, 7, 146, 85, 113, 102, 105,
	39, 4, 46, 36, 33, 103, 122, 142, 22, 66,
	141, 72, 78, 44, 87, 40, 52, 43, 41, 42,
	68, 34, 37, 38, 16, 76, 15, 24, 25, 17,
	18, 68, 21, 46, 80, 71, 19, 29, 24, 25,
	27, 137, 121, 23, 44, 120, 40, 73, 43, 41,
	42, 27, 100, 93, 94, 101, 46, 80, 91, 92,
	63, 109, 128, 90, 110, 111, 33, 44, 129, 153,
	118, 43, 144, 42, 46, 80, 59, 116, 58, 60,
	61, 114, 104, 126, 124, 44, 62, 40, 115, 105,
	41, 131, 127, 158, 57, 59, 134, 58, 60, 61,
	46, 80, 56, 156, 55, 62, 145, 139, 155, 152,
	138, 44, 151, 149, 147, 148, 125, 140, 54, 53,
	56, 75, 55, 74, 150, 143, 154, 136, 135, 157,
	29, 24, 25, 133, 132, 130, 108, 149, 117, 97,
	95, 89, 88, 82, 81, 77, 70, 69, 47, 96,
	49, 123, 119, 84, 106, 99, 64, 98, 26, 14,
	13, 12, 11, 10, 9, 8, 6, 5, 48, 83,
	51, 3, 2, 1, 45,
}
var yyPact = []int{

	-5, -1000, -1000, 39, 143, -1000, -1000, -1000, -1000, -18,
	-1000, -1000, -1000, -1000, -1000, 14, 161, 165, 16, 128,
	103, -1000, 62, 22, 160, 159, 50, 129, -1000, -1000,
	158, 112, 157, -1000, 156, 169, 6, 155, 154, 72,
	56, 51, -1000, -1000, -1000, 153, -1000, 163, 152, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, 11, -1000, -3, -1000, 149, -1000,
	-1000, -1000, 84, 110, -1000, -1000, -1000, 68, -1000, -1000,
	-1000, 45, -10, 91, -1000, 99, 88, 151, -1000, 112,
	168, -1000, 43, -1000, 40, 1, 167, 130, 143, -1000,
	-1000, -1000, 87, -1000, -1000, -1000, 149, -1000, 66, 148,
	112, 147, -1000, 146, 16, 141, 140, 52, -1000, 120,
	7, 4, 138, 82, 16, -1000, -13, -1000, -1000, -1000,
	68, -1000, 86, -1000, 137, -1000, -1000, 125, 122, -1000,
	79, -1000, -1000, -1000, 16, 121, -1000, 116, 112, -1000,
	-1000, -1000, -1000, -1000, 106, -1000, 112, -1000, -1000,
}
var yyPgo = []int{

	0, 194, 2, 4, 8, 0, 193, 192, 191, 190,
	1, 189, 188, 6, 3, 187, 186, 14, 185, 184,
	183, 182, 181, 180, 179, 28, 178, 177, 176, 175,
	10, 7, 174, 5,
}
var yyR1 = []int{

	0, 6, 6, 9, 10, 10, 11, 11, 12, 12,
	13, 13, 2, 2, 2, 2, 14, 14, 1, 1,
	3, 3, 3, 4, 4, 4, 5, 5, 7, 8,
	8, 15, 15, 15, 15, 15, 15, 15, 15, 15,
	16, 16, 16, 17, 17, 17, 18, 18, 18, 18,
	18, 18, 18, 18, 18, 18, 19, 19, 19, 19,
	19, 20, 20, 21, 21, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 22, 22, 25, 25, 26, 26,
	26, 26, 26, 26, 26, 26, 27, 26, 23, 23,
	23, 28, 28, 31, 31, 32, 32, 33, 33, 30,
	29, 29, 24,
}
var yyR2 = []int{

	0, 1, 1, 1, 0, 1, 0, 1, 0, 1,
	0, 1, 2, 4, 2, 4, 1, 1, 1, 1,
	1, 2, 2, 1, 2, 4, 0, 1, 2, 0,
	2, 1, 1, 1, 1, 1, 1, 1, 1, 1,
	3, 5, 7, 2, 4, 6, 3, 5, 3, 5,
	7, 4, 4, 6, 5, 6, 3, 5, 3, 4,
	6, 3, 4, 2, 1, 2, 2, 2, 3, 5,
	5, 6, 6, 1, 1, 2, 1, 2, 2, 2,
	2, 2, 2, 2, 2, 2, 0, 5, 3, 2,
	3, 1, 2, 2, 2, 1, 2, 2, 2, 2,
	0, 1, 1,
}
var yyChk = []int{

	-1000, -6, -7, -8, 26, -15, -16, -17, -18, -19,
	-20, -21, -22, -23, -24, 7, 5, 10, 11, 17,
	-3, 13, -25, 24, 8, 9, -26, 21, -3, 7,
	29, -2, -14, -4, 27, -13, 9, 28, 29, 6,
	21, 24, 25, 23, 19, -1, 8, 7, -12, 5,
	-10, -9, 20, 11, 10, 14, 12, 11, 14, 12,
	15, 16, 22, 18, -28, -30, 7, -31, 29, 7,
	7, -25, -3, 17, 14, 12, -17, 7, -5, -4,
	9, 7, 7, -11, 4, 10, 7, 28, 7, 7,
	11, 22, 23, 22, 23, 7, 6, 7, -27, -29,
	-31, -30, 7, 28, 15, 22, -32, -33, 7, -14,
	-2, -14, -5, 27, 10, 9, 9, 7, -5, 4,
	22, 22, 25, 4, -13, 6, -3, -33, 16, 22,
	7, -5, 7, 7, -10, 7, 7, 9, 10, 7,
	17, 23, 23, 7, 10, -10, 28, -14, -2, -5,
	7, 7, 7, 10, -10, 7, 7, -5, 7,
}
var yyDef = []int{

	29, -2, 1, 2, 0, 30, 31, 32, 33, 34,
	35, 36, 37, 38, 39, -2, 0, -2, 4, 0,
	0, 73, 74, 0, 0, 0, 76, 0, 28, 20,
	0, 26, 0, 43, 0, 6, 0, 0, 0, 11,
	0, 0, 16, 17, 23, 0, 18, 0, 0, 9,
	63, 5, 3, 65, 66, 79, 82, 67, 78, 81,
	84, 85, 86, 75, 100, 89, 0, 91, 0, 21,
	22, 77, 0, 0, 80, 83, 61, 0, 40, 27,
	19, 26, 46, 0, 7, 0, 56, 0, 58, 26,
	68, 12, 0, 14, 0, 24, 48, 10, 0, 88,
	92, 101, 0, 90, 93, 94, 99, 95, 0, 0,
	26, 0, 44, 0, -2, 0, 0, 59, 62, 0,
	0, 0, 0, 0, -2, 11, 0, 96, 97, 98,
	26, 41, 26, 47, 0, 54, 57, 0, 69, 70,
	0, 13, 15, 25, -2, 0, 87, 0, 26, 45,
	53, 60, 71, 72, 0, 55, 26, 42, 50,
}
var yyTok1 = []int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 20, 3, 23, 27, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 25, 3,
	3, 3, 3, 3, 26, 21, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 22, 3, 3,
	24, 3, 3, 3, 29, 3, 3, 28,
}
var yyTok2 = []int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19,
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
		//line datetime.y:56
		{
			yyVAL.intval = 0
		}
	case 13:
		//line datetime.y:59
		{
			yyVAL.intval = 0
		}
	case 14:
		//line datetime.y:62
		{
			yyVAL.intval = 12
		}
	case 15:
		//line datetime.y:65
		{
			yyVAL.intval = 12
		}
	case 18:
		yyVAL.intval = yyS[yypt-0].intval
	case 19:
		yyVAL.intval = yyS[yypt-0].intval
	case 20:
		yyVAL.tval = yyS[yypt-0].tval
	case 21:
		//line datetime.y:76
		{
			yyS[yypt-0].tval.s = "+" + yyS[yypt-0].tval.s
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 22:
		//line datetime.y:80
		{
			yyS[yypt-0].tval.s = "-" + yyS[yypt-0].tval.s
			yyS[yypt-0].tval.i *= -1
			yyVAL.tval = yyS[yypt-0].tval
		}
	case 23:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 24:
		//line datetime.y:88
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
	case 25:
		//line datetime.y:99
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyS[yypt-3].intval*(3600*yyS[yypt-2].tval.i+60*yyS[yypt-0].tval.i))
		}
	case 26:
		//line datetime.y:104
		{
			yyVAL.zoneval = nil
		}
	case 27:
		yyVAL.zoneval = yyS[yypt-0].zoneval
	case 28:
		//line datetime.y:108
		{
			l := yylex.(*dateLexer)
			if !l.state(HAVE_TIME, true) {
				l.time = time.Unix(int64(yyS[yypt-0].tval.i), 0)
			}
		}
	case 40:
		//line datetime.y:133
		{
			l := yylex.(*dateLexer)
			// Hack to allow HHMMam to parse correctly, cos adie is a mong.
			if yyS[yypt-2].tval.l == 3 || yyS[yypt-2].tval.l == 4 {
				l.setTime(yyS[yypt-2].tval.i/100+yyS[yypt-1].intval, yyS[yypt-2].tval.i%100, 0, yyS[yypt-0].zoneval)
			} else {
				l.setTime(yyS[yypt-2].tval.i+yyS[yypt-1].intval, 0, 0, yyS[yypt-0].zoneval)
			}
		}
	case 41:
		//line datetime.y:142
		{
			yylex.(*dateLexer).setTime(yyS[yypt-4].tval.i+yyS[yypt-1].intval, yyS[yypt-2].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 42:
		//line datetime.y:145
		{
			yylex.(*dateLexer).setTime(yyS[yypt-6].tval.i+yyS[yypt-1].intval, yyS[yypt-4].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].zoneval)
		}
	case 43:
		//line datetime.y:152
		{
			yylex.(*dateLexer).setHMS(yyS[yypt-1].tval.i, yyS[yypt-1].tval.l, yyS[yypt-0].zoneval)
		}
	case 44:
		//line datetime.y:155
		{
			yylex.(*dateLexer).setTime(yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, 0, yyS[yypt-0].zoneval)
		}
	case 45:
		//line datetime.y:158
		{
			yylex.(*dateLexer).setTime(yyS[yypt-5].tval.i, yyS[yypt-3].tval.i, yyS[yypt-1].tval.i, yyS[yypt-0].zoneval)
		}
	case 46:
		//line datetime.y:166
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
	case 47:
		//line datetime.y:176
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
	case 48:
		//line datetime.y:189
		{
			// the DDth
			yylex.(*dateLexer).setDay(yyS[yypt-1].tval.i)
		}
	case 49:
		//line datetime.y:193
		{
			// the DDth of Month
			yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 50:
		//line datetime.y:197
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
	case 51:
		//line datetime.y:210
		{
			// DD[th] [of] Month
			yylex.(*dateLexer).setDate(0, yyS[yypt-0].intval, yyS[yypt-3].tval.i)
		}
	case 52:
		//line datetime.y:214
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
	case 53:
		//line datetime.y:224
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
	case 54:
		//line datetime.y:237
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
	case 55:
		//line datetime.y:251
		{
			l := yylex.(*dateLexer)
			if yyS[yypt-0].tval.l == 4 {
				// assume Month [the] DD[th][,] YYYY
				l.setDate(yyS[yypt-0].tval.i, yyS[yypt-5].intval, yyS[yypt-3].tval.i)
			} else if yyS[yypt-0].tval.i > 68 {
				// assume Month [the] DD[th][,] YY, add 1900 if YY > 68
				l.setDate(yyS[yypt-0].tval.i+1900, yyS[yypt-5].intval, yyS[yypt-3].tval.i)
			} else {
				// assume Month [the] DD[th][,] YY, add 2000 otherwise
				l.setDate(yyS[yypt-0].tval.i+2000, yyS[yypt-5].intval, yyS[yypt-3].tval.i)
			}
		}
	case 56:
		//line datetime.y:267
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
	case 57:
		//line datetime.y:281
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
	case 58:
		//line datetime.y:294
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
	case 59:
		//line datetime.y:304
		{
			// assume YYYY-'W'WW
			yylex.(*dateLexer).setWeek(yyS[yypt-3].tval.i, yyS[yypt-0].tval.i, 1)
		}
	case 60:
		//line datetime.y:308
		{
			// assume YYYY-'W'WW-D
			yylex.(*dateLexer).setWeek(yyS[yypt-5].tval.i, yyS[yypt-2].tval.i, yyS[yypt-0].tval.i)
		}
	case 62:
		//line datetime.y:316
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
	case 63:
		//line datetime.y:331
		{
			// Tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-1].intval, 0)
		}
	case 64:
		//line datetime.y:335
		{
			// March
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, 0)
		}
	case 65:
		//line datetime.y:339
		{
			// Next tuesday
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 66:
		//line datetime.y:343
		{
			// Next march
			yylex.(*dateLexer).setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 67:
		//line datetime.y:347
		{
			// +-N Tuesdays
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-1].tval.i)
		}
	case 68:
		//line datetime.y:351
		{
			// 3rd Tuesday 
			yylex.(*dateLexer).setDays(yyS[yypt-0].intval, yyS[yypt-2].tval.i)
		}
	case 69:
		//line datetime.y:355
		{
			// 3rd Tuesday of (implicit this) March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i)
			l.setMonths(yyS[yypt-0].intval, 1)
		}
	case 70:
		//line datetime.y:361
		{
			// 3rd Tuesday of 2012
			yylex.(*dateLexer).setDays(yyS[yypt-2].intval, yyS[yypt-4].tval.i, yyS[yypt-0].tval.i)
		}
	case 71:
		//line datetime.y:365
		{
			// 3rd Tuesday of March 2012
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-1].intval, 1, yyS[yypt-0].tval.i)
		}
	case 72:
		//line datetime.y:371
		{
			// 3rd Tuesday of next March
			l := yylex.(*dateLexer)
			l.setDays(yyS[yypt-3].intval, yyS[yypt-5].tval.i)
			l.setMonths(yyS[yypt-0].intval, yyS[yypt-1].intval)
		}
	case 73:
		//line datetime.y:377
		{
			// yesterday or tomorrow
			d := time.Now().Weekday()
			yylex.(*dateLexer).setDays((7+int(d)+yyS[yypt-0].intval)%7, yyS[yypt-0].intval)
		}
	case 75:
		//line datetime.y:385
		{
			yylex.(*dateLexer).setAgo()
		}
	case 78:
		//line datetime.y:394
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 79:
		//line datetime.y:397
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].intval)
		}
	case 80:
		//line datetime.y:400
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), 1)
		}
	case 81:
		//line datetime.y:403
		{
			// Special-case to handle "week" and "fortnight"
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].tval.i*yyS[yypt-0].intval)
		}
	case 82:
		//line datetime.y:407
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-1].intval*yyS[yypt-0].intval)
		}
	case 83:
		//line datetime.y:410
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-0].intval)
		}
	case 84:
		//line datetime.y:413
		{
			// As we need to be able to separate out YD from HS in ISO durations
			// this becomes a fair bit messier than if Y D H S were just T_OFFSET
			// Because writing "next y" or "two h" would be odd, disallow
			// T_RELATIVE tokens from being used with ISO single-letter notation
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 85:
		//line datetime.y:420
		{
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 86:
		//line datetime.y:423
		{
			// Resolve 'm' ambiguity in favour of minutes outside ISO duration
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 87:
		//line datetime.y:426
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyS[yypt-4].tval.i*7)
		}
	case 90:
		//line datetime.y:434
		{
			yylex.(*dateLexer).addOffset(O_DAY, 7*yyS[yypt-1].tval.i)
		}
	case 93:
		//line datetime.y:444
		{
			// takes care of Y and D
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 94:
		//line datetime.y:448
		{
			yylex.(*dateLexer).addOffset(O_MONTH, yyS[yypt-1].tval.i)
		}
	case 97:
		//line datetime.y:457
		{
			// takes care of H and S
			yylex.(*dateLexer).addOffset(offset(yyS[yypt-0].intval), yyS[yypt-1].tval.i)
		}
	case 98:
		//line datetime.y:461
		{
			yylex.(*dateLexer).addOffset(O_MIN, yyS[yypt-1].tval.i)
		}
	case 102:
		//line datetime.y:473
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
