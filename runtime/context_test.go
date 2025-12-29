package runtime

import (
	"testing"
)

func TestContext_Get(t *testing.T) {
	data := map[string]interface{}{
		"Name": "Alice",
		"User": map[string]interface{}{
			"Email": "alice@example.com",
		},
	}

	ctx := NewContext(data)

	tests := []struct {
		name     string
		path     []string
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "root",
			path:     []string{"."},
			expected: data,
			wantErr:  false,
		},
		{
			name:     "simple field",
			path:     []string{".", "Name"},
			expected: "Alice",
			wantErr:  false,
		},
		{
			name:     "nested field",
			path:     []string{".", "User", "Email"},
			expected: "alice@example.com",
			wantErr:  false,
		},
		{
			name:    "empty path",
			path:    []string{},
			wantErr: true,
		},
		{
			name:    "non-existent field",
			path:    []string{".", "NonExistent"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := ctx.Get(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && val == nil {
				t.Errorf("Get() returned nil for valid path")
			}
		})
	}
}

func TestContext_Scopes(t *testing.T) {
	ctx := NewContext(map[string]interface{}{"root": "value"})

	// Initially no scopes
	if len(ctx.scopes) != 0 {
		t.Error("expected no scopes initially")
	}

	// Push a scope and set a variable
	ctx.PushScope()
	ctx.Set("x", 42)

	val, err := ctx.Get([]string{"x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}

	// Push another scope
	ctx.PushScope()
	ctx.Set("y", "hello")

	val, err = ctx.Get([]string{"y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "hello" {
		t.Errorf("expected 'hello', got %v", val)
	}

	// x should still be accessible
	val, err = ctx.Get([]string{"x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != 42 {
		t.Errorf("expected 42, got %v", val)
	}

	// Pop scope
	ctx.PopScope()

	// y should no longer be accessible
	_, err = ctx.Get([]string{"y"})
	if err == nil {
		t.Error("expected error for popped variable")
	}

	// x should still be there
	_, err = ctx.Get([]string{"x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestContext_GetIndex(t *testing.T) {
	ctx := NewContext(nil)

	tests := []struct {
		name    string
		obj     interface{}
		index   interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:  "slice with int index",
			obj:   []string{"a", "b", "c"},
			index: 1,
			want:  "b",
		},
		{
			name:  "map with string key",
			obj:   map[string]int{"x": 10, "y": 20},
			index: "x",
			want:  10,
		},
		{
			name:    "nil object",
			obj:     nil,
			index:   0,
			wantErr: true,
		},
		{
			name:    "index out of bounds",
			obj:     []int{1, 2, 3},
			index:   5,
			wantErr: true,
		},
		{
			name:    "wrong index type for slice",
			obj:     []int{1, 2, 3},
			index:   "not an int",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctx.GetIndex(tt.obj, tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("GetIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToSlice(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantLen int
		wantOK  bool
	}{
		{
			name:    "slice",
			val:     []int{1, 2, 3},
			wantLen: 3,
			wantOK:  true,
		},
		{
			name:    "array",
			val:     [3]string{"a", "b", "c"},
			wantLen: 3,
			wantOK:  true,
		},
		{
			name:   "nil",
			val:    nil,
			wantOK: false,
		},
		{
			name:   "map",
			val:    map[string]int{"x": 1},
			wantOK: false,
		},
		{
			name:   "string",
			val:    "hello",
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ToSlice(tt.val)
			if ok != tt.wantOK {
				t.Errorf("ToSlice() ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && len(got) != tt.wantLen {
				t.Errorf("ToSlice() len = %v, want %v", len(got), tt.wantLen)
			}
		})
	}
}

func TestToMap(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantLen int
		wantOK  bool
	}{
		{
			name:    "map",
			val:     map[string]int{"x": 1, "y": 2},
			wantLen: 2,
			wantOK:  true,
		},
		{
			name:   "nil",
			val:    nil,
			wantOK: false,
		},
		{
			name:   "slice",
			val:    []int{1, 2, 3},
			wantOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, vals, ok := ToMap(tt.val)
			if ok != tt.wantOK {
				t.Errorf("ToMap() ok = %v, want %v", ok, tt.wantOK)
			}
			if ok && (len(keys) != tt.wantLen || len(vals) != tt.wantLen) {
				t.Errorf("ToMap() len = %v, want %v", len(keys), tt.wantLen)
			}
		})
	}
}

func TestContext_GetField_Struct(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	ctx := NewContext(nil)
	person := Person{Name: "Alice", Age: 30}

	name := ctx.getField(person, "Name")
	if name != "Alice" {
		t.Errorf("expected 'Alice', got %v", name)
	}

	age := ctx.getField(person, "Age")
	if age != 30 {
		t.Errorf("expected 30, got %v", age)
	}

	// Non-existent field
	invalid := ctx.getField(person, "NonExistent")
	if invalid != nil {
		t.Errorf("expected nil for non-existent field, got %v", invalid)
	}
}

func TestContext_GetField_Nil(t *testing.T) {
	ctx := NewContext(nil)
	result := ctx.getField(nil, "field")
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestContext_ScopeDotOverride(t *testing.T) {
	ctx := NewContext(map[string]interface{}{"original": "data"})

	// Initially, "." refers to root data
	val, _ := ctx.Get([]string{".", "original"})
	if val != "data" {
		t.Errorf("expected 'data', got %v", val)
	}

	// In a loop, "." gets redefined
	ctx.PushScope()
	ctx.Set(".", "loop item")

	val, _ = ctx.Get([]string{"."})
	if val != "loop item" {
		t.Errorf("expected 'loop item', got %v", val)
	}

	ctx.PopScope()

	// After popping, "." should refer to root again
	val, _ = ctx.Get([]string{".", "original"})
	if val != "data" {
		t.Errorf("expected 'data', got %v", val)
	}
}
