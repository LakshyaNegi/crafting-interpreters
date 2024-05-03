
// generated code - DO NOT EDIT
package generated

type VisitorExpr interface {
	VisitAssign (assign *Assign) (interface{}, error)
	VisitLogical (logical *Logical) (interface{}, error)
	VisitBinary (binary *Binary) (interface{}, error)
	VisitTernary (ternary *Ternary) (interface{}, error)
	VisitGrouping (grouping *Grouping) (interface{}, error)
	VisitLiteral (literal *Literal) (interface{}, error)
	VisitUnary (unary *Unary) (interface{}, error)
	VisitVarExpr (varexpr *VarExpr) (interface{}, error)
}
