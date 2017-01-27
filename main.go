package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type tokenType int

const (
	_ tokenType = iota
	tokenLeftParen
	tokenRightParen
	tokenLeftBrace
	tokenRightBrace
	tokenComma
	tokenDot
	tokenMinus
	tokenPlus
	tokenSemicolon
	tokenSlash
	tokenStar

	tokenBang
	tokenBangEqual
	tokenEqual
	tokenEqualEqual
	tokenGreater
	tokenGreaterEqual
	tokenLess
	tokenLessEqual

	tokenIdentifier
	tokenString
	tokenNumber

	tokenAnd
	tokenClass
	tokenElse
	tokenFalse
	tokenFun
	tokenFor
	tokenIf
	tokenNil
	tokenOr
	tokenPrint
	tokenReturn
	tokenSuper
	tokenThis
	tokenTrue
	tokenVar
	tokenWhile

	tokenEOF
)

type token struct {
	Type    tokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func newToken(tt tokenType, lex string, lit interface{}, line int) *token {
	return &token{tt, lex, lit, line}
}

func (t token) String() string {
	return fmt.Sprintf("%v %v %v", t.Type, t.Lexeme, t.Literal)
}

var keywords = map[string]tokenType{
	"and":    tokenAnd,
	"class":  tokenClass,
	"else":   tokenElse,
	"false":  tokenFalse,
	"for":    tokenFor,
	"fun":    tokenFun,
	"if":     tokenIf,
	"nil":    tokenNil,
	"or":     tokenOr,
	"print":  tokenPrint,
	"return": tokenReturn,
	"super":  tokenSuper,
	"this":   tokenThis,
	"true":   tokenTrue,
	"var":    tokenVar,
	"while":  tokenWhile,
}

type scanner struct {
	Source string
	Reader *strings.Reader
	Tokens []*token

	start   int
	current int
	line    int
}

func newScanner(s string) *scanner {
	r := strings.NewReader(s)
	return &scanner{Source: s, Reader: r}
}

func (s *scanner) scanTokens() []*token {
	for !s.atEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.Tokens = append(s.Tokens, newToken(tokenEOF, "EOF", nil, s.line))
	return s.Tokens
}

func (s *scanner) atEnd() bool {
	return s.current >= len(s.Source)
}

func (s *scanner) scanToken() {
	ch := s.advance()
	switch ch {
	case '(':
		s.addToken(tokenLeftParen, nil)
	case ')':
		s.addToken(tokenRightParen, nil)
	case '{':
		s.addToken(tokenLeftBrace, nil)
	case '}':
		s.addToken(tokenRightBrace, nil)
	case ',':
		s.addToken(tokenComma, nil)
	case '.':
		s.addToken(tokenDot, nil)
	case '-':
		s.addToken(tokenMinus, nil)
	case '+':
		s.addToken(tokenPlus, nil)
	case ';':
		s.addToken(tokenSemicolon, nil)
	case '*':
		s.addToken(tokenStar, nil)
	case '!':
		s.addToken(ifToken(s.match('='), tokenBangEqual, tokenBang), nil)
	case '=':
		s.addToken(ifToken(s.match('='), tokenEqualEqual, tokenEqual), nil)
	case '<':
		s.addToken(ifToken(s.match('='), tokenLessEqual, tokenLess), nil)
	case '>':
		s.addToken(ifToken(s.match('='), tokenGreaterEqual, tokenLess), nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.atEnd() {
				s.advance()
			}
			break
		}
		s.addToken(tokenSlash, nil)
	case '\n':
		s.line++
	case ' ', '\r', '\t': // ignore white space
	case '"':
		s.string()
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		s.number()
	default:
		if isAlpha(ch) {
			s.identifier()
			break
		}

		problem(s.line, fmt.Errorf("Unexpected character: %q", ch))
	}
}

func ifToken(c bool, a tokenType, b tokenType) tokenType {
	if c {
		return a
	}
	return b
}

func isAlphaNumeric(ch rune) bool {
	return isAlpha(ch) || isDigit(ch)
}

func isAlpha(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (s *scanner) advance() rune {
	ch, sz, err := s.Reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("glox: advance error: %v", err))
	}
	s.current += sz
	return ch
}

func (s *scanner) match(ex rune) bool {
	if s.atEnd() {
		return false
	}

	ch, sz, err := s.Reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("glox: match error: %v", err))
	}
	if ch != ex {
		if err := s.Reader.UnreadRune(); err != nil {
			panic(fmt.Sprintf("glox: match unread error: %v", err))
		}
		return false
	}

	s.current += sz
	return true
}

func (s *scanner) peek() rune {
	if s.atEnd() {
		return 0
	}
	ch, _, err := s.Reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("glox: peek error: %v", err))
	}
	if err := s.Reader.UnreadRune(); err != nil {
		panic(fmt.Sprintf("glox: peek unread error: %v", err))
	}
	return ch
}

func (s *scanner) peekNext() rune {
	var (
		i   int
		err error
		ch  rune
	)

	// scan forward
	for ; i < 2 && !s.atEnd() && err == nil; i++ {
		ch, _, err = s.Reader.ReadRune()
	}

	// unwind
	if _, err = s.Reader.Seek(int64(s.current), 0); err != nil {
		panic(fmt.Sprintf("glox: peekNext read error: %v", err))
	}

	// failed to peek n
	if i < 2 {
		return 0
	}

	return ch
}

func (s *scanner) addToken(t tokenType, lit interface{}) {
	lex := s.Source[s.start:s.current]
	s.Tokens = append(s.Tokens, newToken(t, lex, lit, s.line))
}

func (s *scanner) string() {
	for s.peek() != '"' && !s.atEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.atEnd() {
		problem(s.line, errors.New("string error: unterminated string"))
		return
	}

	s.advance()

	v := s.Source[s.start+1 : s.current-1]
	s.addToken(tokenString, v)
}

func (s *scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	v := s.Source[s.start:s.current]
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		problem(s.line, fmt.Errorf("number error: %v", err))
		return
	}
	s.addToken(tokenNumber, f)
}

func (s *scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	v := s.Source[s.start:s.current]
	if i, ok := keywords[v]; ok {
		s.addToken(i, nil)
		return
	}

	s.addToken(tokenIdentifier, nil)
}

var hadError bool

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: glox [script]")
	} else if len(os.Args) == 2 {
		runFile(os.Args[1])
	} else {
		runPrompt()
	}
}

func runFile(f string) {
	buf, err := ioutil.ReadFile(f)
	if err != nil {
		panic(fmt.Sprintf("glox: runFile error: %v", err))
	}
	run(string(buf))

	if hadError {
		os.Exit(2)
	}
}

func runPrompt() {
	s := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		s.Scan()
		run(s.Text())
		hadError = false
	}
}

func run(s string) {
	scan := newScanner(s)
	tkns := scan.scanTokens()

	for _, t := range tkns {
		fmt.Println(t)
	}
}

func problem(line int, err error) {
	fmt.Printf("error line:%d %v\n", line+1, err)
	hadError = true
}
