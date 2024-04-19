package ast

import (
	"fmt"
	"glox/generated"
)

type Printer interface {
	Print(expr generated.Expr) string
}

type printer struct{}

func NewPrinter() Printer {
	return &printer{}
}

func (p *printer) Print(expr generated.Expr) string {
	return fmt.Sprintf("%v", expr.Accept(p))
}

func (p *printer) VisitBinary(expr *generated.Binary) interface{} {
	return p.parenthesize(expr.Operator.GetLexeme(), expr.Left, expr.Right)
}

func (p *printer) VisitTernary(expr *generated.Ternary) interface{} {
	return p.parenthesize("ternary", expr.Condition, expr.ValueTrue, expr.ValueFalse)
}

func (p *printer) VisitGrouping(expr *generated.Grouping) interface{} {
	return p.parenthesize("group", expr.Expression)
}

func (p *printer) VisitLiteral(expr *generated.Literal) interface{} {
	return expr.Value
}

func (p *printer) VisitUnary(expr *generated.Unary) interface{} {
	return p.parenthesize(expr.Operator.GetLexeme(), expr.Right)
}

func (p *printer) parenthesize(name string, exprs ...generated.Expr) interface{} {
	str := ""

	str += "("
	str += name

	for _, expr := range exprs {
		str += " "
		str += fmt.Sprint(expr.Accept(p))
	}

	str += ")"

	return str
}
