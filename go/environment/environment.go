package environment

import (
	"fmt"
	"glox/lerr"
	"glox/token"
)

type Environment struct {
	enclosing *Environment
	values    map[string]interface{}
}

func NewEnvironment(env *Environment) *Environment {
	return &Environment{
		enclosing: env,
		values:    make(map[string]interface{}),
	}
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Get(name token.Token) (interface{}, error) {
	val, ok := e.values[name.GetLexeme()]
	if ok {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}

	return nil, lerr.NewRuntimeErr(
		name, fmt.Sprintf("Undefined variable '%s'.", name.GetLexeme()))
}

func (e *Environment) Assign(name token.Token, value interface{}) error {
	if _, ok := e.values[name.GetLexeme()]; ok {
		e.values[name.GetLexeme()] = value
		return nil
	}

	if e.enclosing != nil {
		e.enclosing.Assign(name, value)
	}

	return lerr.NewRuntimeErr(name,
		fmt.Sprintf("Undefined variable %s", name.GetLexeme()))
}
