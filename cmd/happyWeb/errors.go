package main

import (
	"github.com/pkg/errors"
)

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code int
	Err  error
}

// Allows StatusError to satisfy the error interface.
func (se StatusError) Error() string {
	return se.Err.Error()
}

// Returns our HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

func newError(code int, msg string, err error) *StatusError {
	if err != nil {
		return &StatusError{Code: code, Err: errors.Errorf(msg+": %s", err)}
	} else {
		return &StatusError{Code: code, Err: errors.Errorf(msg)}
	}
}

func new500Error(msg string, err error) *StatusError {
	return newError(500, msg, err)
}

func new404Error(msg string, err error) *StatusError {
	return newError(404, msg, err)
}
