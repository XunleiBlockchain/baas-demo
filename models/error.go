package models

import "fmt"

type ErrWrapper struct {
	code   int
	errmsg string
}

func (e *ErrWrapper) Code() int {
	return e.code
}

func (e *ErrWrapper) Error() string {
	return e.errmsg
}

func (e *ErrWrapper) Join(err error) *ErrWrapper {
	return &ErrWrapper{
		code: e.code,
		errmsg: fmt.Sprintf("%s %v", e.errmsg, err),
	}
}

var (
	ErrParams = &ErrWrapper{
		code:   -10000,
		errmsg: "params error.",
	}
	ErrServer = &ErrWrapper{
		code:   -10001,
		errmsg: "server error.",
	}
)
