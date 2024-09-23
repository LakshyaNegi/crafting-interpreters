package lerr

import "fmt"

type syntaxErr struct {
	line  int
	msg   string
	where string
}

func (e syntaxErr) Error() string {
	return fmt.Sprintf("[line %d] Error %s: %s ", e.line, e.where, e.msg)
}

func NewSyntaxErr(line int, where string, msg string) error {
	return &syntaxErr{line: line, where: where, msg: msg}
}
