package parser

import (
	"fmt"
	"glox/err"
	"glox/generated"
	"glox/token"
)

type Parser interface {
	Parse() (generated.Expr, error)
}

type parser struct {
	tokens []token.Token
	curr   int
	perr   err.ParseErr
}

func NewParser(tokens []token.Token) Parser {
	return &parser{
		tokens: tokens,
		curr:   0,
	}
}

func (p *parser) Parse() (generated.Expr, error) {
	return p.expression()
}

func (p *parser) expression() (generated.Expr, error) {
	return p.ternary()
}

func (p *parser) ternary() (generated.Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.QUESTION) {
		exprTrue, err := p.expression()
		if err != nil {
			return nil, err
		}

		p.consume(token.COLON, "Expect : after expression in ternary operation")

		exprFalse, err := p.expression()
		if err != nil {
			return nil, err
		}

		expr = generated.NewTernary(expr, exprTrue, exprFalse)
	}

	return expr, nil
}

func (p *parser) equality() (generated.Expr, error) {
	expr, err := p.comparsion()
	if err != nil {
		return nil, err
	}

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right, err := p.comparsion()
		if err != nil {
			return nil, err
		}
		expr = generated.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) comparsion() (generated.Expr, error) {
	expr, err := p.terminal()
	if err != nil {
		return nil, err
	}

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.terminal()
		if err != nil {
			return nil, err
		}
		expr = generated.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) terminal() (generated.Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(token.PLUS, token.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}
		expr = generated.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) factor() (generated.Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}
		expr = generated.NewBinary(expr, operator, right)
	}

	return expr, nil
}

func (p *parser) unary() (generated.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return generated.NewUnary(operator, right), nil
	}

	return p.primary()
}

func (p *parser) primary() (generated.Expr, error) {
	if p.match(token.FALSE) {
		return generated.NewLiteral(false), nil
	}

	if p.match(token.TRUE) {
		return generated.NewLiteral(true), nil
	}

	if p.match(token.NIL) {
		return generated.NewLiteral(nil), nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return generated.NewLiteral(p.previous().GetLiteral()), nil
	}

	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after expression")
		if err != nil {
			return nil, err
		}

		return generated.NewGrouping(expr), nil
	}

	return nil, p.perror(p.peek(), "Expect expression.")
}

func (p *parser) consume(t token.TokenType, msg string) (token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}

	serr := p.perror(p.peek(), msg)
	if serr != nil {
		fmt.Print(serr.Error())
	}

	return nil, err.NewParseErr()
}

func (p *parser) perror(t token.Token, msg string) error {
	if t.GetType() == token.EOF {
		return err.NewSyntaxErr(p.peek().GetLine(), " at end", msg)
	} else {
		return err.NewSyntaxErr(p.peek().GetLine(), " at end", msg)
	}
}

func (p *parser) match(tt ...token.TokenType) bool {
	for _, t := range tt {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().GetType() == t
}

func (p *parser) advance() token.Token {
	if !p.isAtEnd() {
		p.curr++
	}

	return p.previous()
}

func (p *parser) previous() token.Token {
	return p.tokens[p.curr-1]
}

func (p *parser) isAtEnd() bool {
	return p.peek().GetType() == token.EOF
}

func (p *parser) peek() token.Token {
	return p.tokens[p.curr]
}

func (p *parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().GetType() == token.SEMICOLON {
			return
		}

		switch p.peek().GetType() {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
