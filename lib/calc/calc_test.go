package calc

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"unicode"
)

func TestCalc(t *testing.T) {
	tests := []struct {
		i string
		o float64
		e bool
	}{
		{"min(3+4, 2*(1+2))", 6, false},
		{"2+4", 6, false},
		{"answer", 42, false},
		{"pi*e", math.Pi * math.E, false},
		{"2+", 0, true},
	}
	for i, tc := range tests {
		r, err := Calc(tc.i)
		if (err == nil) == tc.e {
			t.Errorf("Bad Calc error state for %d (err=%v)", i, err)
		}
		// Stupid approximate floats.
		if r-tc.o >= 1e-12 {
			t.Errorf("Bad Calc result for %d, expected: %g, got: %g", i, tc.o, r)
		}
	}
}

func makets(data ...float64) *tokenStack {
	ts := ts(len(data))
	for _, n := range data {
		ts.push(&token{T_NUM, "", n})
	}
	return ts
}

func TestFunctionisers(t *testing.T) {
	f1 := func(x float64) float64 { return 0.5 * x }
	f2 := func(x, y float64) float64 { return 2 * x / y }

	ff1 := functionise1(f1)
	ff2 := functionise2(f2)

	// Exercise f1 with a range of good inputs.
	for i, v := range []float64{2, 1e6, 6.75e-3, math.Pi} {
		ts := makets(v)
		if err := ff1.exec(ts); err != nil {
			t.Errorf("f1(%d): Unexpected error result from function.", i)
		} else if tok, err := ts.pop(); err != nil {
			t.Errorf("f1(%d): Stack not updated correctly with result", i)
		} else if exp := f1(v); tok.numval != exp {
			t.Errorf("f1(%d): Func result differed: %f != %f", i, tok.numval, exp)
		}
	}

	// Same for f2
	for i, v := range [][]float64{{2, 4}, {1e6, 2.9e4}, {6.75, 48.1}, {math.Pi, math.E}} {
		ts := makets(v...)
		if err := ff2.exec(ts); err != nil {
			t.Errorf("f2(%d): Unexpected error result from function.", i)
		} else if tok, err := ts.pop(); err != nil {
			t.Errorf("f2(%d): Stack not updated correctly with result", i)
		} else if exp := f2(v[0], v[1]); tok.numval != exp {
			t.Errorf("f2(%d): Func result differed: %f != %f", i, tok.numval, exp)
		}
	}

	// Now try to break things
	zerots := ts(0)
	if err := ff1.exec(zerots); err == nil {
		t.Errorf("ff1 seemed happy with empty stack")
	}
	if err := ff2.exec(zerots); err == nil {
		t.Errorf("ff2 seemed happy with empty stack")
	}
	onets := makets(1)
	if err := ff2.exec(onets); err == nil {
		t.Errorf("ff2 seemed happy with undersized stack")
	}

	// Lastly, check they're not breaking the stack
	fourts := makets(1, 2, 3, 4)
	ff1.exec(fourts) // should pop 4, push 2
	if len(*fourts) != 4 || (*fourts)[3].numval != 2 || (*fourts)[2].numval != 3 {
		t.Errorf("ff1 changed stack size unexpectedly")
	}
	ff2.exec(fourts) // should pop 2,3, push 3
	if len(*fourts) != 3 || (*fourts)[2].numval != 3 || (*fourts)[1].numval != 2 {
		t.Errorf("ff2 changed stack size unexpectedly")
	}
}

// Lexer tests

// This takes care of exercising peek, next, scan, and rewind
func TestLexerLowLevelFuncs(t *testing.T) {
	// Mmm. Unicodey. 42 bytes, 22 chars.
	l := &lexer{input: "This ïş æ 💩♽⛤ “Ŧəšţ”™."}

	// First, peek.
	if l.peek() != 'T' {
		t.Errorf("Lexer appears to be starting in the wrong place")
	}
	// Advance 5 bytes to ï (\u00ef, 0xC3 0xAf)
	l.pos += 5
	if l.peek() != 'ï' {
		t.Errorf("Lexer not decoding two-byte unicode chars")
	}
	// Advance another byte to the middle of ï
	l.pos += 1
	if l.peek() != 0 {
		t.Errorf("Didn't get EOF from bad unicode")
	}

	// Advance to POO, PILE OF!
	l.pos = strings.Index(l.input, "💩")
	if l.peek() != '💩' {
		t.Errorf("Lexer can't decode shit")
	}

	// For the next three chars, make sure peek() and next() are in sync
	for i := 0; i < 3; i++ {
		if string(l.peek()) != l.next() {
			t.Errorf("Peek and next don't agree")
		}
	}

	// We should be at the space before “Ŧəšţ” now.
	if l.next() != " " {
		t.Errorf("Lexer seems out of sync with reality")
	}
	l.rewind()
	if l.next() != " " {
		t.Errorf("Lexer still seems out of sync with reality")
	}
	l.next() // skip opening quote.

	// Test scanning Ŧəšţ
	if l.scan(unicode.IsLetter) != "Ŧəšţ" {
		t.Errorf("Scanning for letters didn't retrieve string")
	}
	l.rewind()
	if l.next() != "Ŧ" {
		t.Errorf("Rewinding scan didn't put lexer in correct place")
	}
}

func TestNumber(t *testing.T) {
	tests := []struct {
		i string  // input
		o float64 // output
		p int     // expected value of lexer.pos afterwards
	}{
		// GOOD CASES
		{"0", 0, 1},
		{"-1", -1, 2},
		{"1.25", 1.25, 4},
		{"-12345.6789", -12345.6789, 11},
		{"1e6", 1e6, 3},
		{"-1.23e45", -1.23e45, 8},
		{"1.23e-45", 1.23e-45, 8},
		// BAD CASES
		{"1e999", 0, 5},   // > MaxFloat
		{"1e-999", 0, 6},  // < MinFloat
		{"NaN", 0, 0},     // should result in ParseFloat("")
		{"a123.45", 0, 0}, //   ""
		// UGLY CASES
		{"0xf00", 0, 1}, // Hex not supported yet
		{"0b010", 0, 1}, // Binary not supported yet
		{"1foo", 1, 1},  // Stops at first non-digit
		{"०१२", 0, 9},   // 012 in devanagari digits
		{"໘໔໓", 0, 9},   // 843 in lao digits
		// I guess those poor devanagari etc. are SOL ;-(
	}

	for i, tc := range tests {
		l := &lexer{input: tc.i}
		if o := l.number(); o != tc.o {
			t.Errorf("number(%d): '%s' result %f != %f", i, tc.i, o, tc.o)
		}
		if l.pos != tc.p {
			t.Errorf("number(%d): '%s' pos %d != %d", i, tc.i, l.pos, tc.p)
		}
	}
}

type tt struct {
	i string    // input
	k tokenKind // token.kind
	s string    // token.strval
	n float64   // token.numval
}

func TestToken(t *testing.T) {
	tests := []tt{
		{"", T_EOF, "", 0},
		{"       ", T_EOF, "", 0},
		{"+", T_OP, "+", 0},
		// We can't test '-' as a standalone operator because the lexer
		// assumes that it's the unary minus at the beginning of a number.
		//		{"-",       T_OP,    "-",     0},
		{"*", T_OP, "*", 0},
		{"/", T_OP, "/", 0},
		{"**", T_OP, "**", 0},
		{"^", T_OP, "^", 0},
		{"%", T_OP, "%", 0},
		{"(", T_LPAR, "(", 0},
		{")", T_RPAR, ")", 0},
		{",", T_COMMA, ",", 0},
		{"1234.5", T_NUM, "", 1234.5},
		{"-1234.5", T_NUM, "", -1234.5},
		{"pie", T_NFI, "pie", 0},
		{"-pie", T_NFI, "-pie", 0},
		{"&", T_NFI, "&", 0},
	}
	// Test all the fucntions are correctly recognised
	for fun, _ := range functionMap {
		tests = append(tests, tt{fun, T_FUNC, fun, 0})
	}
	// Test all the constants are correctly recognised
	for con, val := range ConstMap {
		tests = append(tests, tt{con, T_NUM, "", val})
		tests = append(tests, tt{"-" + con, T_NUM, "", -val})
	}

	for i, tc := range tests {
		l := &lexer{input: tc.i}
		tok := l.token()
		if tok.kind != tc.k {
			t.Errorf("token(%d) '%s' kind mismatch, %d != %d",
				i, tc.i, tok.kind, tc.k)
		}
		if tok.strval != tc.s {
			t.Errorf("token(%d) '%s' str mismatch, %s != %s",
				i, tc.i, tok.strval, tc.s)
		}
		if tok.numval != tc.n {
			t.Errorf("token(%d) '%s' num mismatch, %f != %f",
				i, tc.i, tok.numval, tc.n)
		}
	}
}

func TestTokenMinus(t *testing.T) {
	// Special minus testing
	tests := []struct {
		i string // input
		n int    // expect only the nth token to be a minus
	}{
		{"-e - -3", 2},
		{"atan2(3, -2) - 1", 7},
		{"cos(-pi-2)", 4},
		{"5*6-7", 4},
	}
	for i, tc := range tests {
		l := &lexer{input: tc.i}
		tok := &token{T_NUM, "", 0} // start things off ...
		for j := 1; tok.kind != T_EOF; j++ {
			tok = l.token()
			if j != tc.n && tok.kind == T_OP && tok.strval == "-" {
				t.Errorf("token(%d) '%s' unexpected - operator at token %d",
					i, tc.i, j)
			}
			if j == tc.n && (tok.kind != T_OP ||
				(tok.kind == T_OP && tok.strval != "-")) {
				t.Errorf("token(%d) '%s' unexpected %#v token",
					i, tc.i, tok)
			}
		}
	}
}

// for testing.
func (ts *tokenStack) serialise() string {
	if ts == nil {
		return ""
	}
	s := make([]string, len(*ts))
	for i, t := range *ts {
		if t.kind == T_NUM {
			s[i] = fmt.Sprintf("%g", t.numval)
		} else {
			s[i] = t.strval
		}
	}
	return strings.Join(s, "")
}

func TestTokens(t *testing.T) {
	tests := []string{
		"2+4",
		"(2*2+6^3)",
		"cos(4/3)*atan2(1,2)",
		"(1+((2+(3+(4*(5*6-7)+8)*9))*10))",
		"&D(foo)", // lots of T_NFI
	}
	for i, tc := range tests {
		ts := (&lexer{input: tc}).tokens()
		if s := ts.serialise(); s != tc {
			t.Errorf("Unexpected string output for %d, expected: %s, got: %s",
				i, tc, s)
			t.Errorf("%s", ts)
		}
	}
}

func TestShunt(t *testing.T) {
	// This tests that a set of inputs produces the expected outputs
	tests := []struct {
		i string
		o string
		e bool
	}{
		{"2+4", "24+", false},
		{"(2+4)*6", "24+6*", false},
		{"2+4*6", "246*+", false},
		{"2+3+4+5+6+7+8*9", "23+4+5+6+7+89*+", false},
		{"2+3*4^5", "2345^*+", false},
		{"tan(answer)", "42tan", false},
		{"(1+((2+(3+(4*(5*6-7)+8)*9))*10))", "123456*7-*8+9*++10*+", false},
		{"1*atan2((2+3)*4,5*(6+7))+8", "123+4*567+*atan2*8+", false},
		{"(2+4", "", true},
		{"2+4)", "", true},
		{"cos(,)", "", true},
		{"cos()", "cos", false}, // not a parse error at shunt time
	}
	for i, tc := range tests {
		ts, err := shunt((&lexer{input: tc.i}).tokens())
		if (err == nil) == tc.e {
			t.Errorf("Bad shunt error state for %d (err=%v).", i, err)
		}
		if s := ts.serialise(); s != tc.o {
			t.Errorf("Bad shunted output for %d, expected: %s, got: %s",
				i, tc.o, s)
		}
	}
}

func TestShuntStep(t *testing.T) {
	tests := []struct {
		tok    *token
		si, so string
		ai, ao string
		oi, oo string
		e      bool
	}{
		// An unrecognised token should result in an error and no mutations.
		{&token{T_NFI, "!", 0},
			"+cos", "+cos",
			"cos", "cos",
			"12", "12",
			true},

		// If the token is a number, then add it to the output queue. 
		{&token{T_NUM, "", 5},
			"+cos", "+cos",
			"cos", "cos",
			"12", "125",
			false},

		// If the token is an operator, op1, then:
		// - While there is an operator token, op2, at the top of the stack,
		//   and its precedence is less than or equal to that of op2:
		//   NOTE: short-circuit here because no right-associative operators!
		//   - Pop op2 off the stack, onto the output queue;
		// - Finally, push op1 onto the stack.
		// First -- no operator at top of stack. Just push to stack.
		{&token{T_OP, "^", 0},
			"+cos", "+cos^",
			"cos", "cos",
			"12", "12",
			false},
		// Second -- two operators at top of stack, first equal, second higher.
		// - Pop both operators to output (tok precedence <= stack precedence),
		//   then push token operator to stack.
		{&token{T_OP, "+", 0},
			"+cos+*", "+cos+",
			"cos", "cos",
			"12", "12*+",
			false},
		// Third -- two operators at top of stack, first lower, second equal.
		// - Pop second operator to output, then push token operator to stack.
		{&token{T_OP, "*", 0},
			"+cos+*", "+cos+*",
			"cos", "cos",
			"12", "12*",
			false},
		// Fourth -- two operators at top of stack, first higher, second lower.
		// - Push token operator to stack (tok precedence > stack precedence).
		{&token{T_OP, "*", 0},
			"+cos^+", "+cos^+*",
			"cos", "cos",
			"12", "12",
			false},
		// Fifth -- nothing on stack. Just push to stack.
		{&token{T_OP, "*", 0},
			"", "*",
			"", "",
			"12", "12",
			false},

		// If the token is a function token, then push it onto the stack.
		// Also, push it onto the argcs stack for argument count checking.
		{&token{T_FUNC, "sin", 0},
			"+cos", "+cossin",
			"cos", "cossin",
			"12", "12",
			false},

		// If the token is a left parenthesis, then push it onto the stack.
		{&token{T_LPAR, "(", 0},
			"+cos", "+cos(",
			"cos", "cos",
			"12", "12",
			false},

		// If the token is a right parenthesis:
		// - Until the token at the top of the stack is a left parenthesis,
		//     pop operators off the stack onto the output queue.
		// - Pop the left parenthesis from the stack, but not to output queue.
		// - If the token at the top of the stack is a function token,
		//   pop it onto the output queue.
		// - If the stack runs out without finding a left parenthesis,
		//   then there are mismatched parentheses.
		//
		// NOTE: with the testing setup here it is not possible to test the
		//       argument counting feature. It is tested separately later.
		// First -- two ops before ( on stack, pop both and drop (.
		{&token{T_RPAR, ")", 0},
			"+(*/", "+",
			"", "",
			"12", "12/*",
			false},
		// Second -- op before (, func after, pop both and drop (.
		{&token{T_RPAR, ")", 0},
			"+cos(/", "+",
			"", "", // NOTE: leave argcs empty to avoid breakage
			"12", "12/cos",
			false},
		// Third -- missing (, pop until stack underflow, return error.
		{&token{T_RPAR, ")", 0},
			"+cos/*", "",
			"", "",
			"12", "12*/cos+",
			true},

		// If the token is a function argument separator (e.g., a comma):
		// - Until the token at the top of the stack is a left parenthesis,
		//   pop operators off the stack onto the output queue. If no left 
		//   parentheses are encountered, either the separator was misplaced
		//   or parentheses were mismatched.
		//
		// NOTE: again here it's not possible to test argument counting.
		// First -- two ops before ( on stack, pop both and keep (.
		{&token{T_COMMA, ",", 0},
			"+(*/", "+(",
			"", "",
			"12", "12/*",
			false},
		// Second, missing (, pop until stack underflow, return error.
		{&token{T_COMMA, ",", 0},
			"+cos/*", "",
			"", "",
			"12", "12*/cos+",
			true},
	}

	for i, tc := range tests {
		// Initialise state from inputs
		s := (&lexer{input: tc.si}).tokens()
		a := (&lexer{input: tc.ai}).tokens()
		o := (&lexer{input: tc.oi}).tokens()

		err := shuntStep(tc.tok, s, a, o)
		if (err == nil) == tc.e {
			t.Errorf("Bad step error state for %d (err=%v)", i, err)
		}
		if so := s.serialise(); so != tc.so {
			t.Errorf("Bad step stack output for %d, expected: %s, got: %s",
				i, tc.so, so)
		}
		if ao := a.serialise(); ao != tc.ao {
			t.Errorf("Bad step argcs output for %d, expected: %s, got: %s",
				i, tc.ao, ao)
		}
		if oo := o.serialise(); oo != tc.oo {
			t.Errorf("Bad step output output for %d, expected: %s, got: %s",
				i, tc.oo, oo)
		}
	}
}

func TestShuntStepArgCounting(t *testing.T) {
	// Argument counting (ab)uses the numval field in T_FUNC tokens
	// to count the number of arguments seen for a specific function.
	// Since this information is obscured by ts.serialise() we have to
	// jump through some hoops to check it's working properly.

	// Initial state, parsed 1+cos(2
	s := (&lexer{input: "+cos("}).tokens()
	a := (&lexer{input: "cos"}).tokens()
	ftok := (*a)[0]
	o := (&lexer{input: "12"}).tokens()

	// Correct syntax: T_RPAR
	// Results in 2x stack pop and argcs pop
	err := shuntStep(&token{T_RPAR, ")", 0}, s, a, o)
	// We know that T_RPAR works correctly for the stack and output already.
	if err != nil {
		t.Errorf("Bad step error state for arg count 1, (err=%v)", err)
	}
	if ao := a.serialise(); ao != "" {
		t.Errorf("Bad step argcs output for arg count 1, expected: '', got: %s",
			ao)
	}
	if ftok.numval != 1 {
		t.Errorf("Arg count 1 not correctly incremented to 1.")
	}

	// RESET!
	s = (&lexer{input: "+cos("}).tokens()
	a = (&lexer{input: "cos"}).tokens()
	ftok = (*a)[0]
	o = (&lexer{input: "12"}).tokens()

	// Bad syntax: T_COMMA, T_RPAR
	err = shuntStep(&token{T_COMMA, ",", 0}, s, a, o)
	// The comma itself shouldn't have resulted in any error ...
	if err != nil {
		t.Errorf("Bad step error state for arg count 2, (err=%v)", err)
	}
	//  ... but it should have incremented the arg count.
	if ftok.numval != 1 {
		t.Errorf("Arg count 2 not correctly incremented to 1.")
	}
	err = shuntStep(&token{T_RPAR, ")", 0}, s, a, o)
	if err == nil {
		t.Errorf("Bad step error state for arg count 3, (err=%v)", err)
	}
	if ftok.numval != 2 {
		t.Errorf("Arg count 3 not correctly incremented to 2.")
	}

	// RESET a different way
	s = (&lexer{input: "+atan2("}).tokens()
	a = (&lexer{input: "atan2"}).tokens()
	ftok = (*a)[0]
	o = (&lexer{input: "12"}).tokens()

	// Bad syntax: T_RPAR (expecting 2 arguments)
	err = shuntStep(&token{T_RPAR, ")", 0}, s, a, o)
	if err == nil {
		t.Errorf("Bad step error state for arg count 4, (err=%v)", err)
	}
	if ftok.numval != 1 {
		t.Errorf("Arg count 4 not correctly incremented to 1.")
	}
}
