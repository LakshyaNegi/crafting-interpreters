
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type VarExpr struct {
	Name token.Token
}

func NewVarExpr(
	Name token.Token,
) *VarExpr {
	return &VarExpr {
		Name: Name,
	}
}

func (x *VarExpr) Accept(visitor VisitorExpr) (interface{}, error) {
	return visitor.VisitVarExpr(x)
}
