
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type FunctionStmt struct {
	Name token.Token
	Params []token.Token
	Body []Stmt
}

func NewFunctionStmt(
	Name token.Token,
	Params []token.Token,
	Body []Stmt,
) *FunctionStmt {
	return &FunctionStmt {
		Name: Name,
		Params: Params,
		Body: Body,
	}
}

func (x *FunctionStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitFunctionStmt(x)
}
