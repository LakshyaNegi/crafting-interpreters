
// generated code - DO NOT EDIT
package generated

type VisitorStmt interface {
	VisitBlockStmt (blockstmt *BlockStmt) (interface{}, error)
	VisitIfStmt (ifstmt *IfStmt) (interface{}, error)
	VisitWhileStmt (whilestmt *WhileStmt) (interface{}, error)
	VisitExprStmt (exprstmt *ExprStmt) (interface{}, error)
	VisitPrintStmt (printstmt *PrintStmt) (interface{}, error)
	VisitVarStmt (varstmt *VarStmt) (interface{}, error)
}
