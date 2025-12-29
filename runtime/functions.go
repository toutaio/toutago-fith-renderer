package runtime

import (
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"
)

// Function represents a template function.
type Function func(args ...interface{}) (interface{}, error)

// FunctionRegistry manages available template functions.
type FunctionRegistry struct {
	funcs map[string]Function
}

// NewFunctionRegistry creates a new function registry with built-in functions.
func NewFunctionRegistry() *FunctionRegistry {
	registry := &FunctionRegistry{
		funcs: make(map[string]Function),
	}
	registry.registerBuiltins()
	return registry
}

// Register adds or replaces a function in the registry.
func (r *FunctionRegistry) Register(name string, fn Function) {
	r.funcs[name] = fn
}

// Get retrieves a function by name.
func (r *FunctionRegistry) Get(name string) (Function, bool) {
	fn, ok := r.funcs[name]
	return fn, ok
}

// Call executes a function by name with the given arguments.
func (r *FunctionRegistry) Call(name string, args ...interface{}) (interface{}, error) {
	fn, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("unknown function: %s", name)
	}
	return fn(args...)
}

// AllFunctions returns all registered functions.
func (r *FunctionRegistry) AllFunctions() map[string]Function {
	return r.funcs
}

// registerBuiltins adds all built-in functions to the registry.
func (r *FunctionRegistry) registerBuiltins() {
	// String functions
	r.Register("upper", fnUpper)
	r.Register("lower", fnLower)
	r.Register("title", fnTitle)
	r.Register("trim", fnTrim)
	r.Register("trimPrefix", fnTrimPrefix)
	r.Register("trimSuffix", fnTrimSuffix)
	r.Register("truncate", fnTruncate)
	r.Register("replace", fnReplace)

	// Array functions
	r.Register("join", fnJoin)
	r.Register("len", fnLen)
	r.Register("first", fnFirst)
	r.Register("last", fnLast)

	// Logic functions
	r.Register("default", fnDefault)

	// Encoding functions
	r.Register("urlEncode", fnURLEncode)
	r.Register("htmlEscape", fnHTMLEscape)

	// Date functions
	r.Register("date", fnDate)
}

// ============================================================================
// String Functions
// ============================================================================

func fnUpper(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("upper: expected 1 argument, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("upper: argument must be a string")
	}
	return strings.ToUpper(s), nil
}

func fnLower(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("lower: expected 1 argument, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("lower: argument must be a string")
	}
	return strings.ToLower(s), nil
}

func fnTitle(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("title: expected 1 argument, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("title: argument must be a string")
	}
	// nolint:staticcheck // Using deprecated strings.Title for simplicity
	return strings.Title(s), nil //nolint:staticcheck
}

func fnTrim(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("trim: expected 1 argument, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("trim: argument must be a string")
	}
	return strings.TrimSpace(s), nil
}

func fnTrimPrefix(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("trimPrefix: expected 2 arguments, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("trimPrefix: first argument must be a string")
	}
	prefix, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("trimPrefix: second argument must be a string")
	}
	return strings.TrimPrefix(s, prefix), nil
}

func fnTrimSuffix(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("trimSuffix: expected 2 arguments, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("trimSuffix: first argument must be a string")
	}
	suffix, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("trimSuffix: second argument must be a string")
	}
	return strings.TrimSuffix(s, suffix), nil
}

func fnTruncate(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("truncate: expected 2 arguments, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("truncate: first argument must be a string")
	}

	var maxLen int
	switch v := args[1].(type) {
	case int:
		maxLen = v
	case int64:
		maxLen = int(v)
	case float64:
		maxLen = int(v)
	default:
		return nil, fmt.Errorf("truncate: second argument must be a number")
	}

	if maxLen < 0 {
		return nil, fmt.Errorf("truncate: length must be non-negative")
	}

	runes := []rune(s)
	if len(runes) <= maxLen {
		return s, nil
	}
	return string(runes[:maxLen]) + "...", nil
}

func fnReplace(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		return nil, fmt.Errorf("replace: expected 3 arguments, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("replace: first argument must be a string")
	}
	old, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("replace: second argument must be a string")
	}
	new, ok := args[2].(string)
	if !ok {
		return nil, fmt.Errorf("replace: third argument must be a string")
	}
	return strings.ReplaceAll(s, old, new), nil
}

// ============================================================================
// Array Functions
// ============================================================================

func fnJoin(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("join: expected 2 arguments, got %d", len(args))
	}

	arr, ok := ToSlice(args[0])
	if !ok {
		return nil, fmt.Errorf("join: first argument must be an array or slice")
	}

	sep, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("join: second argument must be a string")
	}

	parts := make([]string, len(arr))
	for i, v := range arr {
		parts[i] = fmt.Sprint(v)
	}

	return strings.Join(parts, sep), nil
}

func fnLen(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("len: expected 1 argument, got %d", len(args))
	}

	if arr, ok := ToSlice(args[0]); ok {
		return len(arr), nil
	}

	if s, ok := args[0].(string); ok {
		return len([]rune(s)), nil
	}

	if keys, _, ok := ToMap(args[0]); ok {
		return len(keys), nil
	}

	return nil, fmt.Errorf("len: argument must be a string, array, or map")
}

func fnFirst(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("first: expected 1 argument, got %d", len(args))
	}

	arr, ok := ToSlice(args[0])
	if !ok {
		return nil, fmt.Errorf("first: argument must be an array or slice")
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("first: array is empty")
	}

	return arr[0], nil
}

func fnLast(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("last: expected 1 argument, got %d", len(args))
	}

	arr, ok := ToSlice(args[0])
	if !ok {
		return nil, fmt.Errorf("last: argument must be an array or slice")
	}

	if len(arr) == 0 {
		return nil, fmt.Errorf("last: array is empty")
	}

	return arr[len(arr)-1], nil
}

// ============================================================================
// Logic Functions
// ============================================================================

func fnDefault(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("default: expected 2 arguments, got %d", len(args))
	}

	val := args[0]
	defaultVal := args[1]

	// Return default if value is nil, empty string, 0, or false
	if val == nil {
		return defaultVal, nil
	}

	switch v := val.(type) {
	case string:
		if v == "" {
			return defaultVal, nil
		}
	case int, int64:
		if v == 0 {
			return defaultVal, nil
		}
	case float64:
		if v == 0.0 {
			return defaultVal, nil
		}
	case bool:
		if !v {
			return defaultVal, nil
		}
	}

	// Check for empty slices
	if arr, ok := ToSlice(val); ok && len(arr) == 0 {
		return defaultVal, nil
	}

	return val, nil
}

// ============================================================================
// Encoding Functions
// ============================================================================

func fnURLEncode(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("urlEncode: expected 1 argument, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("urlEncode: argument must be a string")
	}
	return url.QueryEscape(s), nil
}

func fnHTMLEscape(args ...interface{}) (interface{}, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("htmlEscape: expected 1 argument, got %d", len(args))
	}
	s, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("htmlEscape: argument must be a string")
	}
	return html.EscapeString(s), nil
}

// ============================================================================
// Date Functions
// ============================================================================

func fnDate(args ...interface{}) (interface{}, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("date: expected 2 arguments, got %d", len(args))
	}

	format, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("date: first argument must be a format string")
	}

	var t time.Time
	switch v := args[1].(type) {
	case time.Time:
		t = v
	case string:
		// Try parsing common formats
		var err error
		for _, layout := range []string{time.RFC3339, "2006-01-02", "2006-01-02 15:04:05"} {
			t, err = time.Parse(layout, v)
			if err == nil {
				break
			}
		}
		if err != nil {
			return nil, fmt.Errorf("date: could not parse time string: %v", err)
		}
	default:
		return nil, fmt.Errorf("date: second argument must be a time.Time or string")
	}

	// Convert Go format to common formats
	format = convertDateFormat(format)

	return t.Format(format), nil
}

// convertDateFormat converts common date format strings to Go time format.
func convertDateFormat(format string) string {
	// Important: Replace longer patterns first to avoid partial replacements
	replacements := []struct {
		old string
		new string
	}{
		{"YYYY", "2006"},
		{"YY", "06"},
		{"MM", "01"},
		{"DD", "02"},
		{"hh", "15"},
		{"mm", "04"},
		{"ss", "05"},
	}

	result := format
	for _, r := range replacements {
		result = strings.ReplaceAll(result, r.old, r.new)
	}
	return result
}
