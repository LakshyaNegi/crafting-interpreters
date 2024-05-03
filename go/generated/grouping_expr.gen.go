
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

func (x *Grouping) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitGrouping(x)
}
