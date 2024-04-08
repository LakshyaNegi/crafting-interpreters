
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type Unary struct {
	Expr
	Operator token.Token
	Right Expr
}

func NewUnary(
	Operator token.Token,
	Right Expr,
) *Unary {
	return &Unary{
		Operator: Operator,
		Right: Right,
	}
}

func (x *Unary) accept(visitor Visitor[any]) any {
	return visitor.visitUnary(x)
}
