package runtime

import (
	"fmt"
	"reflect"
	"strings"
)

// Context holds the data context for template execution.
// It provides variable lookup and scope management.
type Context struct {
	data   interface{}              // Root data
	scopes []map[string]interface{} // Stack of local scopes
}

// NewContext creates a new execution context with the given data.
func NewContext(data interface{}) *Context {
	return &Context{
		data:   data,
		scopes: make([]map[string]interface{}, 0),
	}
}

// Get retrieves a value from the context using dot notation.
// Examples: ".", ".Name", ".User.Email", ".Items[0]"
func (c *Context) Get(path []string) (interface{}, error) {
	if len(path) == 0 {
		return nil, fmt.Errorf("empty path")
	}

	// Start with root data if path begins with "."
	var current interface{}
	pathIndex := 0

	if path[0] == "." {
		// Check if "." has been redefined in a scope (e.g., in a range loop)
		found := false
		for i := len(c.scopes) - 1; i >= 0; i-- {
			if val, ok := c.scopes[i]["."]; ok {
				current = val
				found = true
				break
			}
		}
		if !found {
			current = c.data
		}
		pathIndex = 1
	} else {
		// Try to find in scopes first
		varName := path[0]
		found := false
		for i := len(c.scopes) - 1; i >= 0; i-- {
			if val, ok := c.scopes[i][varName]; ok {
				current = val
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("variable not found: %s", varName)
		}
		pathIndex = 1
	}

	// If just ".", return root data
	if pathIndex >= len(path) {
		return current, nil
	}

	// Traverse the path
	for i := pathIndex; i < len(path); i++ {
		current = c.getField(current, path[i])
		if current == nil {
			return nil, fmt.Errorf("nil value at path: %s", strings.Join(path[:i+1], "."))
		}
	}

	return current, nil
}

// getField retrieves a field from a struct, map, or slice using reflection.
func (c *Context) getField(obj interface{}, field string) interface{} {
	if obj == nil {
		return nil
	}

	val := reflect.ValueOf(obj)

	// Dereference pointers
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Struct:
		// Access struct field by name
		fieldVal := val.FieldByName(field)
		if !fieldVal.IsValid() {
			return nil
		}
		return fieldVal.Interface()

	case reflect.Map:
		// Access map by key
		mapKey := reflect.ValueOf(field)
		fieldVal := val.MapIndex(mapKey)
		if !fieldVal.IsValid() {
			return nil
		}
		return fieldVal.Interface()

	default:
		return nil
	}
}

// Set sets a value in the current scope.
func (c *Context) Set(name string, value interface{}) {
	if len(c.scopes) == 0 {
		c.PushScope()
	}
	c.scopes[len(c.scopes)-1][name] = value
}

// PushScope creates a new variable scope.
func (c *Context) PushScope() {
	c.scopes = append(c.scopes, make(map[string]interface{}))
}

// PopScope removes the current variable scope.
func (c *Context) PopScope() {
	if len(c.scopes) > 0 {
		c.scopes = c.scopes[:len(c.scopes)-1]
	}
}

// GetIndex retrieves a value from a slice or map by index/key.
func (c *Context) GetIndex(obj, index interface{}) (interface{}, error) {
	if obj == nil {
		return nil, fmt.Errorf("cannot index nil value")
	}

	val := reflect.ValueOf(obj)

	// Dereference pointers
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, fmt.Errorf("cannot index nil pointer")
		}
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		// Index must be an integer
		idx, ok := index.(int)
		if !ok {
			return nil, fmt.Errorf("array/slice index must be integer, got %T", index)
		}
		if idx < 0 || idx >= val.Len() {
			return nil, fmt.Errorf("index out of bounds: %d (len=%d)", idx, val.Len())
		}
		return val.Index(idx).Interface(), nil

	case reflect.Map:
		// Use the index as map key
		mapKey := reflect.ValueOf(index)
		mapVal := val.MapIndex(mapKey)
		if !mapVal.IsValid() {
			return nil, nil // Key not found, return nil
		}
		return mapVal.Interface(), nil

	default:
		return nil, fmt.Errorf("cannot index type %s", val.Kind())
	}
}

// IsTruthy evaluates whether a value is considered true in template context.
func IsTruthy(val interface{}) bool {
	if val == nil {
		return false
	}

	v := reflect.ValueOf(val)

	switch v.Kind() {
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return v.Float() != 0
	case reflect.String:
		return v.String() != ""
	case reflect.Slice, reflect.Map, reflect.Array:
		return v.Len() > 0
	case reflect.Ptr, reflect.Interface:
		return !v.IsNil()
	default:
		return true
	}
}

// ToSlice converts a value to a slice for iteration.
// Returns the slice and whether the conversion was successful.
func ToSlice(val interface{}) ([]interface{}, bool) {
	if val == nil {
		return nil, false
	}

	v := reflect.ValueOf(val)

	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, false
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result, true
	default:
		return nil, false
	}
}

// ToMap converts a value to a map for iteration.
// Returns the map keys and values, and whether the conversion was successful.
func ToMap(val interface{}) (keys, values []interface{}, ok bool) {
	if val == nil {
		return nil, nil, false
	}

	v := reflect.ValueOf(val)

	// Dereference pointers
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil, nil, false
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Map {
		return nil, nil, false
	}

	mapKeys := v.MapKeys()
	keys = make([]interface{}, len(mapKeys))
	values = make([]interface{}, len(mapKeys))

	for i, k := range mapKeys {
		keys[i] = k.Interface()
		values[i] = v.MapIndex(k).Interface()
	}

	return keys, values, true
}
