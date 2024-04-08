
// generated code - DO NOT EDIT
package expr

import (
	
	"glox/expr"
	
	"glox/token"
	
)

type Unary struct {
	expr.Expr
	
	Operator token.Token
	
	Right expr.Expr
	
}

func NewUnary(
	
	Operator token.Token,
	
	Right expr.Expr,
	
) *Unary {
	return &Unary{
		
		Operator: Operator,
		
		Right: Right,
		
	}
}
