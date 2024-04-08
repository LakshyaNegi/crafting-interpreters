
// generated code - DO NOT EDIT
package generated

import (
)

type Literal struct {
	Expr
	Value interface{}
}

func NewLiteral(
	Value interface{},
) *Literal {
	return &Literal{
		Value: Value,
	}
}

func (x *Literal) accept(visitor Visitor[any]) any {
	return visitor.visitLiteral(x)
}
