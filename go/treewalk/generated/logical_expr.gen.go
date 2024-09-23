
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type Logical struct {
	Left Expr
	Operator token.Token
	Right Expr
}

func NewLogical(
	Left Expr,
	Operator token.Token,
	Right Expr,
) *Logical {
	return &Logical {
		Left: Left,
		Operator: Operator,
		Right: Right,
	}
}

func (x *Logical) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitLogical(x)
}
