
// generated code - DO NOT EDIT
package generated

import (
)

type Grouping struct {
	Expr
	Expression Expr
}

func NewGrouping(
	Expression Expr,
) *Grouping {
	return &Grouping{
		Expression: Expression,
	}
}

func (x *Grouping) accept(visitor Visitor[any]) any {
	return visitor.visitGrouping(x)
}
