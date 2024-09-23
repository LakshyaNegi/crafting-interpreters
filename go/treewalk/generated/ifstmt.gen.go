
// generated code - DO NOT EDIT
package generated

import (
)

type IfStmt struct {
	Condition Expr
	IfBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(
	Condition Expr,
	IfBranch Stmt,
	ElseBranch Stmt,
) *IfStmt {
	return &IfStmt {
		Condition: Condition,
		IfBranch: IfBranch,
		ElseBranch: ElseBranch,
	}
}

func (x *IfStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitIfStmt(x)
}
