
// generated code - DO NOT EDIT
package generated

import (
)

type Grouping struct {
	Expression Expr
}

func NewGrouping(
	Expression Expr,
) *Grouping {
	return &Grouping {
		Expression: Expression,
	}
}

func (x *Grouping) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitGrouping(x)
}
