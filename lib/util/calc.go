package util

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

//------------------------------------------------------------------------------
// Constants, Operators, and Functions for Expressions.

// ConstMap defines constants that can be used in expressions
// passed to Calc(). It's exposed so others can be added dynamically.
var ConstMap = map[string]float64{
	"e": math.E,
	"pi": math.Pi,
	"phi": math.Phi,
	"answer": 42,
}

// precedenceMap defines the precedence of operators Calc() recognises.
// NOTE: These are all binary and left-associative, which simplified
//       some of the code below. If any operators that are not left-
//       associative are added, change precedence() below accordingly.
var precedenceMap = map[string]int{
	"**": 3, "^": 3,
	"*":  2, "/": 2, "%": 2,
	"+":  1, "-": 1,
}

// Return true -- and thus pop the stack in shunt() -- if op1 is left-
// associative and its precedence is less than or equal to that of op2.
func precedence(op1, op2 string) bool {
	// NOTE: If you add a right-associative operator to operatorMap,
	//       this conditional will need to be modified.
	return precedenceMap[op1] <= precedenceMap[op2]
}

// The function type (poor naming, meh) performs the work for any
// function or operator the calculator can perform.
type function struct {
	argc int
	exec func(*tokenStack) error
}

// The functionise functions curry various math functions into function
// structs for the calculator to use. Function function function function
// function function function function mushrooooom mushrooooom.
// Almost all the interesting math functions are f(x) -> z
func functionise1(f func(float64)float64) function {
	return function{1, func(ts *tokenStack) error {
		n, err := ts.getNums(1)
		if err != nil { return err }
		ts.push(&token{T_NUM, "", f(n[0].numval)})
		return nil
	}}
}

// But a lot of the operators are f(x,y) -> z, and so we need two curries.
//     ...
// MMMmmmmmm. Curry.
func functionise2(f func(float64,float64)float64) function {
	return function{2, func(ts *tokenStack) error {
		n, err := ts.getNums(2)
		if err != nil { return err }
		ts.push(&token{T_NUM, "", f(n[0].numval, n[1].numval)})
		return nil
	}}
}

// operatorMap defines the set of operators that can be used in expressions
// passed to Calc(). If you add something to this remember to update the
// precedence map too, or things will get panic()y.
var operatorMap = map[string]function{
	"**": functionise2(math.Pow),
	"^":  functionise2(math.Pow),
	"*":  functionise2(func(x,y float64)float64{return x*y}),
	"/":  functionise2(func(x,y float64)float64{return x/y}),
	"%":  functionise2(math.Mod),
	"+":  functionise2(func(x,y float64)float64{return x+y}),
	"-":  functionise2(func(x,y float64)float64{return x-y}),
}

// functionMap defines the set of functions that can be used in expressions
// passed to Calc. It's just passthroughs to math.* currently.
var functionMap = map[string]function{
	"abs":   functionise1(math.Abs),
	"acos":  functionise1(math.Acos),
	"acosh": functionise1(math.Acosh),
	"asin":  functionise1(math.Asin),
	"asinh": functionise1(math.Asinh),
	"atan":  functionise1(math.Atan),
	"atan2": functionise2(math.Atan2), // 2 args
	"atanh": functionise1(math.Atanh),
	"cbrt":  functionise1(math.Cbrt),
	"ceil":  functionise1(math.Ceil),
	"cos":   functionise1(math.Cos),
	"cosh":  functionise1(math.Cosh),
	"exp":   functionise1(math.Exp),
	"exp2":  functionise1(math.Exp2),
	"floor": functionise1(math.Floor),
	"gamma": functionise1(math.Gamma),
	"hypot": functionise2(math.Hypot), // 2 args
	"int":   functionise1(math.Trunc), // renamed
	"log":   functionise1(math.Log),
	"log10": functionise1(math.Log10),
	"log2":  functionise1(math.Log2),
	"logb":  functionise1(math.Logb),
	"max":   functionise2(math.Max),   // 2 args
	"min":   functionise2(math.Min),   // 2 args
	"sin":   functionise1(math.Sin),
	"sinh":  functionise1(math.Sinh),
	"sqrt":  functionise1(math.Sqrt),
	"tan":   functionise1(math.Tan),
	"tanh":  functionise1(math.Tanh),
}

//------------------------------------------------------------------------------
// The Lexer
//
// (admittedly, this could be combined with the parser in a single struct...)

// The lexer produces tokens of these kinds:
type tokenKind int
const (
	T_EOF tokenKind = iota // end-of-file
	T_NFI                  // no fucking idea what this input is ;-)
	T_NUM                  // a number (fills numval)
	T_OP                   // an operator
	T_FUNC                 // a function
	T_LPAR                 // a left parenthesis (
	T_RPAR                 // a right parenthesis )
	T_COMMA                // a comma ,
)
var kindMap = map[tokenKind]string{
	T_EOF:   "EOF",
	T_NFI:   "NFI",
	T_NUM:   "NUM",
	T_OP:    "OP",
	T_FUNC:  "FUNC",
	T_LPAR:  "LPAR",
	T_RPAR:  "RPAR",
	T_COMMA: "COMMA",
}

type token struct {
	kind tokenKind
	strval string
	numval float64
}

func (t token) String() string {
	switch t.kind {
	case T_EOF, T_LPAR, T_RPAR, T_COMMA:
		return kindMap[t.kind]
	case T_NFI, T_OP, T_FUNC:
		return fmt.Sprintf("%s{%s}", kindMap[t.kind], t.strval)
	case T_NUM:
		return fmt.Sprintf("%s{%g}", kindMap[t.kind], t.numval)
	}
	return ""
}

type lexer struct {
	input string
	start, pos, width int
	binaryMinus bool
}

// peek() returns the utf8 rune that is at lexer.pos in the input string.
// It does not move input.pos; repeated peek()s will return the same rune.
func (l *lexer) peek() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	if r == utf8.RuneError {
		// Treat bad unicode as EOF.
		l.width = 0
		return 0
	}
	return r
}

// next() returns the utf8 rune (in string form, for convenience elsewhere)
// that is at lexer.pos in the input string, then advances lexer.pos past it.
func (l *lexer) next() string {
	l.start = l.pos
	r := l.peek()
	l.pos += l.width
	return string(r)
}

// scan() returns the sequence of runes in the input string anchored
// at lexer.pos that the supplied function returns true for. Usefully,
// unicode.IsDigit et al. fit the required function signature ;-)
func (l *lexer) scan(f func(rune) bool) string {
	l.start = l.pos
	for f(l.peek()) {
		l.pos += l.width
	}
	return l.input[l.start:l.pos]
}

// rewind() undoes the last next() or scan() by resetting lexer.pos.
func (l *lexer) rewind() {
	l.pos = l.start
}

// number() is a higher-level function that extracts a number from the
// input beginning at lexer.pos. A number matches the following regex:
//     -?[0-9]+(.[0-9]+)?([eE]-?[0-9]+)?
func (l *lexer) number() float64 {
	s := l.pos // l.start is reset through the multiple scans
	if l.peek() == '-' { l.pos += l.width }
	l.scan(unicode.IsDigit)
	if l.next() == "." {
		l.scan(unicode.IsDigit)
	} else {
		l.rewind()
	}
	if c := l.next(); c == "e" || c == "E" {
		if l.peek() == '-' { l.pos += l.width }
		l.scan(unicode.IsDigit)
	} else {
		l.rewind()
	}
	l.start = s
	n, err := strconv.ParseFloat(l.input[s:l.pos], 64)
	if err != nil {
		// This might be a bad idea in the long run.
		return 0
	}
	return n
}

// token() produces tokens from the input string for use by the parser.
func (l *lexer) token() (tok *token) {
	l.scan(unicode.IsSpace)
	r := l.peek()

	switch {
	case r == 0:
		tok = &token{T_EOF, "", 0}
	case r == '+' || r == '/' || r == '%' || r == '^':
		tok = &token{T_OP, l.next(), 0}
	case r == '(':
		tok = &token{T_LPAR, l.next(), 0}
	case r == ')':
		tok = &token{T_RPAR, l.next(), 0}
	case r == ',':
		tok = &token{T_COMMA, l.next(), 0}
	case r == '-':
		// could be a prefix - as in -12. This is only valid
		// if the last token was an operator (6 - 4) vs (6 - -4)
		l.next()
		if l.binaryMinus {
			tok = &token{T_OP, "-", 0}
		} else if unicode.IsLetter(l.peek()) {
			// With many apologies, this seemed to be the best place
			// to hack in support for negative constants like "-pi"...
			str := l.scan(unicode.IsLetter)
			if num, ok := ConstMap[str]; ok {
				tok = &token{T_NUM, "", -num}
			} else {
				tok = &token{T_NFI, "-"+str, 0}
			}
		} else {
			l.rewind()
			tok = &token{T_NUM, "", l.number()}
		}
	case r == '*':
		// ** is often the power operator
		l.next()
		if l.peek() == '*' {
			tok = &token{T_OP, "**", 0}
		} else {
			tok = &token{T_OP, "*", 0}
		}
	case unicode.IsLetter(r):
		// could be a constant or a defined function
		str := l.scan(unicode.IsLetter)
		if _, ok := functionMap[str]; ok {
		// since we know our defined functions, let's just check
			// we need special case checking for atan2, log2, and
			// log10, because IsLetter doesn't match the 2 / 10...
			if c := l.next(); c == "2" {
				str += "2"
			} else if c == "1" && l.peek() == '0' {
				l.next()
				str += "10"
			} else {
				l.rewind()
			}
			tok = &token{T_FUNC, str, 0}
		} else if num, ok := ConstMap[str]; ok {
		// since we know our defined constants, ...
			tok = &token{T_NUM, "", num}
		} else {
			// keeping this simple here, error handling for
			// unknown strings can be done up in the parser
			tok = &token{T_NFI, str, 0}
		}
	case unicode.IsDigit(r):
		tok = &token{T_NUM, "", l.number()}
	default:
		tok = &token{T_NFI, l.next(), 0}
	}
	switch tok.kind {
	case T_OP, T_LPAR, T_COMMA:
		l.binaryMinus = false
	default:
		l.binaryMinus = true
	}
	return
}

func (l *lexer) tokens() *tokenStack {
	ts := ts(len(l.input))
	for tok := l.token(); tok.kind != T_EOF; tok = l.token() {
		ts.push(tok)
	}
	return ts
}

//------------------------------------------------------------------------------
// The "Parser"

// To perform both the shunting-yard infix->rpn conversion and the subsequent
// calculation, we need a stack to push and pop from. A slice of token pointers.
type tokenStack []*token

// This makes us a new stack which can hold size elements before resizing.
func ts(size int) *tokenStack {
	ts := tokenStack(make([]*token, 0, size))
	return &ts
}

func (ts *tokenStack) String() string {
	s := make([]string, len(*ts))
	for i, tok := range(*ts) {
		s[i] = fmt.Sprintf("%2d: %s", i, tok)
	}
	return strings.Join(s, "\n")
}

// push() adds an item to the stack
func (ts *tokenStack) push(t *token) {
	*ts = append(*ts, t)
}

// pop() removes an item from the stack
func (ts *tokenStack) pop() (t *token, e error) {
	l := len(*ts)
	if l == 0 {
		return nil, fmt.Errorf("stack underflow")
	}
	*ts, t = (*ts)[:l-1], (*ts)[l-1]
	return t, nil
}

// getNums() pops n T_NUM tokens from the stack and returns them in a slice.
// It's a "helper" function for the functioniseX duo above.
//    ...
// OK, this is kind of horrific and probably comes from bad design decisions
// I APOLOGISE FOR NOTHING
func (ts *tokenStack) getNums(n int) ([]*token, error) {
	nums := make([]*token, n)
	for n--; n >= 0; n-- {
		tok, err := ts.pop()
		if err != nil || tok.kind != T_NUM {
			return nil, fmt.Errorf("syntax error")
		}
		// Work backwards here because preserving stack ordering makes
		// functionise2's argument ordering more intuitive.
		nums[n] = tok
	}
	return nums, nil
}

// shunt() implements a version of Dijkstra's Shunting-Yard algorithm:
//     http://en.wikipedia.org/wiki/Shunting-yard_algorithm
func shunt(input *tokenStack) (*tokenStack, error) {
	stack := ts(len(*input))
	output := ts(len(*input))
	// This is abusing the "numval" field of T_FUNC tokens to check
	// they have been given the correct number of arguments. This
	// should hopefully avoid a number of incorrect calc results.
	argcs := ts(5)
	for _, tok := range *input {
		if err := shuntStep(tok, stack, argcs, output); err != nil {
			return nil, err
		}
	}
	for len(*stack) > 0 {
		tok, _ := stack.pop()
		if tok.kind == T_LPAR {
			return nil, fmt.Errorf("stack overflow")
		}
		output.push(tok)
	}
	return output, nil
}

// shuntStep breaks out a single shunt to make it easier to test
func shuntStep(tok *token, stack, argcs, output *tokenStack) error {
	switch tok.kind {
	case T_NFI:
		return fmt.Errorf("Unrecognised '%#v' in expression", tok)
	case T_NUM:
		output.push(tok)
	case T_OP:
		for {
			top, err := stack.pop()
			if err != nil { break }
			if top.kind == T_OP && precedence(tok.strval, top.strval) {
				output.push(top)
			} else {
				stack.push(top)
				break
			}
		}
		stack.push(tok)
	case T_FUNC:
		stack.push(tok)
		argcs.push(tok)
	case T_LPAR:
		stack.push(tok)
	case T_RPAR:
		for {
			top, err := stack.pop()
			if err != nil { return err }
			if top.kind == T_LPAR { break }
			output.push(top)
		}
		if top, err := stack.pop(); err == nil {
			if top.kind == T_FUNC {
				output.push(top)
				if f, err := argcs.pop(); err == nil {
					f.numval++
					if int(f.numval) != functionMap[f.strval].argc {
						return fmt.Errorf("Incorrect number of arguments" +
							" for function '%s'.", f.strval)
					}
				}
			} else {
				stack.push(top)
			}
		}
	case T_COMMA:
		if l := len(*argcs); l > 0 {
			(*argcs)[l-1].numval++
		}
		for {
			top, err := stack.pop()
			if err != nil { return err }
			if top.kind == T_LPAR {
				stack.push(top)
				break
			}
			output.push(top)
		}
	}
	return nil
}

// calc() takes the rpn token list and applies the functions and
// operators to the numbers to get a result.
func calc(ops *tokenStack) (float64, error) {
	stack := ts(len(*ops))
	for _, tok := range *ops {
		switch tok.kind {
		// only NUM, OP, FUNC should have made it this far
		case T_NUM:
			stack.push(tok)
		case T_OP:
			op := operatorMap[tok.strval]
			err := op.exec(stack)
			if err != nil { return 0, err }
		case T_FUNC:
			f := functionMap[tok.strval]
			err := f.exec(stack)
			if err != nil { return 0, err }
		default:
			return 0, fmt.Errorf("token %#v in calc", tok)
		}
	}
	if len(*stack) != 1 {
		return 0, fmt.Errorf("syntax error, too many numbers")
	}
	ret, _ := stack.pop()
	return ret.numval, nil
}

//------------------------------------------------------------------------------
// The Interface

// Calc("some arbitrary maths string") -> (answer or zero, nil or error)
func Calc(input string) (float64, error) {
	l := &lexer{input: input}
	// Dijkstra's shunting-yard algorithm to RPN input
	ops, err := shunt(l.tokens());
	if  err != nil {
		return 0, err
	}
	// Perform RPN calculation
	return calc(ops)
}
