
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type Binary struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewBinary(
	Left Expr,
	Operator token.Token,
	Right Expr,
) *Binary {
	return &Binary {
		Left: Left,
		Operator: Operator,
		Right: Right,
	}
}

func (x *Binary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitBinary(x)
}
