package interpreter

import "time"

var _ LoxCallable = (*clock)(nil)

type clock struct{}

func (c clock) Arity() int {
	return 0
}

func (c clock) Call(_ Interpreter, _ []interface{}) (interface{}, error) {
	return float64(time.Now().UnixMilli()), nil
}

func (c clock) String() string {
	return "<native fn>"
}
