
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type Call struct {
	Callee Expr
	Paren token.Token
	Arguments []Expr
}

func NewCall(
	Callee Expr,
	Paren token.Token,
	Arguments []Expr,
) *Call {
	return &Call {
		Callee: Callee,
		Paren: Paren,
		Arguments: Arguments,
	}
}

func (x *Call) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitCall(x)
}
