
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type VarStmt struct {
	Name token.Token
	Initializer Expr
}

func NewVarStmt(
	Name token.Token,
	Initializer Expr,
) *VarStmt {
	return &VarStmt {
		Name: Name,
		Initializer: Initializer,
	}
}

func (x *VarStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitVarStmt(x)
}
