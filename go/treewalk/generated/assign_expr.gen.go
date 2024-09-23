
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type Assign struct {
	Name token.Token
	Value Expr
}

func NewAssign(
	Name token.Token,
	Value Expr,
) *Assign {
	return &Assign {
		Name: Name,
		Value: Value,
	}
}

func (x *Assign) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitAssign(x)
}
