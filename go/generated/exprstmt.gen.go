
// generated code - DO NOT EDIT
package generated

import (
)

type ExprStmt struct {
	Expr Expr
}

func NewExprStmt(
	Expr Expr,
) *ExprStmt {
	return &ExprStmt {
		Expr: Expr,
	}
}

func (x *ExprStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitExprStmt(x)
}
