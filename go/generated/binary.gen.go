
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type Binary struct {
	Expr
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinary(
	Left Expr,
	Operator token.Token,
	Right Expr,
) *Binary {
	return &Binary{
		Left: Left,
		Operator: Operator,
		Right: Right,
	}
}

func (x *Binary) accept(visitor Visitor[any]) any {
	return visitor.visitBinary(x)
}
