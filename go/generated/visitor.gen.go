
// generated code - DO NOT EDIT
package generated

type Visitor[T any] interface {
	visitBinary (binary Expr) T
	visitGrouping (grouping Expr) T
	visitLiteral (literal Expr) T
	visitUnary (unary Expr) T
}
