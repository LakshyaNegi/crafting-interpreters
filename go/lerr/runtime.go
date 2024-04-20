package lerr

import "glox/token"

type RuntimeErr struct {
	token token.Token
	msg   string
}

func (e RuntimeErr) Error() string {
	return e.msg
}

func NewRuntimeErr(token token.Token, msg string) error {
	return &RuntimeErr{
		token: token,
		msg:   msg,
	}
}
