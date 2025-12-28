package fith

import (
	"testing"
	"testing/fstest"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.TemplateDir != "templates" {
		t.Errorf("expected TemplateDir = 'templates', got %q", config.TemplateDir)
	}

	if config.LeftDelimiter != "{{" {
		t.Errorf("expected LeftDelimiter = '{{', got %q", config.LeftDelimiter)
	}

	if config.RightDelimiter != "}}" {
		t.Errorf("expected RightDelimiter = '}}', got %q", config.RightDelimiter)
	}

	if !config.CacheEnabled {
		t.Error("expected CacheEnabled = true")
	}

	if config.AutoEscape {
		t.Error("expected AutoEscape = false")
	}

	if config.StrictMode {
		t.Error("expected StrictMode = false")
	}

	if config.MaxIncludeDepth != 100 {
		t.Errorf("expected MaxIncludeDepth = 100, got %d", config.MaxIncludeDepth)
	}

	expectedExtensions := []string{".html", ".tpl", ".txt"}
	if len(config.Extensions) != len(expectedExtensions) {
		t.Errorf("expected %d extensions, got %d", len(expectedExtensions), len(config.Extensions))
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "valid config with template dir",
			config: Config{
				TemplateDir:     "templates",
				LeftDelimiter:   "{{",
				RightDelimiter:  "}}",
				MaxIncludeDepth: 100,
			},
			wantErr: false,
		},
		{
			name: "valid config with template fs",
			config: Config{
				TemplateFS:      fstest.MapFS{},
				LeftDelimiter:   "{{",
				RightDelimiter:  "}}",
				MaxIncludeDepth: 100,
			},
			wantErr: false,
		},
		{
			name: "invalid - no template source",
			config: Config{
				LeftDelimiter:   "{{",
				RightDelimiter:  "}}",
				MaxIncludeDepth: 100,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty left delimiter",
			config: Config{
				TemplateDir:     "templates",
				LeftDelimiter:   "",
				RightDelimiter:  "}}",
				MaxIncludeDepth: 100,
			},
			wantErr: true,
		},
		{
			name: "invalid - empty right delimiter",
			config: Config{
				TemplateDir:     "templates",
				LeftDelimiter:   "{{",
				RightDelimiter:  "",
				MaxIncludeDepth: 100,
			},
			wantErr: true,
		},
		{
			name: "invalid - same delimiters",
			config: Config{
				TemplateDir:     "templates",
				LeftDelimiter:   "{{",
				RightDelimiter:  "{{",
				MaxIncludeDepth: 100,
			},
			wantErr: true,
		},
		{
			name: "invalid - zero max include depth",
			config: Config{
				TemplateDir:     "templates",
				LeftDelimiter:   "{{",
				RightDelimiter:  "}}",
				MaxIncludeDepth: 0,
			},
			wantErr: true,
		},
		{
			name: "auto-fills extensions",
			config: Config{
				TemplateDir:     "templates",
				LeftDelimiter:   "{{",
				RightDelimiter:  "}}",
				MaxIncludeDepth: 100,
				Extensions:      []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check that empty extensions are filled
			if tt.name == "auto-fills extensions" && err == nil {
				if len(tt.config.Extensions) == 0 {
					t.Error("expected Extensions to be filled with defaults")
				}
			}
		})
	}
}

func TestConfigApplyDefaults(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		check  func(*testing.T, *Config)
	}{
		{
			name:   "applies default template dir",
			config: Config{},
			check: func(t *testing.T, c *Config) {
				if c.TemplateDir != "templates" {
					t.Errorf("expected TemplateDir = 'templates', got %q", c.TemplateDir)
				}
			},
		},
		{
			name: "applies default extensions",
			config: Config{
				TemplateDir: "custom",
			},
			check: func(t *testing.T, c *Config) {
				expected := []string{".html", ".tpl", ".txt"}
				if len(c.Extensions) != len(expected) {
					t.Errorf("expected %d extensions, got %d", len(expected), len(c.Extensions))
				}
			},
		},
		{
			name: "applies default delimiters",
			config: Config{
				TemplateDir: "templates",
			},
			check: func(t *testing.T, c *Config) {
				if c.LeftDelimiter != "{{" {
					t.Errorf("expected LeftDelimiter = '{{', got %q", c.LeftDelimiter)
				}
				if c.RightDelimiter != "}}" {
					t.Errorf("expected RightDelimiter = '}}', got %q", c.RightDelimiter)
				}
			},
		},
		{
			name: "applies default max include depth",
			config: Config{
				TemplateDir: "templates",
			},
			check: func(t *testing.T, c *Config) {
				if c.MaxIncludeDepth != 100 {
					t.Errorf("expected MaxIncludeDepth = 100, got %d", c.MaxIncludeDepth)
				}
			},
		},
		{
			name: "preserves custom values",
			config: Config{
				TemplateDir:     "custom",
				Extensions:      []string{".tmpl"},
				LeftDelimiter:   "{%",
				RightDelimiter:  "%}",
				MaxIncludeDepth: 50,
			},
			check: func(t *testing.T, c *Config) {
				if c.TemplateDir != "custom" {
					t.Errorf("expected TemplateDir = 'custom', got %q", c.TemplateDir)
				}
				if len(c.Extensions) != 1 || c.Extensions[0] != ".tmpl" {
					t.Errorf("expected Extensions = ['.tmpl'], got %v", c.Extensions)
				}
				if c.LeftDelimiter != "{%" {
					t.Errorf("expected LeftDelimiter = '{%%', got %q", c.LeftDelimiter)
				}
				if c.RightDelimiter != "%}" {
					t.Errorf("expected RightDelimiter = '%%}', got %q", c.RightDelimiter)
				}
				if c.MaxIncludeDepth != 50 {
					t.Errorf("expected MaxIncludeDepth = 50, got %d", c.MaxIncludeDepth)
				}
			},
		},
		{
			name: "prefers TemplateFS over TemplateDir for defaults",
			config: Config{
				TemplateFS: fstest.MapFS{},
			},
			check: func(t *testing.T, c *Config) {
				if c.TemplateFS == nil {
					t.Error("expected TemplateFS to be preserved")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.config.applyDefaults()
			tt.check(t, &tt.config)
		})
	}
}

func TestConfigWithFS(t *testing.T) {
	testFS := fstest.MapFS{
		"test.html": &fstest.MapFile{
			Data: []byte("test content"),
		},
	}

	config := Config{
		TemplateFS: testFS,
	}

	config.applyDefaults()

	if config.TemplateFS == nil {
		t.Error("expected TemplateFS to be set")
	}

	if config.TemplateDir != "" {
		t.Errorf("expected TemplateDir to remain empty when TemplateFS is set, got %q", config.TemplateDir)
	}
}
