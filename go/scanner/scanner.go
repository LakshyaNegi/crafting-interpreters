package scanner

import (
	"glox/err"
	"glox/token"
	"strconv"
)

type Scanner interface {
	ScanTokens() ([]token.Token, error)
}

type scanner struct {
	source  string
	tokens  []token.Token
	start   int
	current int
	line    int
}

func NewScanner(source string) Scanner {
	return &scanner{
		source: source,
		tokens: make([]token.Token, 0),
	}
}

func (s *scanner) ScanTokens() ([]token.Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", nil, s.line))

	return s.tokens, nil
}

func (s *scanner) scanToken() error {
	c := s.advance()
	switch c {

	// single character
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)

	// logical operators
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL, nil)
		} else {
			s.addToken(token.BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL, nil)
		} else {
			s.addToken(token.EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL, nil)
		} else {
			s.addToken(token.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL, nil)
		} else {
			s.addToken(token.GREATER, nil)
		}

	// slash or comment
	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}

	// meaningless chars
	case ' ':
	case '\r':
	case '\t':
		// Ignore whitespace.
		break

	case '\n':
		s.line++

	// string
	case '"':
		s.strScan()
	default:
		if s.isDigit(c) {
			s.numScan()
		} else if s.isAlpha(c) {
			s.idenScan()
		} else {

			return err.NewSyntaxErr(s.line, "", "Unexpected character.")
		}
	}

	return nil
}

func (s *scanner) addToken(tokenType token.TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.NewToken(tokenType, text, literal, s.line))
}

func (s *scanner) strScan() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}

		s.advance()
	}

	s.advance()

	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

func (s *scanner) numScan() error {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()

		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		return err
	}

	s.addToken(token.NUMBER, num)

	return nil
}

func (s *scanner) idenScan() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}

	tokenType := token.IDENTIFIER

	text := s.source[s.start:s.current]
	if _, ok := token.Keywords[text]; ok {
		tokenType = token.Keywords[text]
	}

	s.addToken(tokenType, nil)
}

func (s *scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != byte(expected) {
		return false
	}

	defer func() {
		s.current++
	}()

	return true
}

func (s *scanner) advance() rune {
	defer func() {
		s.current++
	}()

	return rune(s.source[s.current])
}

func (s *scanner) peek() rune {
	if s.isAtEnd() {
		return rune(0)
	}

	return rune(s.source[s.current])
}

func (s *scanner) peekNext() rune {
	if s.current+1 > len(s.source) {
		return rune(0)
	}

	return rune(s.source[s.current+1])
}

func (s *scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *scanner) isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *scanner) isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		(c == '_')
}

func (s *scanner) isAlphaNumeric(c rune) bool {
	return s.isAlpha(c) || s.isDigit(c)
}
