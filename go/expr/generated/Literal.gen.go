
// generated code - DO NOT EDIT
package expr

import (
	
	"glox/expr"
	
)

type Literal struct {
	expr.Expr
	
	Value interface{}
	
}

func NewLiteral(
	
	Value interface{},
	
) *Literal {
	return &Literal{
		
		Value: Value,
		
	}
}
