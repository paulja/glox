package glox

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// LineError represents an error that has a line number, typically used for
// parse errors. Implements the error interface.
type LineError struct {
	Line int

	error error
}

func (se LineError) Error() string {
	return se.error.Error()
}

func newLineError(n int, err error) LineError {
	return LineError{n, err}
}

// Scanner is a Lox language token scanner. The scanner sends Token type values
// down the T channel and error down the E channel. The Done channel is
// signalled when the scanner has finished processing the source.
type Scanner struct {
	source string
	reader *strings.Reader

	E    chan error
	T    chan *Token
	Done chan bool

	start   int
	current int
	line    int
}

// NewScanner creates a new Scanner.
func NewScanner(s string) *Scanner {
	r := strings.NewReader(s)
	return &Scanner{
		source: s,
		reader: r,
		E:      make(chan error),
		T:      make(chan *Token),
		Done:   make(chan bool),
	}
}

// Scan processes the source input string, emitting tokens and errors as
// appropriate, signalling the Done channel when complete.
func (s *Scanner) Scan() {
	for !s.atEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.T <- NewToken(TokenEOF, "EOF", nil, s.line)
	s.Done <- true
}

func (s *Scanner) atEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	ch := s.advance()
	switch ch {
	case '(':
		s.addToken(TokenLeftParen, nil)
	case ')':
		s.addToken(TokenRightParen, nil)
	case '{':
		s.addToken(TokenLeftBrace, nil)
	case '}':
		s.addToken(TokenRightBrace, nil)
	case ',':
		s.addToken(TokenComma, nil)
	case '.':
		s.addToken(TokenDot, nil)
	case '-':
		s.addToken(TokenMinus, nil)
	case '+':
		s.addToken(TokenPlus, nil)
	case ';':
		s.addToken(TokenSemicolon, nil)
	case '*':
		s.addToken(TokenStar, nil)
	case '!':
		s.addToken(ifToken(s.match('='), TokenBangEqual, TokenBang), nil)
	case '=':
		s.addToken(ifToken(s.match('='), TokenEqualEqual, TokenEqual), nil)
	case '<':
		s.addToken(ifToken(s.match('='), TokenLessEqual, TokenLess), nil)
	case '>':
		s.addToken(ifToken(s.match('='), TokenGreaterEqual, TokenLess), nil)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.atEnd() {
				s.advance()
			}
			break
		}
		s.addToken(TokenSlash, nil)
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
		s.E <- newLineError(s.line, fmt.Errorf("Unexpected character: %q", ch))
	}
}

func ifToken(c bool, a TokenType, b TokenType) TokenType {
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

func (s *Scanner) advance() rune {
	ch, sz, err := s.reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("advance error: %v", err))
	}
	s.current += sz
	return ch
}

func (s *Scanner) match(ex rune) bool {
	if s.atEnd() {
		return false
	}

	ch, sz, err := s.reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("match error: %v", err))
	}
	if ch != ex {
		if err := s.reader.UnreadRune(); err != nil {
			panic(fmt.Sprintf("match unread error: %v", err))
		}
		return false
	}

	s.current += sz
	return true
}

func (s *Scanner) peek() rune {
	if s.atEnd() {
		return 0
	}
	ch, _, err := s.reader.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("peek error: %v", err))
	}
	if err := s.reader.UnreadRune(); err != nil {
		panic(fmt.Sprintf("peek unread error: %v", err))
	}
	return ch
}

func (s *Scanner) peekNext() rune {
	var (
		i   int
		err error
		ch  rune
	)

	// scan forward
	for ; i < 2 && !s.atEnd() && err == nil; i++ {
		ch, _, err = s.reader.ReadRune()
	}

	// unwind
	if _, err = s.reader.Seek(int64(s.current), 0); err != nil {
		panic(fmt.Sprintf("peekNext read error: %v", err))
	}

	// failed to peek n
	if i < 2 {
		return 0
	}

	return ch
}

func (s *Scanner) addToken(t TokenType, lit interface{}) {
	lex := s.source[s.start:s.current]
	s.T <- NewToken(t, lex, lit, s.line)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.atEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.atEnd() {
		s.E <- newLineError(s.line, errors.New("string error: unterminated string"))
		return
	}

	s.advance()

	v := s.source[s.start+1 : s.current-1]
	s.addToken(TokenString, v)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	v := s.source[s.start:s.current]
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		s.E <- newLineError(s.line, fmt.Errorf("number error: %v", err))
		return
	}
	s.addToken(TokenNumber, f)
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	v := s.source[s.start:s.current]
	if i, ok := keywords[v]; ok {
		s.addToken(i, nil)
		return
	}

	s.addToken(TokenIdentifier, nil)
}
