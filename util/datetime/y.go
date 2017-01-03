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
const T_SECOND = 57350
const T_INTEGER = 57351
const T_PLUS = 57352
const T_MINUS = 57353
const T_MIDTIME = 57354
const T_MONTHNAME = 57355
const T_DAYNAME = 57356
const T_DAYS = 57357
const T_DAYSHIFT = 57358
const T_OFFSET = 57359
const T_ISOYD = 57360
const T_ISOHS = 57361
const T_RELATIVE = 57362
const T_AGO = 57363
const T_ZONE = 57364

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"T_OF",
	"T_THE",
	"T_IGNORE",
	"T_DAYQUAL",
	"T_SECOND",
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

//line datetime.y:511

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 16,
	4, 10,
	8, 27,
	13, 10,
	14, 27,
	15, 27,
	17, 27,
	18, 27,
	19, 27,
	25, 27,
	-2, 110,
	-1, 18,
	9, 8,
	-2, 72,
	-1, 121,
	9, 4,
	-2, 59,
	-1, 150,
	9, 4,
	-2, 57,
}

const yyNprod = 112
const yyPrivate = 57344

var yyTokenNames []string
var yyStates []string

const yyLast = 216

var yyAct = [...]int{

	119, 54, 35, 36, 55, 20, 114, 39, 87, 74,
	59, 72, 34, 120, 109, 111, 7, 94, 4, 129,
	148, 93, 112, 43, 147, 37, 50, 40, 110, 101,
	102, 128, 64, 56, 127, 80, 86, 75, 48, 95,
	44, 82, 47, 45, 46, 22, 38, 41, 42, 17,
	15, 84, 26, 16, 27, 28, 24, 18, 19, 73,
	31, 21, 70, 62, 25, 158, 32, 98, 30, 63,
	65, 23, 61, 66, 67, 78, 99, 100, 150, 121,
	68, 107, 75, 108, 50, 88, 69, 165, 143, 116,
	64, 79, 117, 118, 37, 134, 48, 111, 125, 62,
	47, 135, 46, 123, 112, 122, 65, 164, 61, 66,
	67, 161, 160, 131, 50, 88, 68, 157, 137, 156,
	133, 155, 69, 140, 149, 142, 48, 33, 27, 28,
	50, 88, 146, 141, 139, 138, 151, 136, 115, 154,
	152, 153, 48, 124, 44, 105, 47, 45, 46, 50,
	88, 103, 159, 81, 162, 97, 26, 33, 27, 28,
	96, 48, 154, 44, 31, 90, 45, 81, 25, 62,
	26, 145, 30, 89, 62, 144, 60, 62, 61, 58,
	57, 60, 25, 61, 83, 85, 61, 77, 76, 51,
	132, 104, 53, 130, 126, 92, 113, 106, 71, 29,
	14, 13, 12, 11, 10, 9, 8, 6, 5, 52,
	91, 3, 2, 1, 163, 49,
}
var yyPact = [...]int{

	-11, -1000, -1000, 44, 118, -1000, -1000, -1000, -1000, -20,
	-1000, -1000, -1000, -1000, -1000, -1000, 16, 180, 187, 10,
	166, 55, 41, 50, -1000, -1000, -1000, 179, 178, 148,
	169, -1000, -1000, -1000, 176, 104, 164, -1000, 156, 191,
	8, 151, 146, 53, 51, 4, -1000, -1000, -1000, 142,
	-1000, 184, 136, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, 5, -1000, -3, -1000, 129, -1000, -1000, -1000, 91,
	161, -1000, -1000, -1000, -1000, 74, -1000, -1000, -1000, 120,
	-17, 66, -1000, 94, 92, 134, -1000, 104, 190, -1000,
	9, -1000, 6, -9, 189, 183, -1000, -1000, -1000, 79,
	-1000, -1000, -1000, 129, -1000, 76, 128, 104, 126, -1000,
	125, 10, 124, 116, 77, -1000, 162, -2, -6, 115,
	65, 10, -1000, -1000, -1000, -1000, 74, -1000, 139, -1000,
	112, -1000, -1000, 110, 108, -1000, 52, -1000, -1000, -1000,
	10, 103, 102, 104, 98, -1000, -1000, -1000, -1000, 78,
	-1000, 104, -1000, -1000, -1000, -1000,
}
var yyPgo = [...]int{

	0, 215, 2, 10, 5, 61, 214, 8, 0, 213,
	212, 211, 4, 1, 210, 209, 7, 3, 208, 207,
	16, 206, 205, 204, 203, 202, 201, 200, 45, 199,
	198, 197, 11, 9, 196, 6,
}
var yyR1 = [...]int{

	0, 9, 9, 12, 13, 13, 14, 14, 15, 15,
	16, 16, 6, 6, 3, 3, 4, 4, 4, 2,
	2, 2, 2, 17, 17, 1, 1, 5, 5, 5,
	7, 7, 7, 8, 8, 10, 11, 11, 18, 18,
	18, 18, 18, 18, 18, 18, 18, 18, 19, 19,
	19, 20, 20, 20, 21, 21, 21, 21, 21, 21,
	21, 21, 21, 21, 22, 22, 22, 22, 22, 23,
	23, 24, 24, 24, 24, 24, 24, 24, 24, 24,
	24, 25, 25, 28, 28, 29, 29, 29, 29, 29,
	29, 29, 29, 29, 29, 29, 26, 26, 26, 30,
	30, 33, 33, 34, 34, 35, 35, 32, 31, 31,
	27, 27,
}
var yyR2 = [...]int{

	0, 1, 1, 1, 0, 1, 0, 1, 0, 1,
	0, 1, 0, 1, 1, 1, 1, 1, 1, 2,
	4, 2, 4, 1, 1, 1, 1, 1, 2, 2,
	1, 2, 4, 0, 1, 2, 0, 2, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 3, 5,
	7, 2, 4, 7, 3, 5, 3, 5, 7, 4,
	4, 6, 5, 6, 3, 5, 3, 4, 6, 3,
	4, 2, 1, 2, 2, 2, 3, 5, 5, 6,
	6, 1, 2, 1, 2, 2, 2, 2, 2, 2,
	2, 2, 2, 2, 2, 1, 3, 2, 3, 1,
	2, 2, 2, 1, 2, 2, 2, 2, 0, 1,
	1, 1,
}
var yyChk = [...]int{

	-1000, -9, -10, -11, 29, -18, -19, -20, -21, -22,
	-23, -24, -25, -26, -27, 6, 9, 5, 13, 14,
	-4, -5, -28, 27, 12, 20, 8, 10, 11, -29,
	24, 16, -5, 9, 32, -2, -17, -7, 30, -16,
	11, 31, 32, 7, 24, 27, 28, 26, 22, -1,
	10, 9, -15, 5, -13, -12, 23, 14, 13, -3,
	15, 17, 8, 14, -3, 15, 18, 19, 25, 31,
	21, -30, -32, 9, -33, 32, 9, 9, -28, -5,
	-4, 5, -3, 15, -20, 9, -8, -7, 11, 9,
	9, -14, 4, 13, 9, 31, 9, 9, 14, 25,
	26, 25, 26, 9, 7, 9, -31, -33, -32, 9,
	31, 18, 25, -34, -35, 9, -17, -2, -17, -8,
	30, 13, 11, 11, 9, -8, 4, 25, 25, 28,
	4, -16, 7, -35, 19, 25, 9, -8, 9, 9,
	-13, 9, 9, 11, 13, 9, -4, 26, 26, 9,
	13, -12, -17, -2, -8, 9, 9, 9, 13, -13,
	9, 9, -8, -6, 9, 9,
}
var yyDef = [...]int{

	36, -2, 1, 2, 0, 37, 38, 39, 40, 41,
	42, 43, 44, 45, 46, 47, -2, 17, -2, 4,
	0, 0, 81, 0, 111, 16, 18, 0, 0, 83,
	0, 95, 35, 27, 0, 33, 0, 51, 0, 6,
	0, 0, 0, 11, 0, 0, 23, 24, 30, 0,
	25, 0, 0, 9, 71, 5, 3, 73, 74, 86,
	89, 14, 15, 75, 85, 88, 91, 92, 93, 94,
	82, 108, 97, 0, 99, 0, 28, 29, 84, 0,
	0, 17, 87, 90, 69, 0, 48, 34, 26, 33,
	54, 0, 7, 0, 64, 0, 66, 33, 76, 19,
	0, 21, 0, 31, 56, 10, 96, 100, 109, 0,
	98, 101, 102, 107, 103, 0, 0, 33, 0, 52,
	0, -2, 0, 0, 67, 70, 0, 0, 0, 0,
	0, 60, 11, 104, 105, 106, 33, 49, 33, 55,
	0, 62, 65, 0, 77, 78, 0, 20, 22, 32,
	-2, 0, 0, 33, 12, 61, 68, 79, 80, 0,
	63, 33, 50, 53, 13, 58,
}
var yyTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 23, 3, 26, 30, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 28, 3,
	3, 3, 3, 3, 29, 24, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 25, 3, 3,
	27, 3, 3, 3, 32, 3, 3, 31,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22,
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
	case 15:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:57
		{
			yyVAL.intval = int(O_SEC)
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:62
		{
			yyVAL.intval = 0
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:62
		{
			yyVAL.intval = 2
		}
	case 19:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:65
		{
			yyVAL.intval = 0
		}
	case 20:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:68
		{
			yyVAL.intval = 0
		}
	case 21:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:71
		{
			yyVAL.intval = 12
		}
	case 22:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:74
		{
			yyVAL.intval = 12
		}
	case 28:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:85
		{
			yyDollar[2].tval.s = "+" + yyDollar[2].tval.s
			yyVAL.tval = yyDollar[2].tval
		}
	case 29:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:89
		{
			yyDollar[2].tval.s = "-" + yyDollar[2].tval.s
			yyDollar[2].tval.i *= -1
			yyVAL.tval = yyDollar[2].tval
		}
	case 31:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:97
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
	case 32:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:108
		{
			yyVAL.zoneval = time.FixedZone("WTF", yyDollar[1].intval*(3600*yyDollar[2].tval.i+60*yyDollar[4].tval.i))
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line datetime.y:113
		{
			yyVAL.zoneval = nil
		}
	case 35:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:117
		{
			yylex.(*dateLexer).setUnix(int64(yyDollar[2].tval.i))
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:140
		{
			l := yylex.(*dateLexer)
			// Hack to allow HHMMam to parse correctly, cos adie is a mong.
			if yyDollar[1].tval.l == 3 || yyDollar[1].tval.l == 4 {
				l.setTime(ampm(yyDollar[1].tval.i/100, yyDollar[2].intval), yyDollar[1].tval.i%100, 0, yyDollar[3].zoneval)
			} else {
				l.setTime(ampm(yyDollar[1].tval.i, yyDollar[2].intval), 0, 0, yyDollar[3].zoneval)
			}
		}
	case 49:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:149
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i+yyDollar[4].intval, yyDollar[3].tval.i, 0, yyDollar[5].zoneval)
		}
	case 50:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line datetime.y:152
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i+yyDollar[6].intval, yyDollar[3].tval.i, yyDollar[5].tval.i, yyDollar[7].zoneval)
		}
	case 51:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:159
		{
			yylex.(*dateLexer).setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, yyDollar[2].zoneval)
		}
	case 52:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:162
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i, yyDollar[3].tval.i, 0, yyDollar[4].zoneval)
		}
	case 53:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line datetime.y:165
		{
			yylex.(*dateLexer).setTime(yyDollar[1].tval.i, yyDollar[3].tval.i, yyDollar[5].tval.i, yyDollar[6].zoneval)
			// Hack to make time.ANSIC, time.UnixDate and time.RubyDate parse
			if yyDollar[7].tval.l == 4 {
				yylex.(*dateLexer).setYear(yyDollar[7].tval.i)
			}
		}
	case 54:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:177
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
	case 55:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:187
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
	case 56:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:200
		{
			// the DDth
			yylex.(*dateLexer).setDay(yyDollar[2].tval.i)
		}
	case 57:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:204
		{
			// the DDth of Month
			yylex.(*dateLexer).setDate(0, yyDollar[5].intval, yyDollar[2].tval.i)
		}
	case 58:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line datetime.y:208
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
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:221
		{
			// DD[th] [of] Month
			yylex.(*dateLexer).setDate(0, yyDollar[4].intval, yyDollar[1].tval.i)
		}
	case 60:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:225
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
	case 61:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:235
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
	case 62:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:248
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
	case 63:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:262
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
	case 64:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:280
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
	case 65:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:294
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
	case 66:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:307
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
	case 67:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:317
		{
			// assume YYYY-'W'WW
			yylex.(*dateLexer).setWeek(yyDollar[1].tval.i, yyDollar[4].tval.i, 1)
		}
	case 68:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:321
		{
			// assume YYYY-'W'WW-D
			yylex.(*dateLexer).setWeek(yyDollar[1].tval.i, yyDollar[4].tval.i, yyDollar[6].tval.i)
		}
	case 70:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line datetime.y:329
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
	case 71:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:344
		{
			// Tuesday
			yylex.(*dateLexer).setDays(yyDollar[1].intval, 0)
		}
	case 72:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:348
		{
			// March
			yylex.(*dateLexer).setMonths(yyDollar[1].intval, 0)
		}
	case 73:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:352
		{
			// Next tuesday
			yylex.(*dateLexer).setDays(yyDollar[2].intval, yyDollar[1].intval)
		}
	case 74:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:356
		{
			// Next march
			yylex.(*dateLexer).setMonths(yyDollar[2].intval, yyDollar[1].intval)
		}
	case 75:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:360
		{
			// +-N Tuesdays
			yylex.(*dateLexer).setDays(yyDollar[2].intval, yyDollar[1].tval.i)
		}
	case 76:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:364
		{
			// 3rd Tuesday (of implicit this month)
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(0, 0)
		}
	case 77:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:370
		{
			// 3rd Tuesday of (implicit this) March
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(yyDollar[5].intval, 0)
		}
	case 78:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line datetime.y:376
		{
			// 3rd Tuesday of 2012
			yylex.(*dateLexer).setDays(yyDollar[3].intval, yyDollar[1].tval.i, yyDollar[5].tval.i)
		}
	case 79:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:380
		{
			// 3rd Tuesday of March 2012
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(yyDollar[5].intval, 0, yyDollar[6].tval.i)
		}
	case 80:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line datetime.y:386
		{
			// 3rd Tuesday of next March
			l := yylex.(*dateLexer)
			l.setDays(yyDollar[3].intval, yyDollar[1].tval.i)
			l.setMonths(yyDollar[6].intval, yyDollar[5].intval)
		}
	case 82:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:395
		{
			yylex.(*dateLexer).setAgo()
		}
	case 85:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:404
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 86:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:407
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].intval)
		}
	case 87:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:410
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), 1)
		}
	case 88:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:413
		{
			// Special-case to handle "week" and "fortnight"
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].tval.i*yyDollar[2].intval)
		}
	case 89:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:417
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].intval*yyDollar[2].intval)
		}
	case 90:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:420
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[2].intval)
		}
	case 91:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:423
		{
			// As we need to be able to separate out YD from HS in ISO durations
			// this becomes a fair bit messier than if Y D H S were just T_OFFSET
			// Because writing "next y" or "two h" would be odd, disallow
			// T_RELATIVE tokens from being used with ISO single-letter notation
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 92:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:430
		{
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 93:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:433
		{
			// Resolve 'm' ambiguity in favour of minutes outside ISO duration
			yylex.(*dateLexer).addOffset(O_MIN, yyDollar[1].tval.i)
		}
	case 94:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:437
		{
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].tval.i*7)
		}
	case 95:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:440
		{
			// yesterday or tomorrow
			yylex.(*dateLexer).addOffset(O_DAY, yyDollar[1].intval)
		}
	case 98:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line datetime.y:449
		{
			yylex.(*dateLexer).addOffset(O_DAY, 7*yyDollar[2].tval.i)
		}
	case 101:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:459
		{
			// takes care of Y and D
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 102:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:463
		{
			yylex.(*dateLexer).addOffset(O_MONTH, yyDollar[1].tval.i)
		}
	case 105:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:472
		{
			// takes care of H and S
			yylex.(*dateLexer).addOffset(offset(yyDollar[2].intval), yyDollar[1].tval.i)
		}
	case 106:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line datetime.y:476
		{
			yylex.(*dateLexer).addOffset(O_MIN, yyDollar[1].tval.i)
		}
	case 110:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:488
		{
			l := yylex.(*dateLexer)
			switch yyDollar[1].tval.l {
			case 8:
				// assume ISO 8601 YYYYMMDD
				l.setYMD(yyDollar[1].tval.i, yyDollar[1].tval.l)
			case 7:
				// assume ISO 8601 ordinal YYYYDDD
				l.setDate(yyDollar[1].tval.i/1000, 1, yyDollar[1].tval.i%1000)
			case 6:
				// assume ISO 8601 HHMMSS with no zone
				l.setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, nil)
			case 4:
				// Assuming HHMM because that's more useful on IRC.
				l.setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, nil)
			case 2:
				// assume HH with no zone
				l.setHMS(yyDollar[1].tval.i, yyDollar[1].tval.l, nil)
			}
		}
	case 111:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line datetime.y:508
		{
			yylex.(*dateLexer).setHMS(yyDollar[1].intval, 2, nil)
		}
	}
	goto yystack /* stack new state and value */
}
