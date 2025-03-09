package models

import (
	"errors"
	"strings"
)

// ESError represents an error with a user-friendly message and detailed error chain.
type ESError struct {
	Message      string `json:"message"` // Readable message for the end user
	err          error  `json:"-"`       // The actual error (not serialized)
	DisplayError string `json:"error"`   // The generated error message for API response
}

// NewESError creates a new ESError with the given message and error, unwrapping the error chain.
func NewESError(message string, err error) ESError {
	esError := ESError{
		Message: message,
		err:     err,
	}
	esError.unwrapError()
	return esError
}

// unwrapError populates DisplayError with the unwrapped error chain, if any.
func (e *ESError) unwrapError() {
	if e.err == nil {
		e.DisplayError = ""
		return
	}

	var errChain []string
	currentErr := e.err
	for currentErr != nil {
		errChain = append(errChain, currentErr.Error())
		currentErr = errors.Unwrap(currentErr)
	}
	e.DisplayError = strings.Join(errChain, " -> ")
}
