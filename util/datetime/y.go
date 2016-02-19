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
const T_MIDTIME = 57353
const T_MONTHNAME = 57354
const T_DAYNAME = 57355
const T_DAYS = 57356
const T_DAYSHIFT = 57357
const T_OFFSET = 57358
const T_ISOYD = 57359
const T_ISOHS = 57360
const T_RELATIVE = 57361
const T_AGO = 57362
const T_ZONE = 57363

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"T_OF",
	"T_THE",
	"T_IGNORE",
	"T_DAYQUAL",
	"T_INTEGER",
	"T_PLUS",
	"T_MINUS",
	"T_MIDTIME",
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
	"','",
	"'A'",
	"'M'",
	"'.'",
	"'P'",
	"':'",
	"'@'",
	"'/'",
	"'W'",
	"'T'",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line datetime.y:502

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 16,
	4, 10,
	12, 10,
	13, 22,
	14, 22,
	16, 22,
	17, 22,
	18, 22,
	24, 22,
	-2, 105,
	-1, 18,
	8, 8,
	-2, 67,
	-1, 116,
	8, 4,
	-2, 54,
	-1, 145,
	8, 4,
	-2, 52,
}

const yyNprod = 107
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 199

var yyAct = [...]int{

	114, 52, 33, 34, 53, 109, 37, 70, 82, 68,
	106, 32, 115, 41, 7, 48, 38, 107, 89, 4,
	143, 124, 88, 105, 142, 35, 104, 46, 123, 42,
	22, 45, 43, 44, 81, 36, 39, 40, 17, 15,
	90, 16, 25, 26, 24, 18, 19, 79, 29, 71,
	59, 61, 20, 60, 62, 63, 28, 69, 74, 23,
	61, 64, 60, 62, 63, 96, 97, 65, 94, 95,
	64, 122, 54, 66, 106, 102, 65, 103, 48, 83,
	71, 107, 138, 93, 111, 48, 83, 112, 113, 35,
	46, 153, 42, 120, 45, 43, 44, 46, 48, 83,
	21, 45, 129, 44, 58, 30, 57, 126, 130, 78,
	46, 77, 42, 132, 128, 43, 145, 116, 135, 118,
	48, 83, 31, 25, 26, 31, 25, 26, 75, 29,
	117, 146, 46, 76, 149, 147, 148, 28, 160, 140,
	56, 55, 58, 139, 57, 159, 156, 154, 155, 157,
	141, 152, 151, 150, 144, 137, 136, 149, 134, 133,
	131, 110, 119, 100, 98, 92, 91, 85, 84, 80,
	73, 72, 49, 127, 99, 51, 125, 121, 87, 108,
	101, 67, 27, 14, 13, 12, 11, 10, 9, 8,
	6, 5, 50, 86, 3, 2, 1, 158, 47,
}
var yyPact = [...]int{

	-9, -1000, -1000, 33, 117, -1000, -1000, -1000, -1000, -20,
	-1000, -1000, -1000, -1000, -1000, -1000, 6, 164, 170, 50,
	128, 37, 53, 49, -1000, 163, 162, 114, 95, -1000,
	-1000, -1000, 161, 111, 160, -1000, 159, 174, 10, 158,
	157, 70, 44, 41, -1000, -1000, -1000, 156, -1000, 167,
	155, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, 18, -1000, -7,
	-1000, 153, -1000, -1000, -1000, 46, 90, -1000, -1000, -1000,
	76, -1000, -1000, -1000, 69, -17, 105, -1000, 120, 109,
	154, -1000, 111, 173, -1000, 47, -1000, 4, -6, 172,
	166, -1000, -1000, -1000, 57, -1000, -1000, -1000, 153, -1000,
	84, 152, 111, 151, -1000, 150, 50, 148, 147, 72,
	-1000, 131, -1, -5, 146, 104, 50, -1000, -1000, -1000,
	-1000, 76, -1000, 89, -1000, 145, -1000, -1000, 144, 143,
	-1000, 79, -1000, -1000, -1000, 50, 140, 138, 111, 137,
	-1000, -1000, -1000, -1000, 130, -1000, 111, -1000, -1000, -1000,
	-1000,
}
var yyPgo = [...]int{

	0, 198, 2, 100, 197, 8, 0, 196, 195, 194,
	4, 1, 193, 192, 6, 3, 191, 190, 14, 189,
	188, 187, 186, 185, 184, 183, 30, 182, 181, 180,
	9, 7, 179, 5,
}
var yyR1 = [...]int{

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
	33, 33, 30, 29, 29, 25, 25,
}
var yyR2 = [...]int{

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
	2, 2, 2, 0, 1, 1, 1,
}
var yyChk = [...]int{

	-1000, -7, -8, -9, 28, -16, -17, -18, -19, -20,
	-21, -22, -23, -24, -25, 6, 8, 5, 12, 13,
	19, -3, -26, 26, 11, 9, 10, -27, 23, 15,
	-3, 8, 31, -2, -15, -5, 29, -14, 10, 30,
	31, 7, 23, 26, 27, 25, 21, -1, 9, 8,
	-13, 5, -11, -10, 22, 13, 12, 16, 14, 13,
	16, 14, 17, 18, 24, 30, 20, -28, -30, 8,
	-31, 31, 8, 8, -26, -3, 19, 16, 14, -18,
	8, -6, -5, 10, 8, 8, -12, 4, 12, 8,
	30, 8, 8, 13, 24, 25, 24, 25, 8, 7,
	8, -29, -31, -30, 8, 30, 17, 24, -32, -33,
	8, -15, -2, -15, -6, 29, 12, 10, 10, 8,
	-6, 4, 24, 24, 27, 4, -14, 7, -33, 18,
	24, 8, -6, 8, 8, -11, 8, 8, 10, 12,
	8, 19, 25, 25, 8, 12, -10, -15, -2, -6,
	8, 8, 8, 12, -11, 8, 8, -6, -4, 8,
	8,
}
var yyDef = [...]int{

	31, -2, 1, 2, 0, 32, 33, 34, 35, 36,
	37, 38, 39, 40, 41, 42, -2, 0, -2, 4,
	0, 0, 76, 0, 106, 0, 0, 78, 0, 90,
	30, 22, 0, 28, 0, 46, 0, 6, 0, 0,
	0, 11, 0, 0, 18, 19, 25, 0, 20, 0,
	0, 9, 66, 5, 3, 68, 69, 81, 84, 70,
	80, 83, 86, 87, 88, 89, 77, 103, 92, 0,
	94, 0, 23, 24, 79, 0, 0, 82, 85, 64,
	0, 43, 29, 21, 28, 49, 0, 7, 0, 59,
	0, 61, 28, 71, 14, 0, 16, 0, 26, 51,
	10, 91, 95, 104, 0, 93, 96, 97, 102, 98,
	0, 0, 28, 0, 47, 0, -2, 0, 0, 62,
	65, 0, 0, 0, 0, 0, 55, 11, 99, 100,
	101, 28, 44, 28, 50, 0, 57, 60, 0, 72,
	73, 0, 15, 17, 27, -2, 0, 0, 28, 12,
	56, 63, 74, 75, 0, 58, 28, 45, 48, 13,
	53,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 22, 3, 25, 29, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 27, 3,
	3, 3, 3, 3, 28, 23, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 24, 3, 3,
	26, 3, 3, 3, 31, 3, 3, 30,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
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

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
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
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
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
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
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
			if yyn < 0 || yyn == yytoken {
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
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
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
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
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
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
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
		yyDollar = yyS[yypt-0 : yypt+1]
		//line datetime.y:54
		{
			yyVAL.tval = textint{}
		}
	case 14:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:57
		{
			yyVAL.intval = 0
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:60
		{
			yyVAL.intval = 0
		}
	case 16:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:63
		{
			yyVAL.intval = 12
		}
	case 17:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:66
		{
			yyVAL.intval = 12
		}
	case 23:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:77
		{
			yyDollar[2].tval.s = "+" + yyDollar[2].tval.s
			yyVAL.tval = yyDollar[2].tval
		}
	case 24:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:81
		{
			yyDollar[2].tval.s = "-" + yyDollar[2].tval.s
			yyDollar[2].tval.i *= -1
			yyVAL.tval = yyDollar[2].tval
		}
	case 26:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:89
		{
			hrs, mins := yyDollar[2].tval.i, 0
			if yyDollar[2].tval.l == 4 {
				hrs, mins = (yyDollar[2].tval.i / 100), (yyDollar[2].tval.i % 100)
			} else if yyDollar[2].tval.l == 2 {
				hrs *= 100
			} else {
				yylex.Error("Invalid timezone offset " + yyDollar[2].tval.s)
			}
			yyVAL.zoneval = time.FixedZone("WTF", yyDollar[1].intval*(3600*hrs+60*mins))
		}
	case 27:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:100
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyDollar[1].intval*(3600*yyDollar[2].tval.i+60*yyDollar[4].tval.i))
		}
	case 28:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line datetime.y:105
		{
			yyVAL.zoneval = nil
		}
	case 30:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:109
		{
			yylex.(*dateLexer).setUnix(int64(yyDollar[2].tval.i))
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:132
		{
			l := yylex.(*dateLexer)
			// Hack to allow HHMMam to parse correctly, cos adie is a mong.
			if yyDollar[1].tval.l == 3 || yyDollar[1].tval.l == 4 {
				l.setTime(ampm(yyDollar[1].tval.i/100, yyDollar[2].intval), yyDollar[1].tval.i%100, 0, yyDollar[3].zoneval)
			} else {
				l.setTime(ampm(yyDollar[1].tval.i, yyDollar[2].intval), 0, 0, yyDollar[3].zoneval)
			}
		}
	case 44:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:141
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i+yyDollar[4].intval, yyDollar[3].tval.i, 0, yyDollar[5].zoneval)
		}
	case 45:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line datetime.y:144
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i+yyDollar[6].intval, yyDollar[3].tval.i, yyDollar[5].tval.i, yyDollar[7].zoneval)
		}
	case 46:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:151
		{
			yylex.(*dateLexer).setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, yyDollar[2].zoneval)
		}
	case 47:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:154
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i, yyDollar[3].tval.i, 0, yyDollar[4].zoneval)
		}
	case 48:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line datetime.y:157
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i, yyDollar[3].tval.i, yyDollar[5].tval.i, yyDollar[6].zoneval)
			// Hack to make time.ANSIC, time.UnixDate and time.RubyDate parse
			if yyDollar[7].tval.l == 4 {
				yylex.(*dateLexer).setYear(yyDollar[7].tval.i)
			}
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:169
		{
			l := yylex.(*dateLexer)
			if yyDollar[3].tval.l == 4 {
				// assume we have MM/YYYY
				l.setDate(yyDollar[3].tval.i, yyDollar[1].tval.i, 1)
			} else {
				// assume we have DD/MM (too bad, americans)
				l.setDate(0, yyDollar[3].tval.i, yyDollar[1].tval.i)
			}
		}
	case 50:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:179
		{
			l := yylex.(*dateLexer)
			if yyDollar[5].tval.l == 4 {
				// assume we have DD/MM/YYYY
				l.setDate(yyDollar[5].tval.i, yyDollar[3].tval.i, yyDollar[1].tval.i)
			} else if yyDollar[5].tval.i > 68 {
				// assume we have DD/MM/YY, add 1900 if YY > 68
				l.setDate(yyDollar[5].tval.i+1900, yyDollar[3].tval.i, yyDollar[1].tval.i)
			} else {
				// assume we have DD/MM/YY, add 2000 otherwise
				l.setDate(yyDollar[5].tval.i+2000, yyDollar[3].tval.i, yyDollar[1].tval.i)
			}
		}
	case 51:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:192
		{
			// the DDth
			yylex.(*dateLexer).setDay(yyDollar[2].tval.i)
		}
	case 52:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:196
		{
			// the DDth of Month
			yylex.(*dateLexer).setDate(0, yyDollar[5].intval, yyDollar[2].tval.i)
		}
	case 53:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line datetime.y:200
		{
			l := yylex.(*dateLexer)
			if yyDollar[7].tval.l == 4 {
				// the DDth of Month[,] YYYY
				l.setDate(yyDollar[7].tval.i, yyDollar[5].intval, yyDollar[2].tval.i)
			} else if yyDollar[7].tval.i > 68 {
				// the DDth of Month[,] YY, add 1900 if YY > 68
				l.setDate(yyDollar[7].tval.i+1900, yyDollar[5].intval, yyDollar[2].tval.i)
			} else {
				// the DDth of Month[,] YY, add 2000 otherwise
				l.setDate(yyDollar[7].tval.i+2000, yyDollar[5].intval, yyDollar[2].tval.i)
			}
		}
	case 54:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:213
		{
			// DD[th] [of] Month
			yylex.(*dateLexer).setDate(0, yyDollar[4].intval, yyDollar[1].tval.i)
		}
	case 55:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:217
		{
			l := yylex.(*dateLexer)
			if yyDollar[3].tval.l == 4 {
				// assume Month YYYY
				l.setDate(yyDollar[3].tval.i, yyDollar[1].intval, 1)
			} else {
				// assume Month [the] DD[th]
				l.setDate(0, yyDollar[1].intval, yyDollar[3].tval.i)
			}
		}
	case 56:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:227
		{
			l := yylex.(*dateLexer)
			if yyDollar[6].tval.l == 4 {
				// assume DD[th] [of] Month[,] YYYY
				l.setDate(yyDollar[6].tval.i, yyDollar[4].intval, yyDollar[1].tval.i)
			} else if yyDollar[6].tval.i > 68 {
				// assume DD[th] [of] Month[,] YY, add 1900 if YY > 68
				l.setDate(yyDollar[6].tval.i+1900, yyDollar[4].intval, yyDollar[1].tval.i)
			} else {
				// assume DD[th] [of] Month[,] YY, add 2000 otherwise
				l.setDate(yyDollar[6].tval.i+2000, yyDollar[4].intval, yyDollar[1].tval.i)
			}
		}
	case 57:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:240
		{
			// RFC 850, srsly :(
			l := yylex.(*dateLexer)
			if yyDollar[5].tval.l == 4 {
				// assume DD-Mon-YYYY
				l.setDate(yyDollar[5].tval.i, yyDollar[3].intval, yyDollar[1].tval.i)
			} else if yyDollar[5].tval.i > 68 {
				// assume DD-Mon-YY, add 1900 if YY > 68
				l.setDate(yyDollar[5].tval.i+1900, yyDollar[3].intval, yyDollar[1].tval.i)
			} else {
				// assume DD-Mon-YY, add 2000 otherwise
				l.setDate(yyDollar[5].tval.i+2000, yyDollar[3].intval, yyDollar[1].tval.i)
			}
		}
	case 58:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:254
		{
			// comma cannot be optional here; T_MONTHNAME T_INTEGER T_INTEGER
			// can easily match [March 02 10]:30:00 and break parsing.
			l := yylex.(*dateLexer)
			if yyDollar[6].tval.l == 4 {
				// assume Month [the] DD[th], YYYY
				l.setDate(yyDollar[6].tval.i, yyDollar[1].intval, yyDollar[3].tval.i)
			} else if yyDollar[6].tval.i > 68 {
				// assume Month [the] DD[th], YY, add 1900 if YY > 68
				l.setDate(yyDollar[6].tval.i+1900, yyDollar[1].intval, yyDollar[3].tval.i)
			} else {
				// assume Month [the] DD[th], YY, add 2000 otherwise
				l.setDate(yyDollar[6].tval.i+2000, yyDollar[1].intval, yyDollar[3].tval.i)
			}
		}
	case 59:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:272
		{
			l := yylex.(*dateLexer)
			if yyDollar[1].tval.l == 4 && yyDollar[3].tval.l == 3 {
				// assume we have YYYY-DDD
				l.setDate(yyDollar[1].tval.i, 1, yyDollar[3].tval.i)
			} else if yyDollar[1].tval.l == 4 {
				// assume we have YYYY-MM
				l.setDate(yyDollar[1].tval.i, yyDollar[3].tval.i, 1)
			} else {
				// assume we have MM-DD (not strictly ISO compliant)
				// this is for americans, because of DD/MM above ;-)
				l.setDate(0, yyDollar[1].tval.i, yyDollar[3].tval.i)
			}
		}
	case 60:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:286
		{
			l := yylex.(*dateLexer)
			if yyDollar[1].tval.l == 4 {
				// assume we have YYYY-MM-DD
				l.setDate(yyDollar[1].tval.i, yyDollar[3].tval.i, yyDollar[5].tval.i)
			} else if yyDollar[1].tval.i > 68 {
				// assume we have YY-MM-DD, add 1900 if YY > 68
				l.setDate(yyDollar[1].tval.i+1900, yyDollar[3].tval.i, yyDollar[5].tval.i)
			} else {
				// assume we have YY-MM-DD, add 2000 otherwise
				l.setDate(yyDollar[1].tval.i+2000, yyDollar[3].tval.i, yyDollar[5].tval.i)
			}
		}
	case 61:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:299
		{
			l := yylex.(*dateLexer)
			wday, week := 1, yyDollar[3].tval.i
			if yyDollar[3].tval.l == 3 {
				// assume YYYY'W'WWD
				wday = week % 10
				week = week / 10
			}
			l.setWeek(yyDollar[1].tval.i, week, wday)
		}
	case 62:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:309
		{
			// assume YYYY-'W'WW
			yylex.(*dateLexer).setWeek(yyDollar[1].tval.i, yyDollar[4].tval.i, 1)
		}
	case 63:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:313
		{
			// assume YYYY-'W'WW-D
			yylex.(*dateLexer).setWeek(yyDollar[1].tval.i, yyDollar[4].tval.i, yyDollar[6].tval.i)
		}
	case 65:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:321
		{
			// this goes here because the YYYYMMDD and HHMMSS forms of the
			// ISO 8601 format date and time are handled by 'integer' below.
			l := yylex.(*dateLexer)
			if yyDollar[1].tval.l == 8 {
				// assume ISO 8601 YYYYMMDD
				l.setYMD(yyDollar[1].tval.i, yyDollar[1].tval.l)
			} else if yyDollar[1].tval.l == 7 {
				// assume ISO 8601 ordinal YYYYDDD
				l.setDate(yyDollar[1].tval.i/1000, 1, yyDollar[1].tval.i%1000)
			}
			l.setHMS(yyDollar[3].tval.i, yyDollar[3].tval.l, yyDollar[4].zoneval)
		}
	case 66:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:336
		{
			// Tuesday
			yylex.(*dateLexer).setDays(yyDollar[1].intval, 0)
		}
	case 67:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:340
		{
			// March
			yylex.(*dateLexer).setMonths(yyDollar[1].intval, 0)
		}
	case 68:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:344
		{
			// Next tuesday
			yylex.(*dateLexer).setDays(yyDollar[2].intval, yyDollar[1].intval)
		}
	case 69:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:348
		{
			// Next march
			yylex.(*dateLexer).setMonths(yyDollar[2].intval, yyDollar[1].intval)
		}
	case 70:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:352
		{
			// +-N Tuesdays
			yylex.(*dateLexer).setDays(yyDollar[2].intval, yyDollar[1].tval.i)
		}
	case 71:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:356
		{
			// 3rd Tuesday (of implicit this month)
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(0, 0)
		}
	case 72:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:362
		{
			// 3rd Tuesday of (implicit this) March
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(yyDollar[5].intval, 0)
		}
	case 73:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:368
		{
			// 3rd Tuesday of 2012
			yylex.(*dateLexer).setDays(yyDollar[3].intval, yyDollar[1].tval.i, yyDollar[5].tval.i)
		}
	case 74:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:372
		{
			// 3rd Tuesday of March 2012
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(yyDollar[5].intval, 0, yyDollar[6].tval.i)
		}
	case 75:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:378
		{
			// 3rd Tuesday of next March
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(yyDollar[6].intval, yyDollar[5].intval)
		}
	case 77:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:387
		{
			yylex.(*dateLexer).setAgo()
		}
	case 80:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:396
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 81:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:399
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].intval)
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:402
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), 1)
		}
	case 83:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:405
		{
			// Special-case to handle "week" and "fortnight"
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].tval.i*yyDollar[2].intval)
		}
	case 84:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:409
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].intval*yyDollar[2].intval)
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:412
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[2].intval)
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:415
		{
			// As we need to be able to separate out YD from HS in ISO durations
			// this becomes a fair bit messier than if Y D H S were just T_OFFSET
			// Because writing "next y" or "two h" would be odd, disallow
			// T_RELATIVE tokens from being used with ISO single-letter notation
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:422
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:425
		{
			// Resolve 'm' ambiguity in favour of minutes outside ISO duration
			yylex.(*dateLexer).addOffset(O_MIN, yyDollar[1].tval.i)
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:429
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].tval.i*7)
		}
	case 90:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:432
		{
			// yesterday or tomorrow
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].intval)
		}
	case 93:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:441
		{
			yylex.(*dateLexer).addOffset(O_DAY, 7*yyDollar[2].tval.i)
		}
	case 96:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:451
		{
			// takes care of Y and D
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 97:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:455
		{
			yylex.(*dateLexer).addOffset(O_MONTH, yyDollar[1].tval.i)
		}
	case 100:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:464
		{
			// takes care of H and S
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 101:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:468
		{
			yylex.(*dateLexer).addOffset(O_MIN, yyDollar[1].tval.i)
		}
	case 105:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:480
		{
			l := yylex.(*dateLexer)
			if yyDollar[1].tval.l == 8 {
				// assume ISO 8601 YYYYMMDD
				l.setYMD(yyDollar[1].tval.i, yyDollar[1].tval.l)
			} else if yyDollar[1].tval.l == 7 {
				// assume ISO 8601 ordinal YYYYDDD
				l.setDate(yyDollar[1].tval.i/1000, 1, yyDollar[1].tval.i%1000)
			} else if yyDollar[1].tval.l == 6 {
				// assume ISO 8601 HHMMSS with no zone
				l.setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, nil)
			} else if yyDollar[1].tval.l == 4 {
				// Assuming HHMM because that's more useful on IRC.
				l.setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, nil)
			} else if yyDollar[1].tval.l == 2 {
				// assume HH with no zone
				l.setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, nil)
			}
		}
	case 106:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:499
		{
			yylex.(*dateLexer).setHMS(yyDollar[1].intval, 2, nil)
		}
	}
	goto yystack /* stack new state and value */
}
