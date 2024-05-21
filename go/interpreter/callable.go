package interpreter

type LoxCallable interface {
	Arity() int
	Call(Interpreter, []interface{}) (interface{}, error)
	String() string
}
