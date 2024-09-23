
// generated code - DO NOT EDIT
package generated

import (
)

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(
	Statements []Stmt,
) *BlockStmt {
	return &BlockStmt {
		Statements: Statements,
	}
}

func (x *BlockStmt) Accept(visitor VisitorStmt) (interface{}, error) {
	return visitor.VisitBlockStmt(x)
}
