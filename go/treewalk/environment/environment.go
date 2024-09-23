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

func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	env := e.ancestor(distance)
	val, ok := env.values[name]
	if ok {
		return val, nil
	}

	return nil, lerr.NewRuntimeErr(
		nil, fmt.Sprintf("Undefined variable '%s'.", name))
}

func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.enclosing
	}

	return env
}

func (e *Environment) AssignAt(distance int, name token.Token, value interface{}) error {
	env := e.ancestor(distance)
	if _, ok := env.values[name.GetLexeme()]; ok {
		env.values[name.GetLexeme()] = value
		return nil
	}

	return lerr.NewRuntimeErr(
		name, fmt.Sprintf("Undefined variable '%s'.", name.GetLexeme()))
}
