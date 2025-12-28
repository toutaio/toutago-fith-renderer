// Package fith provides error types for the Fíth template engine.
package fith

import (
	"fmt"
)

// Error represents a template error with context information.
type Error struct {
	Type    ErrorType
	Message string
	Slug    string
	Line    int
	Column  int
	Cause   error
}

// ErrorType categorizes template errors.
type ErrorType int

const (
	// ErrorTypeUnknown represents an unknown error type.
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeTemplate indicates a template syntax error.
	ErrorTypeTemplate
	// ErrorTypeCompilation indicates a compilation error.
	ErrorTypeCompilation
	// ErrorTypeRuntime indicates a runtime execution error.
	ErrorTypeRuntime
	// ErrorTypeLoader indicates a template loading error.
	ErrorTypeLoader
	// ErrorTypeFunction indicates a function-related error.
	ErrorTypeFunction
)

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Slug != "" && e.Line > 0 {
		return fmt.Sprintf("%s [%s:%d:%d]: %s", e.Type, e.Slug, e.Line, e.Column, e.Message)
	}
	if e.Slug != "" {
		return fmt.Sprintf("%s [%s]: %s", e.Type, e.Slug, e.Message)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause error.
func (e *Error) Unwrap() error {
	return e.Cause
}

// String returns a string representation of the error type.
func (t ErrorType) String() string {
	switch t {
	case ErrorTypeTemplate:
		return "TemplateError"
	case ErrorTypeCompilation:
		return "CompilationError"
	case ErrorTypeRuntime:
		return "RuntimeError"
	case ErrorTypeLoader:
		return "LoaderError"
	case ErrorTypeFunction:
		return "FunctionError"
	default:
		return "UnknownError"
	}
}

// NewError creates a new Error with the given type and message.
func NewError(errType ErrorType, message string) *Error {
	return &Error{
		Type:    errType,
		Message: message,
	}
}

// NewErrorWithLocation creates a new Error with location information.
func NewErrorWithLocation(errType ErrorType, message, slug string, line, column int) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Slug:    slug,
		Line:    line,
		Column:  column,
	}
}

// WrapError wraps an existing error with additional context.
func WrapError(errType ErrorType, message string, cause error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Cause:   cause,
	}
}
