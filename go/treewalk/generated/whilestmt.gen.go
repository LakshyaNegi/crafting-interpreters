
// generated code - DO NOT EDIT
package generated

import (
)

type WhileStmt struct {
	Condition Expr
	Stmt Stmt
}

func NewWhileStmt(
	Condition Expr,
	Stmt Stmt,
) *WhileStmt {
	return &WhileStmt {
		Condition: Condition,
		Stmt: Stmt,
	}
}

func (x *WhileStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitWhileStmt(x)
}
