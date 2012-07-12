package util

import (
	"fmt"
	"math"
	"strconv"
	"unicode"
	"unicode/utf8"
)

type tokenKind int
const (
	T_EOF tokenKind = iota
	T_NFI
	T_NUM
	T_OP
	T_FUNC
	T_LPAR
	T_RPAR
	T_COMMA
)

type token struct {
	kind tokenKind
	strval string
	numval float64
}

type tokenStack []token

func ts(l int) *tokenStack {
	s := tokenStack(make([]token, 0, l))
	return &s
}

func (ts *tokenStack) l() int {
	return len(*ts)
}

func (ts *tokenStack) push(t token) {
	*ts = append(*ts, t)
}

func (ts *tokenStack) pop() (t token, e error) {
	l := ts.l()
	if l == 0 {
		return token{}, fmt.Errorf("stack underflow")
	}
	*ts, t = (*ts)[:l-1], (*ts)[l-1]
	return t, nil
}

func (ts *tokenStack) shift() (t token, e error) {
	if ts.l() == 0 {
		return token{}, fmt.Errorf("stack underflow")
	}
	t, *ts = (*ts)[0], (*ts)[1:]
	return t, nil
}

// OK, this is kind of horrific and probably comes from bad design decisions
// I APOLOGISE FOR NOTHING
func (ts *tokenStack) getNums(n int) ([]token, error) {
	nums := make([]token, 0, n)
	for ; n > 0; n-- {
		tok, err := ts.pop()
		if err != nil || tok.kind != T_NUM {
			return nil, fmt.Errorf("syntax error")
		}
		nums = append(nums, tok)
	}
	return nums, nil
}

var constMap = map[string]float64{
	"e": math.E,
	"pi": math.Pi,
	"phi": math.Phi,
	"answer": 42,
}

var precedenceMap = map[string]int{
	"**": 3, "^": 3,
	"*": 2, "/": 2, "%": 2,
	"+": 1, "-": 1,
}
	
func precedence(op1, op2 string) bool {
	// Return true (i.e. pop the stack) if op1 is left-associative and its
	// precedence is less than or equal to that of o2. All our operators 
	// are left-associative, so this is easy :-)
	return precedenceMap[op1] <= precedenceMap[op2]
}

type funcop func(*tokenStack) error
type function struct {
	argc int
	exec funcop
}

var operatorMap = map[string]funcop{
	"**": f_power,
	"^": f_power,
	"*": f_mult,
	"/": f_div,
	"%": f_mod,
	"+": f_plus,
	"-": f_minus,
}

func f_power(ts *tokenStack) error {
	n, err := ts.getNums(2)
	if err != nil { return err }
	ts.push(token{T_NUM, "", math.Pow(n[0].numval, n[1].numval)})
	return nil
}

func f_mult(ts *tokenStack) error {
	n, err := ts.getNums(2)
	if err != nil { return err }
	ts.push(token{T_NUM, "", n[1].numval * n[0].numval})
	return nil
}

func f_div(ts *tokenStack) error {
	n, err := ts.getNums(2)
	if err != nil { return err }
	ts.push(token{T_NUM, "", n[1].numval / n[0].numval})
	return nil
}

func f_mod(ts *tokenStack) error {
	n, err := ts.getNums(2)
	if err != nil { return err }
	ts.push(token{T_NUM, "", math.Mod(n[1].numval, n[0].numval)})
	return nil
}

func f_plus(ts *tokenStack) error {
	n, err := ts.getNums(2)
	if err != nil { return err }
	ts.push(token{T_NUM, "", n[1].numval + n[0].numval})
	return nil
}

func f_minus(ts *tokenStack) error {
	n, err := ts.getNums(2)
	if err != nil { return err }
	ts.push(token{T_NUM, "", n[1].numval - n[0].numval})
	return nil
}

// fortunately, almost all the interesting math functions are f(x) -> x
func functionise(f func(float64)float64) function {
	return function{1, func(ts *tokenStack) error {
		n, err := ts.getNums(1)
		if err != nil { return err }
		ts.push(token{T_NUM, "", f(n[0].numval)})
		return nil
	}}
}

var functionMap = map[string]function{
	"cos": functionise(math.Cos),
}

type lexer struct {
	input string
	start, pos, width int
	unaryMinus bool
}

func (l *lexer) peek() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return 0
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	return r
}

func (l *lexer) next() string {
	l.start = l.pos
	r := l.peek()
	l.pos += l.width
	return string(r)
}

func (l *lexer) rewind() {
	l.pos = l.start
}

func (l *lexer) scan(f func(rune) bool) string {
	l.start = l.pos
	for f(l.peek()) {
		l.pos += l.width
	}
	return l.input[l.start:l.pos]
}

func (l *lexer) number() float64 {
	// we say a number is: [0-9]+(.[0-9]+)?([eE][0-9]+)?
	s := l.pos // l.start is reset through the multiple scans
	if l.peek() == '-' { l.next() }
	l.scan(unicode.IsDigit)
	if l.next() == "." {
		l.scan(unicode.IsDigit)
	} else {
		l.rewind()
	}
	c := l.next()
	if c == "e" || c == "E" {
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

func (l *lexer) token() (tok token) {
	l.scan(unicode.IsSpace)
	r := l.peek()

	switch {
	case r == 0:
		tok = token{T_EOF, "", 0}
	case r == '-':
		// could be a prefix - as in -12. This is only valid
		// if the last token was an operator (6 - 4) vs (6 - -4)
		if l.unaryMinus {
			tok = token{T_NUM, "", l.number()}
		} else {
			tok = token{T_OP, l.next(), 0}
		}
	case r == '*':
		// ** is often the power operator
		l.next()
		if l.peek() == '*' {
			tok = token{T_OP, "**", 0}
		} else {
			tok = token{T_OP, "*", 0}
		}
	case r == '+' || r == '/' || r == '%' || r == '^':
		tok = token{T_OP, l.next(), 0}
	case r == '(':
		tok = token{T_LPAR, l.next(), 0}
	case r == ')':
		tok = token{T_RPAR, l.next(), 0}
	case r == ',':
		tok = token{T_COMMA, l.next(), 0}
	case unicode.IsLetter(r):
		// could be a constant or a defined function
		str := l.scan(unicode.IsLetter)
		if _, ok := functionMap[str]; ok {
		// since we know our defined functions, let's just check
			tok = token{T_FUNC, str, 0}
		} else if num, ok := constMap[str]; ok {
		// since we know our defined constants, ...
			tok = token{T_NUM, "", num}
		} else {
			// keeping this simple here, error handling for unknown strings
			// can be done further up in the parser
			tok = token{T_NFI, str, 0}
		}
	case unicode.IsNumber(r):
		tok = token{T_NUM, "", l.number()}
	default:
		tok = token{T_NFI, l.next(), 0}
	}
	switch tok.kind {
	case T_OP, T_LPAR, T_COMMA:
		l.unaryMinus = true
	default:
		l.unaryMinus = false
	}
	return
}

type parser struct {
	l *lexer
	ops *tokenStack
}

func (p *parser) shunt() error {
	stack := ts(len(p.l.input))
	// This is abusing the "numval" field of T_FUNC tokens to check
	// they have been given the correct number of arguments. This
	// should hopefully avoid a number of incorrect calc results.
	argcs := ts(5)
SHUNT:
	for {
		tok := p.l.token()
		switch tok.kind {
		case T_EOF:
			break SHUNT
		case T_NFI:
			return fmt.Errorf("Unrecognised '%#v' in expression", tok)
		case T_NUM:
			p.ops.push(tok)
		case T_OP:
			for {
				top, err := stack.pop()
				if err != nil { break }
				if top.kind == T_OP && precedence(tok.strval, top.strval) {
					p.ops.push(top)
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
				if top.kind == T_LPAR {
					break
				}
				p.ops.push(top)
			}
			if top, err := stack.pop(); err == nil {
				if top.kind == T_FUNC {
					p.ops.push(top)
					if f, err := argcs.pop(); err == nil {
						f.numval++
						fmt.Printf("func %s: expected %d, got %f\n",
							f.strval, functionMap[f.strval].argc, f.numval)
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
			if argcs.l() > 0 {
				(*argcs)[argcs.l()-1].numval++
			}
			for {
				top, err := stack.pop()
				if err != nil { return err }
				if top.kind == T_LPAR {
					stack.push(top)
					break
				}
				p.ops.push(top)
			}
		}
	}
	for stack.l() > 0 {
		tok, _ := stack.pop()
		if tok.kind == T_LPAR {
			return fmt.Errorf("stack overflow")
		}
		p.ops.push(tok)
	}
	return nil
}

func (p *parser) calc() (float64, error) {
	stack := ts(len(p.l.input))
	for p.ops.l() > 0 {
		tok, _ := p.ops.shift()
		switch tok.kind {
		// only NUM, OP, FUNC should have made it this far
		case T_NUM:
			stack.push(tok)
		case T_OP:
			op := operatorMap[tok.strval]
			err := op(stack)
			if err != nil { return 0, err }
		case T_FUNC:
			f := functionMap[tok.strval]
			err := f.exec(stack)
			if err != nil { return 0, err }
		default:
			return 0, fmt.Errorf("token %#v in calc", tok)
		}
	}
	if stack.l() != 1 {
		return 0, fmt.Errorf("Syntax error, overfull stack")
	}
	ret, _ := stack.pop()
	return ret.numval, nil
}

func Calc(input string) (float64, error) {
	p := &parser{
		l: &lexer{input: input},
		ops: ts(len(input)),
	}
	// Dijkstra's shunting-yard algorithm to RPN input
	if err := p.shunt(); err != nil {
		return 0, err
	}
	fmt.Printf("Shunted. Operation list:\n")
	for i, v := range *p.ops {
		fmt.Printf("%2d: %#v\n", i, v)
	}
	// Perform RPN calculation
	return p.calc()
}

