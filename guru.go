// Package guru provides Go errors with a Guru Meditation Code.
package guru

import (
	"errors"
	"fmt"
)

// coder is the main interface to errors in this package.
type coder interface {
	Code() int
}

type withCode struct {
	error
	code int
}

func (e *withCode) Unwrap() error                { return e.error }
func (e *withCode) Code() int                    { return e.code }
func (e withCode) Format(s fmt.State, verb rune) { fmt.Fprintf(s, "error %v: %v", e.code, e.error) }

type wrapped struct {
	msg  string
	code int
	error
}

func (e *wrapped) Error() string { return e.msg }
func (e *wrapped) Unwrap() error { return e.error }
func (e *wrapped) Code() int     { return e.code }
func (e wrapped) Format(s fmt.State, verb rune) {
	fmt.Fprintf(s, "error %v: %v", e.code, e.error)
	if e.msg != "" {
		fmt.Fprintf(s, ": %v", e.msg)
	}
}

// New returns a new error message with an error code.
func New(code int, msg string) error {
	return &withCode{
		error: errors.New(msg),
		code:  code,
	}
}

// Errorf returns a new error message with an error code.
func Errorf(code int, format string, args ...interface{}) error {
	return &withCode{
		error: fmt.Errorf(format, args...),
		code:  code,
	}
}

// WithCode wraps an existing error with the provided error code. It will return
// nil if err is nil.
func WithCode(code int, err error) error {
	if err == nil {
		return nil
	}
	return &withCode{
		error: err,
		code:  code,
	}
}

// Wrap returns an error annotating err with an error code, and the supplied
// message. It will return nil if err is nil.
func Wrap(code int, err error, msg string) error {
	if err == nil {
		return nil
	}
	return &wrapped{
		msg:   msg,
		code:  code,
		error: err,
	}
}

// Wrapf returns an error annotating err with an error code, and the format
// specifier. It will return nil if err is nil.
func Wrapf(code int, err error, msg string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wrapped{
		msg:   fmt.Sprintf(msg, args...),
		code:  code,
		error: err,
	}
}

// Code extracts the highest-level error code from the error or the errors it
// wraps. It will return 0 if the error does not implement the coder interface.
func Code(err error) int {
	for {
		if sc, ok := err.(coder); ok {
			return sc.Code()
		}
		err := errors.Unwrap(err)
		if err == nil {
			break
		}
	}
	return 0
}
