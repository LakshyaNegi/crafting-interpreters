
// generated code - DO NOT EDIT
package generated

import (
)

type Ternary struct {
	Condition Expr
	ValueTrue Expr
	ValueFalse Expr
}

func NewTernary(
	Condition Expr,
	ValueTrue Expr,
	ValueFalse Expr,
) *Ternary {
	return &Ternary {
		Condition: Condition,
		ValueTrue: ValueTrue,
		ValueFalse: ValueFalse,
	}
}

func (x *Ternary) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitTernary(x)
}
