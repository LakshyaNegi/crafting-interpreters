
// generated code - DO NOT EDIT
package generated

import (
)

type PrintStmt struct {
	Expr Expr
}

func NewPrintStmt(
	Expr Expr,
) *PrintStmt {
	return &PrintStmt {
		Expr: Expr,
	}
}

func (x *PrintStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitPrintStmt(x)
}
