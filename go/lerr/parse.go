package lerr

type ParseErr struct {
}

func (e ParseErr) Error() string {
	return "Parse error"
}

func NewParseErr() error {
	return &ParseErr{}
}
