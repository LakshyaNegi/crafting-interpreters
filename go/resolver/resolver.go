package resolver

import (
	"glox/generated"
	"glox/interpreter"
	"glox/lerr"
	"glox/token"
)

type functionType string

const (
	FunctionTypeNone     functionType = "none"
	FunctionTypeFunction functionType = "function"
)

type Resolver interface {
	generated.VisitorStmt
	generated.VisitorExpr
	Resolve([]generated.Stmt) error
}

type resolver struct {
	interpreter  interpreter.Interpreter
	scopes       []map[string]bool
	currFunction functionType
}

func NewResolver(in interpreter.Interpreter) Resolver {
	return &resolver{
		interpreter:  in,
		scopes:       []map[string]bool{},
		currFunction: FunctionTypeNone,
	}
}

func (r *resolver) VisitBlockStmt(stmt *generated.BlockStmt) (interface{}, error) {
	r.beginScope()

	err := r.resolveStmts(stmt.Statements)
	if err != nil {
		return nil, err
	}

	r.endScope()

	return nil, nil
}

func (r *resolver) Resolve(stmts []generated.Stmt) error {
	return r.resolveStmts(stmts)
}

func (r *resolver) resolveStmts(stmts []generated.Stmt) error {
	for _, stmt := range stmts {
		_, err := r.resolveStmt(stmt)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *resolver) resolveStmt(stmt generated.Stmt) (interface{}, error) {
	return stmt.Accept(r)
}

func (r *resolver) resolveExpr(expr generated.Expr) (interface{}, error) {
	return expr.Accept(r)
}

func (r *resolver) beginScope() {
	r.scopes = append(r.scopes, map[string]bool{})
}

func (r *resolver) endScope() {
	r.scopes = r.scopes[:len(r.scopes)-1]
}

func (r *resolver) VisitVarStmt(stmt *generated.VarStmt) (interface{}, error) {
	err := r.declare(stmt.Name)
	if err != nil {
		return nil, err
	}

	if stmt.Initializer != nil {
		_, err := r.resolveExpr(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	r.define(stmt.Name)

	return nil, nil
}

func (r *resolver) declare(name token.Token) error {
	if len(r.scopes) == 0 {
		return nil
	}

	scope := r.scopes[len(r.scopes)-1]
	if _, ok := scope[name.GetLexeme()]; ok {
		return lerr.NewSyntaxErr(name.GetLine(), "", "Variable with this name already declared in this scope.")
	}
	scope[name.GetLexeme()] = false

	return nil
}

func (r *resolver) define(name token.Token) {
	if len(r.scopes) == 0 {
		return
	}

	scope := r.scopes[len(r.scopes)-1]
	scope[name.GetLexeme()] = true
}

func (r *resolver) VisitVarExpr(expr *generated.VarExpr) (interface{}, error) {
	if len(r.scopes) != 0 {
		scope := r.scopes[len(r.scopes)-1]
		if _, ok := scope[expr.Name.GetLexeme()]; ok && !scope[expr.Name.GetLexeme()] {
			return nil, lerr.NewRuntimeErr(expr.Name, "Cannot read local variable in its own initializer.")
		}
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *resolver) resolveLocal(expr generated.Expr, name token.Token) {
	for i := len(r.scopes) - 1; i >= 0; i-- {
		if _, ok := r.scopes[i][name.GetLexeme()]; ok {
			r.interpreter.Resolve(expr, len(r.scopes)-1-i)
			return
		}
	}
}

func (r *resolver) VisitAssign(expr *generated.Assign) (interface{}, error) {
	_, err := r.resolveExpr(expr.Value)
	if err != nil {
		return nil, err
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *resolver) VisitFunctionStmt(stmt *generated.FunctionStmt) (interface{}, error) {
	err := r.declare(stmt.Name)
	if err != nil {
		return nil, err
	}

	r.define(stmt.Name)

	err = r.resolveFunction(stmt, FunctionTypeFunction)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *resolver) resolveFunction(stmt *generated.FunctionStmt, typ functionType) error {
	enclosingFunction := r.currFunction
	r.currFunction = typ

	r.beginScope()
	for _, param := range stmt.Params {
		err := r.declare(stmt.Name)
		if err != nil {
			return err
		}
		r.define(param)
	}

	err := r.resolveStmts(stmt.Body)
	if err != nil {
		return nil
	}

	r.endScope()

	r.currFunction = enclosingFunction

	return nil
}

func (r *resolver) VisitExprStmt(stmt *generated.ExprStmt) (interface{}, error) {
	return r.resolveExpr(stmt.Expr)
}

func (r *resolver) VisitIfStmt(stmt *generated.IfStmt) (interface{}, error) {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveStmt(stmt.IfBranch)
	if err != nil {
		return nil, err
	}

	if stmt.ElseBranch != nil {
		_, err = r.resolveStmt(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *resolver) VisitWhileStmt(stmt *generated.WhileStmt) (interface{}, error) {
	_, err := r.resolveExpr(stmt.Condition)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveStmt(stmt.Stmt)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *resolver) VisitReturnStmt(stmt *generated.ReturnStmt) (interface{}, error) {
	if r.currFunction == FunctionTypeNone {
		return nil, lerr.NewSyntaxErr(stmt.Keyword.GetLine(), "", "Cannot return from top-level code.")
	}

	if stmt.Value != nil {
		_, err := r.resolveExpr(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (r *resolver) VisitPrintStmt(stmt *generated.PrintStmt) (interface{}, error) {
	return r.resolveExpr(stmt.Expr)
}

func (r *resolver) VisitLogical(expr *generated.Logical) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *resolver) VisitBinary(expr *generated.Binary) (interface{}, error) {
	_, err := r.resolveExpr(expr.Left)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.Right)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *resolver) VisitTernary(expr *generated.Ternary) (interface{}, error) {
	_, err := r.resolveExpr(expr.Condition)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.ValueTrue)
	if err != nil {
		return nil, err
	}

	_, err = r.resolveExpr(expr.ValueFalse)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *resolver) VisitGrouping(expr *generated.Grouping) (interface{}, error) {
	return r.resolveExpr(expr.Expression)
}

func (r *resolver) VisitLiteral(expr *generated.Literal) (interface{}, error) {
	return nil, nil
}

func (r *resolver) VisitUnary(expr *generated.Unary) (interface{}, error) {
	return r.resolveExpr(expr.Right)
}

func (r *resolver) VisitCall(expr *generated.Call) (interface{}, error) {
	_, err := r.resolveExpr(expr.Callee)
	if err != nil {
		return nil, err
	}

	for _, arg := range expr.Arguments {
		_, err := r.resolveExpr(arg)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
