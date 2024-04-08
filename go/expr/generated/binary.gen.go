
// generated code - DO NOT EDIT
package expr

import (
	
	"glox/expr"
	
	"glox/token"
	
)

type Binary struct {
	expr.Expr
	
	Left expr.Expr
	
	Operator token.Token
	
	Right expr.Expr
	
}

func NewBinary(
	
	Left expr.Expr,
	
	Operator token.Token,
	
	Right expr.Expr,
	
) *Binary {
	return &Binary{
		
		Left: Left,
		
		Operator: Operator,
		
		Right: Right,
		
	}
}
