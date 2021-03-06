package calc

import (
	"fmt"
	"math"
	"strings"
	"unicode"

	"github.com/fluffle/sp0rkle/util"
)

//------------------------------------------------------------------------------
// Constants, Operators, and Functions for Expressions.

// A TokenMap allows the calc parser to recognise arbitrary
// strings in the input and do maths on them.
type TokenMap map[string]float64

// constMap defines constants that can be used in expressions
// passed to Calc(). These are statically compiled in.
var constMap = TokenMap{
	"e":      math.E,
	"pi":     math.Pi,
	"phi":    math.Phi,
	"answer": 42,
}

// precedenceMap defines the precedence of operators Calc() recognises.
var precedenceMap = map[string]struct {
	prec int
	lAss bool
}{
	"**": {3, false},
	"^":  {3, false},
	"*":  {2, true},
	"/":  {2, true},
	"%":  {2, true},
	"+":  {1, true},
	"-":  {1, true},
}

// Return true -- and thus pop the stack in shunt() -- if op1 is left-
// associative and its precedence is less than or equal to that of op2,
// or if op1 is right-associative and its precedence is less than op2.
func precedence(op1, op2 string) bool {
	o1 := precedenceMap[op1]
	o2 := precedenceMap[op2]
	return o1.prec < o2.prec || (o1.lAss && o1.prec == o2.prec)
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
func functionise1(f func(float64) float64) function {
	return function{1, func(ts *tokenStack) error {
		n, err := ts.getNums(1)
		if err != nil {
			return err
		}
		ts.push(&token{T_NUM, "", f(n[0].numval)})
		return nil
	}}
}

// But a lot of the operators are f(x,y) -> z, and so we need two curries.
//     ...
// MMMmmmmmm. Curry.
func functionise2(f func(float64, float64) float64) function {
	return function{2, func(ts *tokenStack) error {
		n, err := ts.getNums(2)
		if err != nil {
			return err
		}
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
	"*":  functionise2(func(x, y float64) float64 { return x * y }),
	"/":  functionise2(func(x, y float64) float64 { return x / y }),
	"%":  functionise2(math.Mod),
	"+":  functionise2(func(x, y float64) float64 { return x + y }),
	"-":  functionise2(func(x, y float64) float64 { return x - y }),
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
	"max":   functionise2(math.Max), // 2 args
	"min":   functionise2(math.Min), // 2 args
	"sin":   functionise1(math.Sin),
	"sinh":  functionise1(math.Sinh),
	"sqrt":  functionise1(math.Sqrt),
	"tan":   functionise1(math.Tan),
	"tanh":  functionise1(math.Tanh),
}

//------------------------------------------------------------------------------
// The Lexer

// The lexer produces tokens of these kinds:
type tokenKind int

const (
	T_EOF   tokenKind = iota // end-of-file
	T_NFI                    // no fucking idea what this input is ;-)
	T_NUM                    // a number (fills numval)
	T_OP                     // an operator
	T_FUNC                   // a function
	T_LPAR                   // a left parenthesis (
	T_RPAR                   // a right parenthesis )
	T_COMMA                  // a comma ,
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
	kind   tokenKind
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
	// we extend the basic lexer util to produce tokens here.
	*util.Lexer
	binaryMinus bool
	toks        TokenMap
}

func calcLexer(i string) *lexer {
	return &lexer{Lexer: &util.Lexer{Input: i}, toks: constMap}
}

// token() produces tokens from the input string for use by the parser.
func (l *lexer) token() (tok *token) {
	l.Scan(unicode.IsSpace)
	r := l.Peek()

	switch {
	case r == 0:
		tok = &token{T_EOF, "", 0}
	case r == '+' || r == '/' || r == '%' || r == '^':
		tok = &token{T_OP, l.Next(), 0}
	case r == '(':
		tok = &token{T_LPAR, l.Next(), 0}
	case r == ')':
		tok = &token{T_RPAR, l.Next(), 0}
	case r == ',':
		tok = &token{T_COMMA, l.Next(), 0}
	case r == '-':
		// could be a prefix - as in -12. This is only valid
		// if the last token was an operator (6 - 4) vs (6 - -4)
		l.Next()
		if l.binaryMinus {
			tok = &token{T_OP, "-", 0}
		} else if unicode.IsLetter(l.Peek()) {
			// With many apologies, this seemed to be the best place
			// to hack in support for negative constants like "-pi"...
			str := l.Scan(unicode.IsLetter)
			if num, ok := l.toks[str]; ok {
				tok = &token{T_NUM, "", -num}
			} else {
				tok = &token{T_NFI, "-" + str, 0}
			}
		} else {
			l.Rewind()
			tok = &token{T_NUM, "", l.Number()}
		}
	case r == '*':
		// ** is often the power operator
		l.Next()
		if l.Peek() == '*' {
			l.Next()
			tok = &token{T_OP, "**", 0}
		} else {
			tok = &token{T_OP, "*", 0}
		}
	case unicode.IsLetter(r):
		// could be a constant or a defined function
		str := l.Scan(unicode.IsLetter)
		if _, ok := functionMap[str]; ok {
			// since we know our defined functions, let's just check
			// we need special case checking for atan2, log2, and
			// log10, because IsLetter doesn't match the 2 / 10...
			if c := l.Next(); c == "2" {
				str += "2"
			} else if c == "1" && l.Peek() == '0' {
				l.Next()
				str += "10"
			} else {
				l.Rewind()
			}
			tok = &token{T_FUNC, str, 0}
		} else if num, ok := l.toks[str]; ok {
			// since we know our defined constants, ...
			tok = &token{T_NUM, "", num}
		} else {
			// keeping this simple here, error handling for
			// unknown strings can be done up in the parser
			tok = &token{T_NFI, str, 0}
		}
	case unicode.IsDigit(r):
		tok = &token{T_NUM, "", l.Number()}
	default:
		tok = &token{T_NFI, l.Next(), 0}
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
	ts := ts(len(l.Lexer.Input))
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
	for i, tok := range *ts {
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
			if err != nil {
				break
			}
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
			if err != nil {
				return err
			}
			if top.kind == T_LPAR {
				break
			}
			output.push(top)
		}
		if top, err := stack.pop(); err == nil {
			if top.kind == T_FUNC {
				output.push(top)
				if f, err := argcs.pop(); err == nil {
					f.numval++
					if int(f.numval) != functionMap[f.strval].argc {
						return fmt.Errorf("Incorrect number of arguments"+
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
			if err != nil {
				return err
			}
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
			if err != nil {
				return 0, err
			}
		case T_FUNC:
			f := functionMap[tok.strval]
			err := f.exec(stack)
			if err != nil {
				return 0, err
			}
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
func Calc(input string, tokens TokenMap) (float64, error) {
	l := calcLexer(input)
	if tokens != nil {
		for k, v := range constMap {
			tokens[k] = v
		}
		l.toks = tokens
	}
	// Dijkstra's shunting-yard algorithm to RPN input
	ops, err := shunt(l.tokens())
	if err != nil {
		return 0, err
	}
	// Perform RPN calculation
	return calc(ops)
}
