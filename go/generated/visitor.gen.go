
// generated code - DO NOT EDIT
package generated

type Visitor interface {
	VisitBinary (binary *Binary) (interface{}, error)
	VisitTernary (ternary *Ternary) (interface{}, error)
	VisitGrouping (grouping *Grouping) (interface{}, error)
	VisitLiteral (literal *Literal) (interface{}, error)
	VisitUnary (unary *Unary) (interface{}, error)
}
