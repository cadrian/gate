package errors

import (
	"errors"
	"fmt"
	"runtime"
)

type StackError struct {
	Nested error
	StackTrace string
}

func (self StackError) Error() string {
	return fmt.Sprintf("%s\n%s", self.Nested.Error(), self.StackTrace)
}

func newerror(err error) error {
	const size = 4096
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	return StackError{
		Nested: err,
		StackTrace: string(buf),
	}
}

func New(message string) error {
	return newerror(errors.New(message))
}

func Newf(format string, args... interface{}) error {
	return newerror(errors.New(fmt.Sprintf(format, args...)))
}

func Decorated(err error) error {
	return newerror(err)
}
