
// generated code - DO NOT EDIT
package expr

import (
	
	"glox/expr"
	
)

type Grouping struct {
	expr.Expr
	
	Expression expr.Expr
	
}

func NewGrouping(
	
	Expression expr.Expr,
	
) *Grouping {
	return &Grouping{
		
		Expression: Expression,
		
	}
}
