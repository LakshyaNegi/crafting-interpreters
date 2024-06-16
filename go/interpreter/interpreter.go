package interpreter

import (
	"fmt"
	"glox/environment"
	"glox/generated"
	"glox/lerr"
	"glox/token"
	"os"
	"reflect"
)

type Interpreter interface {
	Interpret([]generated.Stmt)
	GetGlobalEnv() *environment.Environment
	generated.VisitorExpr
	generated.VisitorStmt
	Resolve(generated.Expr, int)
	ExecuteBlock([]generated.Stmt, *environment.Environment) (interface{}, error)
}

type interpreter struct {
	GlobalEnv *environment.Environment
	Env       *environment.Environment
	Locals    map[generated.Expr]int
}

func NewInterpreter() Interpreter {
	g := environment.NewEnvironment(nil)

	g.Define("clock", &clock{})

	in := &interpreter{
		GlobalEnv: g,
	}

	in.Env = in.GlobalEnv

	return in
}

func (i *interpreter) GetGlobalEnv() *environment.Environment {
	return i.GlobalEnv
}

func (i *interpreter) ExecuteBlock(stmts []generated.Stmt, Env *environment.Environment) (interface{}, error) {
	return i.executeBlock(stmts, Env)
}

func (i *interpreter) Resolve(expr generated.Expr, depth int) {
	i.Locals[expr] = depth
}

func (i *interpreter) Interpret(stmts []generated.Stmt) {
	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			fmt.Printf("Error while interpreting : %v\n", err)
			os.Exit(70)
		}
	}
}

func (i *interpreter) execute(stmt generated.Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *interpreter) executeBlock(stmts []generated.Stmt, Env *environment.Environment) (interface{}, error) {
	previousEnv := i.Env
	defer func(pre *environment.Environment) {
		i.Env = pre
	}(previousEnv)

	i.Env = Env

	for _, stmt := range stmts {
		_, err := i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *interpreter) VisitReturnStmt(returnstmt *generated.ReturnStmt) (interface{}, error) {
	var value interface{}
	var err error

	if returnstmt.Value != nil {
		value, err = i.evaluate(returnstmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return value, &Return{Value: value}
}

func (i *interpreter) VisitExprStmt(exprstmt *generated.ExprStmt) (interface{}, error) {
	return i.evaluate(exprstmt.Expr)
}

func (i *interpreter) VisitIfStmt(ifstmt *generated.IfStmt) (interface{}, error) {
	val, err := i.evaluate(ifstmt.Condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(val) {
		return i.execute(ifstmt.IfBranch)
	} else if ifstmt.ElseBranch != nil {
		return i.execute(ifstmt.ElseBranch)
	}

	return nil, nil
}

func (i *interpreter) VisitPrintStmt(printstmt *generated.PrintStmt) (interface{}, error) {
	value, err := i.evaluate(printstmt.Expr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", i.stringify(value))

	return nil, nil
}

func (i *interpreter) VisitBlockStmt(blockstmt *generated.BlockStmt) (interface{}, error) {
	_, err := i.executeBlock(blockstmt.Statements, environment.NewEnvironment(i.Env))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *interpreter) VisitWhileStmt(whilestmt *generated.WhileStmt) (interface{}, error) {
	c, err := i.evaluate(whilestmt.Condition)
	if err != nil {
		return nil, err
	}

	for i.isTruthy(c) {
		c, err = i.evaluate(whilestmt.Condition)
		if err != nil {
			return nil, err
		}

		_, err := i.execute(whilestmt.Stmt)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (i *interpreter) VisitVarStmt(varstmt *generated.VarStmt) (interface{}, error) {
	var value interface{}
	var err error

	if varstmt.Initializer != nil {
		value, err = i.evaluate(varstmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.Env.Define(varstmt.Name.GetLexeme(), value)
	return nil, nil
}

func (i *interpreter) VisitFunctionStmt(funstmt *generated.FunctionStmt) (interface{}, error) {
	function := &fun{Declaration: funstmt, Closure: i.Env}
	i.Env.Define(funstmt.Name.GetLexeme(), function)
	return nil, nil
}

func (i *interpreter) VisitAssign(assign *generated.Assign) (interface{}, error) {
	value, err := i.evaluate(assign.Value)
	if err != nil {
		return nil, err
	}

	dist, ok := i.Locals[assign]
	if ok {
		i.Env.AssignAt(dist, assign.Name, value)
	} else {
		i.GlobalEnv.Assign(assign.Name, value)
	}
	return value, nil
}

func (i *interpreter) VisitLogical(logical *generated.Logical) (interface{}, error) {
	left, err := i.evaluate(logical.Left)
	if err != nil {
		return nil, err
	}

	if logical.Operator.GetType() == token.OR {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(logical.Right)
}

func (i *interpreter) VisitCall(call *generated.Call) (interface{}, error) {
	callee, err := i.evaluate(call.Callee)
	if err != nil {
		return nil, err
	}

	args := make([]interface{}, 0)
	for _, arg := range call.Arguments {
		value, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}

		args = append(args, value)
	}

	function, ok := (callee).(LoxCallable)
	if !ok {
		return nil, lerr.NewRuntimeErr(call.Paren, "Can only call functions and classes.")
	}

	if len(args) != function.Arity() {
		return nil, lerr.NewRuntimeErr(call.Paren, fmt.Sprintf("Expected %d arguments but got %d.", function.Arity(), len(args)))
	}

	return function.Call(i, args)
}

func (i *interpreter) VisitBinary(binary *generated.Binary) (interface{}, error) {
	left, err := i.evaluate(binary.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(binary.Right)
	if err != nil {
		return nil, err
	}

	switch binary.Operator.GetType() {

	case token.MINUS:
		err := i.checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}

		return left.(float64) - right.(float64), nil

	case token.SLASH:
		err := i.checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil

	case token.STAR:
		err := i.checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil

	case token.PLUS:
		if reflect.TypeOf(right).Kind() == reflect.Float64 &&
			reflect.TypeOf(left).Kind() == reflect.Float64 {
			return left.(float64) + right.(float64), nil
		}

		if reflect.TypeOf(right).Kind() == reflect.String &&
			reflect.TypeOf(left).Kind() == reflect.String {
			return left.(string) + right.(string), nil
		}

		return nil, lerr.NewRuntimeErr(binary.Operator, "Operands must be two numbers or two strings.")
	case token.GREATER:
		err := i.checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil

	case token.GREATER_EQUAL:
		err := i.checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil

	case token.LESS:
		err := i.checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil

	case token.LESS_EQUAL:
		err := i.checkNumberOperands(binary.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil

	case token.BANG_EQUAL:
		return !isEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return isEqual(left, right), nil
	}

	return nil, nil
}

func (i *interpreter) VisitTernary(ternary *generated.Ternary) (interface{}, error) {
	cond, err := i.evaluate(ternary.Condition)
	if err != nil {
		return nil, err
	}
	valueTrue, err := i.evaluate(ternary.ValueTrue)
	if err != nil {
		return nil, err
	}
	valueFalse, err := i.evaluate(ternary.ValueFalse)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(cond) {
		return valueTrue, nil
	}

	return valueFalse, nil
}

func (i *interpreter) VisitGrouping(grouping *generated.Grouping) (interface{}, error) {
	return i.evaluate(grouping)
}

func (i *interpreter) VisitLiteral(literal *generated.Literal) (interface{}, error) {
	return literal.Value, nil
}

func (i *interpreter) VisitUnary(unary *generated.Unary) (interface{}, error) {
	right, err := i.evaluate(unary.Right)
	if err != nil {
		return nil, err
	}

	switch unary.Operator.GetType() {
	case token.MINUS:
		err := i.checkNumberOperands(unary.Operator, right)
		if err != nil {
			return nil, err
		}
		return -right.(float64), nil
	case token.BANG:
		return !i.isTruthy(right.(float64)), nil
	}

	return nil, nil
}

func (i *interpreter) VisitVarExpr(varexpr *generated.VarExpr) (interface{}, error) {
	return i.lookUpVariable(varexpr.Name, varexpr)
}

func (i *interpreter) lookUpVariable(name token.Token, expr generated.Expr) (interface{}, error) {
	distance, ok := i.Locals[expr]
	if ok {
		return i.Env.GetAt(distance, name.GetLexeme())
	}

	return i.GlobalEnv.Get(name)
}

func (i *interpreter) evaluate(expr generated.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *interpreter) isTruthy(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if reflect.TypeOf(obj).Kind() == reflect.Bool {
		return obj.(bool)
	}

	return true
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	return reflect.ValueOf(a).Equal(reflect.ValueOf(b))
}

func (i *interpreter) checkNumberOperands(operator token.Token, operands ...interface{}) error {
	for _, operand := range operands {
		if reflect.TypeOf(operand).Kind() != reflect.Float64 {
			return lerr.NewRuntimeErr(operator, "Operand(s) must be a number(s).")

		}
	}

	return nil
}

func (i *interpreter) stringify(obj interface{}) string {
	if obj == nil {
		return "nil"
	}

	if reflect.TypeOf(obj).Kind() == reflect.Float64 {
		// could remove .0 from end
		return fmt.Sprintf("%f", obj.(float64))
	}

	return fmt.Sprintf("%v", obj)
}
