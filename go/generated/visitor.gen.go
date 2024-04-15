
// generated code - DO NOT EDIT
package generated

type Visitor interface {
	VisitBinary (binary *Binary) interface{}
	VisitGrouping (grouping *Grouping) interface{}
	VisitLiteral (literal *Literal) interface{}
	VisitUnary (unary *Unary) interface{}
}
