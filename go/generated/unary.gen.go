
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type Unary struct {
	Operator token.Token
	Right Expr
}

func NewUnary(
	Operator token.Token,
	Right Expr,
) *Unary {
	return &Unary {
		Operator: Operator,
		Right: Right,
	}
}

func (x *Unary) Accept(visitor Visitor) interface{} {
	return visitor.VisitUnary(x)
}
