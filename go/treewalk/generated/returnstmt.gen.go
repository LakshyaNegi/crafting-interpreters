
// generated code - DO NOT EDIT
package generated

import (
	"glox/token"
)

type ReturnStmt struct {
	Keyword token.Token
	Value Expr
}

func NewReturnStmt(
	Keyword token.Token,
	Value Expr,
) *ReturnStmt {
	return &ReturnStmt {
		Keyword: Keyword,
		Value: Value,
	}
}

func (x *ReturnStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitReturnStmt(x)
}
