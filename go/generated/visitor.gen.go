
// generated code - DO NOT EDIT
package generated

type Visitor interface {
	VisitBinary (binary *Binary) interface{}
	VisitTernary (ternary *Ternary) interface{}
	VisitGrouping (grouping *Grouping) interface{}
	VisitLiteral (literal *Literal) interface{}
	VisitUnary (unary *Unary) interface{}
}
