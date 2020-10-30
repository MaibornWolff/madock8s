package errors

import (
	"fmt"
)

type ErrInvalidInput struct {
	Method   string
	Argument string
}

func (e *ErrInvalidInput) Error() string {
	return fmt.Sprintf("invalid input %v in method %v", e.Argument, e.Method)
}

func Panic(e error) {
	if e != nil {
		panic(e)
	}
}
