package interpreter

import (
	"fmt"
	"glox/generated"
	"glox/lerr"
	"glox/token"
	"os"
	"reflect"
)

type Interpreter interface {
	Interpret(generated.Expr)
}

type interpreter struct {
}

func NewInterpreter() Interpreter {
	return &interpreter{}
}

func (i *interpreter) Interpret(expr generated.Expr) {
	value, err := i.evaluate(expr)
	if err != nil {
		fmt.Printf("Error while interpreting : %v\n", err)
		os.Exit(70)
	}

	fmt.Printf("%v\n", i.stringify(value))
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
