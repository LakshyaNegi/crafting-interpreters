
// generated code - DO NOT EDIT
package generated

import (
)

type Literal struct {
	Value interface{}
}

func NewLiteral(
	Value interface{},
) *Literal {
	return &Literal {
		Value: Value,
	}
}

func (x *Literal) Accept(visitor Visitor) interface{} {
	return visitor.VisitLiteral(x)
}
