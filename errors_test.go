package fith

import (
	"errors"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(ErrorTypeTemplate, "test error")

	if err.Type != ErrorTypeTemplate {
		t.Errorf("expected Type = ErrorTypeTemplate, got %v", err.Type)
	}

	if err.Message != "test error" {
		t.Errorf("expected Message = 'test error', got %q", err.Message)
	}
}

func TestNewErrorWithLocation(t *testing.T) {
	err := NewErrorWithLocation(ErrorTypeRuntime, "runtime error", "test.html", 10, 5)

	if err.Type != ErrorTypeRuntime {
		t.Errorf("expected Type = ErrorTypeRuntime, got %v", err.Type)
	}

	if err.Message != "runtime error" {
		t.Errorf("expected Message = 'runtime error', got %q", err.Message)
	}

	if err.Slug != "test.html" {
		t.Errorf("expected Slug = 'test.html', got %q", err.Slug)
	}

	if err.Line != 10 {
		t.Errorf("expected Line = 10, got %d", err.Line)
	}

	if err.Column != 5 {
		t.Errorf("expected Column = 5, got %d", err.Column)
	}
}

func TestWrapError(t *testing.T) {
	cause := errors.New("underlying error")
	err := WrapError(ErrorTypeCompilation, "compilation failed", cause)

	if err.Type != ErrorTypeCompilation {
		t.Errorf("expected Type = ErrorTypeCompilation, got %v", err.Type)
	}

	if err.Message != "compilation failed" {
		t.Errorf("expected Message = 'compilation failed', got %q", err.Message)
	}

	if err.Cause != cause {
		t.Errorf("expected Cause to be set")
	}
}

func TestErrorError(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "error with location",
			err: &Error{
				Type:    ErrorTypeTemplate,
				Message: "syntax error",
				Slug:    "test.html",
				Line:    5,
				Column:  10,
			},
			want: "TemplateError [test.html:5:10]: syntax error",
		},
		{
			name: "error with slug only",
			err: &Error{
				Type:    ErrorTypeLoader,
				Message: "not found",
				Slug:    "missing.html",
			},
			want: "LoaderError [missing.html]: not found",
		},
		{
			name: "error without location",
			err: &Error{
				Type:    ErrorTypeFunction,
				Message: "invalid argument",
			},
			want: "FunctionError: invalid argument",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("cause error")
	err := &Error{
		Type:    ErrorTypeRuntime,
		Message: "wrapped",
		Cause:   cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, cause)
	}
}

func TestErrorTypeString(t *testing.T) {
	tests := []struct {
		errType ErrorType
		want    string
	}{
		{ErrorTypeTemplate, "TemplateError"},
		{ErrorTypeCompilation, "CompilationError"},
		{ErrorTypeRuntime, "RuntimeError"},
		{ErrorTypeLoader, "LoaderError"},
		{ErrorTypeFunction, "FunctionError"},
		{ErrorTypeUnknown, "UnknownError"},
		{ErrorType(999), "UnknownError"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.errType.String()
			if got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}
