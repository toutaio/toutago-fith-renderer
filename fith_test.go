package fith

import (
	"embed"
	"os"
	"path/filepath"
	"testing"
)

//go:embed testdata/*.html
var testFS embed.FS

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with template dir",
			config: Config{
				TemplateDir: "templates",
			},
			wantErr: false,
		},
		{
			name: "valid config with embed fs",
			config: Config{
				TemplateFS: testFS,
			},
			wantErr: false,
		},
		{
			name: "empty config gets defaults",
			config: Config{
				LeftDelimiter:  "{{",
				RightDelimiter: "}}",
			},
			wantErr: false,
		},
		{
			name: "invalid config - same delimiters",
			config: Config{
				TemplateDir:    "templates",
				LeftDelimiter:  "{{",
				RightDelimiter: "{{",
			},
			wantErr: true,
		},
		{
			name: "zero max include depth gets defaulted",
			config: Config{
				TemplateDir: "templates",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewWithDefaults(t *testing.T) {
	engine, err := NewWithDefaults()
	if err != nil {
		t.Fatalf("NewWithDefaults() error = %v", err)
	}

	if engine.config.TemplateDir != "templates" {
		t.Errorf("expected default TemplateDir = 'templates', got %q", engine.config.TemplateDir)
	}

	if !engine.config.CacheEnabled {
		t.Error("expected CacheEnabled = true by default")
	}
}

func TestNewWithFS(t *testing.T) {
	engine, err := NewWithFS(testFS)
	if err != nil {
		t.Fatalf("NewWithFS() error = %v", err)
	}

	if engine.config.TemplateFS == nil {
		t.Error("expected TemplateFS to be set")
	}
}

func TestNewWithDir(t *testing.T) {
	tmpDir := t.TempDir()

	engine, err := NewWithDir(tmpDir)
	if err != nil {
		t.Fatalf("NewWithDir() error = %v", err)
	}

	if engine.config.TemplateDir != tmpDir {
		t.Errorf("expected TemplateDir = %q, got %q", tmpDir, engine.config.TemplateDir)
	}
}

func TestRenderString(t *testing.T) {
	engine, err := NewWithDefaults()
	if err != nil {
		t.Fatalf("NewWithDefaults() error = %v", err)
	}

	tests := []struct {
		name     string
		template string
		data     interface{}
		want     string
		wantErr  bool
	}{
		{
			name:     "simple variable",
			template: "Hello {{.Name}}!",
			data: map[string]interface{}{
				"Name": "World",
			},
			want:    "Hello World!",
			wantErr: false,
		},
		{
			name:     "nested data",
			template: "User: {{.User.Name}}",
			data: map[string]interface{}{
				"User": map[string]interface{}{
					"Name": "Alice",
				},
			},
			want:    "User: Alice",
			wantErr: false,
		},
		{
			name:     "if statement",
			template: "{{if .Show}}Visible{{end}}",
			data: map[string]interface{}{
				"Show": true,
			},
			want:    "Visible",
			wantErr: false,
		},
		{
			name:     "if else statement",
			template: "{{if .Show}}Yes{{else}}No{{end}}",
			data: map[string]interface{}{
				"Show": false,
			},
			want:    "No",
			wantErr: false,
		},
		{
			name:     "range loop",
			template: "{{range .Items}}{{.}} {{end}}",
			data: map[string]interface{}{
				"Items": []string{"a", "b", "c"},
			},
			want:    "a b c ",
			wantErr: false,
		},
		{
			name:     "function call",
			template: "{{upper .Name}}",
			data: map[string]interface{}{
				"Name": "hello",
			},
			want:    "HELLO",
			wantErr: false,
		},
		{
			name:     "filter pipeline",
			template: "{{.Name | upper | trim}}",
			data: map[string]interface{}{
				"Name": " world ",
			},
			want:    "WORLD",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.RenderString(tt.template, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("RenderString() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRender(t *testing.T) {
	// Create temp directory with test templates
	tmpDir := t.TempDir()

	// Create test templates
	templates := map[string]string{
		"simple.html":  "Hello {{.Name}}!",
		"nested.html":  "User: {{.User.Name}}",
		"with-if.html": "{{if .Show}}Visible{{else}}Hidden{{end}}",
	}

	for name, content := range templates {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test template: %v", err)
		}
	}

	engine, err := NewWithDir(tmpDir)
	if err != nil {
		t.Fatalf("NewWithDir() error = %v", err)
	}

	tests := []struct {
		name    string
		slug    string
		data    interface{}
		want    string
		wantErr bool
	}{
		{
			name: "simple template",
			slug: "simple",
			data: map[string]interface{}{
				"Name": "World",
			},
			want:    "Hello World!",
			wantErr: false,
		},
		{
			name: "nested data",
			slug: "nested",
			data: map[string]interface{}{
				"User": map[string]interface{}{
					"Name": "Alice",
				},
			},
			want:    "User: Alice",
			wantErr: false,
		},
		{
			name: "conditional - true",
			slug: "with-if",
			data: map[string]interface{}{
				"Show": true,
			},
			want:    "Visible",
			wantErr: false,
		},
		{
			name: "conditional - false",
			slug: "with-if",
			data: map[string]interface{}{
				"Show": false,
			},
			want:    "Hidden",
			wantErr: false,
		},
		{
			name:    "non-existent template",
			slug:    "nonexistent",
			data:    map[string]interface{}{},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := engine.Render(tt.slug, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRegisterFunction(t *testing.T) {
	engine, err := NewWithDefaults()
	if err != nil {
		t.Fatalf("NewWithDefaults() error = %v", err)
	}

	// Register custom function
	engine.RegisterFunction("double", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, nil
		}
		if n, ok := args[0].(int); ok {
			return n * 2, nil
		}
		return nil, nil
	})

	template := "{{double .Value}}"
	data := map[string]interface{}{
		"Value": 5,
	}

	got, err := engine.RenderString(template, data)
	if err != nil {
		t.Errorf("RenderString() error = %v", err)
		return
	}

	want := "10"
	if got != want {
		t.Errorf("RenderString() = %q, want %q", got, want)
	}
}

func TestExists(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test template
	path := filepath.Join(tmpDir, "test.html")
	if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	engine, err := NewWithDir(tmpDir)
	if err != nil {
		t.Fatalf("NewWithDir() error = %v", err)
	}

	if !engine.Exists("test") {
		t.Error("Exists() = false, want true for existing template")
	}

	if engine.Exists("nonexistent") {
		t.Error("Exists() = true, want false for non-existent template")
	}
}

func TestClearCache(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test template
	path := filepath.Join(tmpDir, "test.html")
	if err := os.WriteFile(path, []byte("Hello {{.Name}}!"), 0644); err != nil {
		t.Fatalf("failed to create test template: %v", err)
	}

	engine, err := NewWithDir(tmpDir)
	if err != nil {
		t.Fatalf("NewWithDir() error = %v", err)
	}

	data := map[string]interface{}{"Name": "World"}

	// First render - should cache
	_, err = engine.Render("test", data)
	if err != nil {
		t.Fatalf("first Render() error = %v", err)
	}

	// Clear cache
	engine.ClearCache()

	// Second render - should work after cache clear
	_, err = engine.Render("test", data)
	if err != nil {
		t.Fatalf("second Render() error = %v", err)
	}
}

func TestConfig(t *testing.T) {
	config := Config{
		TemplateDir:     "templates",
		LeftDelimiter:   "{{",
		RightDelimiter:  "}}",
		CacheEnabled:    true,
		MaxIncludeDepth: 50,
	}

	engine, err := New(config)
	if err != nil {
		t.Fatalf("New() error = %v", err)
	}

	got := engine.Config()

	if got.TemplateDir != config.TemplateDir {
		t.Errorf("Config().TemplateDir = %q, want %q", got.TemplateDir, config.TemplateDir)
	}

	if got.MaxIncludeDepth != config.MaxIncludeDepth {
		t.Errorf("Config().MaxIncludeDepth = %d, want %d", got.MaxIncludeDepth, config.MaxIncludeDepth)
	}
}

func BenchmarkRenderString(b *testing.B) {
	engine, err := NewWithDefaults()
	if err != nil {
		b.Fatalf("NewWithDefaults() error = %v", err)
	}

	template := "Hello {{.Name}}! You have {{.Count}} messages."
	data := map[string]interface{}{
		"Name":  "World",
		"Count": 5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.RenderString(template, data)
		if err != nil {
			b.Fatalf("RenderString() error = %v", err)
		}
	}
}

func BenchmarkRenderWithCache(b *testing.B) {
	tmpDir := b.TempDir()

	path := filepath.Join(tmpDir, "bench.html")
	content := "Hello {{.Name}}! You have {{.Count}} messages."
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		b.Fatalf("failed to create test template: %v", err)
	}

	engine, err := NewWithDir(tmpDir)
	if err != nil {
		b.Fatalf("NewWithDir() error = %v", err)
	}

	data := map[string]interface{}{
		"Name":  "World",
		"Count": 5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := engine.Render("bench", data)
		if err != nil {
			b.Fatalf("Render() error = %v", err)
		}
	}
}
