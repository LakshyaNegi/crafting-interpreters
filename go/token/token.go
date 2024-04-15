package token

import "fmt"

type Token interface {
	Show() string
	GetLexeme() string
}

type token struct {
	tokenType TokenType
	lexeme    string
	literal   interface{}
	line      int
}

func (t *token) Show() string {
	return fmt.Sprintf("%v : %s : %v\n", t.tokenType, t.lexeme, t.literal)
}

func (t *token) GetLexeme() string {
	return t.lexeme
}

func NewToken(
	tokenType TokenType,
	lexeme string,
	literal interface{},
	line int,
) Token {
	return &token{
		tokenType: tokenType,
		lexeme:    lexeme,
		literal:   literal,
		line:      line,
	}
}
