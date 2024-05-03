package parser

import (
	"fmt"

	"glox/generated"
	"glox/lerr"
	"glox/token"
)

type Parser interface {
	Parse() ([]generated.Stmt, error)
}

type parser struct {
	tokens []token.Token
	curr   int
}

func NewParser(tokens []token.Token) Parser {
	return &parser{
		tokens: tokens,
		curr:   0,
	}
}

func (p *parser) Parse() ([]generated.Stmt, error) {
	var stmts []generated.Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}

	return stmts, nil
}

func (p *parser) declaration() (generated.Stmt, error) {
	if p.match(token.VAR) {
		vd, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, nil
		}

		return vd, nil
	}

	return p.statement()
}

func (p *parser) varDeclaration() (generated.Stmt, error) {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.\n")
	if err != nil {
		return nil, err
	}

	var initializer generated.Expr
	if p.match(token.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	p.consume(token.SEMICOLON, "Expect ; after variable declaration.\n")

	return generated.NewVarStmt(name, initializer), nil
}

func (p *parser) statement() (generated.Stmt, error) {
	if p.match(token.FOR) {
		return p.forStmt()
	} else if p.match(token.IF) {
		return p.ifStmt()
	} else if p.match(token.PRINT) {
		return p.printStmt()
	} else if p.match(token.WHILE) {
		return p.whileStmt()
	} else if p.match(token.LEFT_BRACE) {
		block, err := p.blockstmt()
		if err != nil {
			return nil, err
		}

		return generated.NewBlockStmt(block), nil
	}

	return p.expressionStmt()
}

func (p *parser) blockstmt() ([]generated.Stmt, error) {
	stmts := []generated.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		d, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, d)
	}

	p.consume(token.RIGHT_BRACE, "Expect '}' after block.")

	return stmts, nil
}

func (p *parser) ifStmt() (generated.Stmt, error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after if.")

	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.LEFT_PAREN, "Expect ')' after if condition.")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch generated.Stmt
	if p.match(token.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return generated.NewIfStmt(condition, thenBranch, elseBranch), nil
}

func (p *parser) whileStmt() (generated.Stmt, error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}

	p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	return generated.NewWhileStmt(condition, body), nil
}

func (p *parser) forStmt() (generated.Stmt, error) {
	p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")

	var initializer generated.Stmt
	var err error
	if p.match(token.SEMICOLON) {
		initializer = nil
	} else if p.match(token.VAR) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStmt()
	}
	if err != nil {
		return nil, err
	}

	var condition generated.Expr
	if !p.check(token.SEMICOLON) {
		condition, err = p.expression()
	}
	if err != nil {
		return nil, err
	}

	p.consume(token.SEMICOLON, "Expect ';' after loop condition.")

	var increment generated.Expr
	if !p.check(token.RIGHT_PAREN) {
		increment, err = p.expression()
	}
	if err != nil {
		return nil, err
	}

	p.consume(token.RIGHT_PAREN, "Expect ')' after loop clauses.")

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = generated.NewBlockStmt(
			[]generated.Stmt{
				body, generated.NewExprStmt(increment),
			})
	}

	if condition == nil {
		condition = generated.NewLiteral(true)
	}

	body = generated.NewWhileStmt(condition, body)

	if initializer != nil {
		body = generated.NewBlockStmt([]generated.Stmt{initializer, body})
	}

	return body, nil
}

func (p *parser) printStmt() (generated.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.SEMICOLON, "Expect ; after value.")
	if err != nil {
		return nil, err
	}

	return generated.NewPrintStmt(expr), nil
}

func (p *parser) expressionStmt() (generated.Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}

	_, err = p.consume(token.SEMICOLON, "Expect ; after expression.")
	if err != nil {
		return nil, err
	}

	return generated.NewExprStmt(expr), nil
}

func (p *parser) expression() (generated.Expr, error) {
	return p.assignment()
}

func (p *parser) assignment() (generated.Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if v, ok := expr.(*generated.VarExpr); ok {
			name := v.Name
			return generated.NewAssign(name, value), nil
		}

		return nil, p.perror(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *parser) ternary() (generated.Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	for p.match(token.QUESTION) {
		exprTrue, err := p.expression()
		if err != nil {
			return nil, err
		}

		_, err = p.consume(token.COLON, "Expect : after expression in ternary operation")
		if err != nil {
			return nil, err
		}

		exprFalse, err := p.expression()
		if err != nil {
			return nil, err
		}

		expr = generated.NewTernary(expr, exprTrue, exprFalse)
	}

	return expr, nil
}

func (p *parser) or() (generated.Expr, error) {
	left, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(token.OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		left = generated.NewLogical(left, operator, right)
	}

	return left, nil
}

func (p *parser) and() (generated.Expr, error) {
	left, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(token.AND) {
		operator := p.previous()

		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		left = generated.NewLogical(left, operator, right)
	}

	return left, nil
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

	if p.match(token.IDENTIFIER) {
		return generated.NewVarExpr(p.previous()), nil
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

	return nil, lerr.NewParseErr()
}

func (p *parser) perror(t token.Token, msg string) error {
	if t.GetType() == token.EOF {
		return lerr.NewSyntaxErr(p.peek().GetLine(), " at end", msg)
	} else {
		return lerr.NewSyntaxErr(p.peek().GetLine(), " at end", msg)
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
