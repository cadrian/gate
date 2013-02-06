// This file is part of Gate.
// Copyright (C) 2012-2013 Cyril Adrian <cyril.adrian@gmail.com>
//
// Gate is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// Gate is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.	 See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Gate.  If not, see <http://www.gnu.org/licenses/>.

// Getting a stacktrace from Go errors.
package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// An error and the corresponding stacktrace
type StackError struct {
	Nested     error
	StackTrace string
}

func (self StackError) Error() string {
	return fmt.Sprintf("%s\n%s", self.Nested.Error(), self.StackTrace)
}

func newerror(err error) error {
	const size = 16384
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]

	// skip the frames inside this package (4 irrelevant lines)
	stack := strings.Split(string(buf), "\n")
	stacktrace := fmt.Sprintf("Traceback of %s\n%s", stack[0], strings.Join(stack[5:], "\n"))

	return StackError{
		Nested:     err,
		StackTrace: stacktrace,
	}
}

// Create a brand new error using the provided message
func New(message string) error {
	return newerror(errors.New(message))
}

// Create a brand new error using the provided format and arguments
func Newf(format string, args ...interface{}) error {
	return newerror(errors.New(fmt.Sprintf(format, args...)))
}

// Decorate an error
func Decorated(err error) error {
	switch err.(type) {
	case StackError:
		return err
	}
	return newerror(err)
}
