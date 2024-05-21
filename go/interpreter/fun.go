package interpreter

import (
	"glox/environment"
	"glox/generated"
)

var _ LoxCallable = (*fun)(nil)

type fun struct {
	Declaration *generated.FunctionStmt
	Closure     *environment.Environment
}

func (f fun) Arity() int {
	return len(f.Declaration.Params)
}

func (f fun) Call(in Interpreter, args []interface{}) (interface{}, error) {
	env := environment.NewEnvironment(f.Closure)

	for i, param := range f.Declaration.Params {
		env.Define(param.GetLexeme(), args[i])
	}

	_, err := in.ExecuteBlock(f.Declaration.Body, env)
	if err != nil {
		if _, ok := err.(*Return); ok {
			return err.(*Return).Value, nil
		}

		return nil, err
	}

	return nil, nil
}

func (f fun) String() string {
	return "<native fn " + f.Declaration.Name.GetLexeme() + ">"
}
