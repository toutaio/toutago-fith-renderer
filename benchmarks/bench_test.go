package benchmarks

import (
	"testing"

	"github.com/toutaio/toutago-fith-renderer/compiler"
	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
	"github.com/toutaio/toutago-fith-renderer/runtime"
)

// Benchmark simple variable rendering
func BenchmarkSimpleVariable(b *testing.B) {
	template := "Hello, {{.Name}}!"
	data := map[string]interface{}{
		"Name": "World",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lex := lexer.New(template)
		p := parser.New(lex)
		tmpl, _ := p.Parse()
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark with parsing cached (typical case)
func BenchmarkSimpleVariable_Cached(b *testing.B) {
	template := "Hello, {{.Name}}!"
	data := map[string]interface{}{
		"Name": "World",
	}

	// Parse once
	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark complex template with loops and conditionals
func BenchmarkComplexTemplate(b *testing.B) {
	template := `
{{if .ShowHeader}}
<header>
  <h1>{{.Title}}</h1>
</header>
{{end}}
<ul>
{{range .Items}}
  <li>{{.}} - {{@index}}</li>
{{end}}
</ul>
`
	data := map[string]interface{}{
		"ShowHeader": true,
		"Title":      "My Page",
		"Items":      []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lex := lexer.New(template)
		p := parser.New(lex)
		tmpl, _ := p.Parse()
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark complex template with caching
func BenchmarkComplexTemplate_Cached(b *testing.B) {
	template := `
{{if .ShowHeader}}
<header>
  <h1>{{.Title}}</h1>
</header>
{{end}}
<ul>
{{range .Items}}
  <li>{{.}} - {{@index}}</li>
{{end}}
</ul>
`
	data := map[string]interface{}{
		"ShowHeader": true,
		"Title":      "My Page",
		"Items":      []string{"Item 1", "Item 2", "Item 3", "Item 4", "Item 5"},
	}

	// Parse once
	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark nested data access
func BenchmarkNestedAccess(b *testing.B) {
	template := "{{.User.Profile.Name}} - {{.User.Profile.Email}}"
	data := map[string]interface{}{
		"User": map[string]interface{}{
			"Profile": map[string]interface{}{
				"Name":  "John Doe",
				"Email": "john@example.com",
			},
		},
	}

	// Parse once
	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark filter pipeline
func BenchmarkFilterPipeline(b *testing.B) {
	template := "{{.Text | upper | trim}}"
	data := map[string]interface{}{
		"Text": "  hello world  ",
	}

	// Parse once
	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark large loop
func BenchmarkLargeLoop(b *testing.B) {
	template := `{{range .Items}}{{.}}{{end}}`

	// Create 100 items
	items := make([]string, 100)
	for i := 0; i < 100; i++ {
		items[i] = "Item"
	}

	data := map[string]interface{}{
		"Items": items,
	}

	// Parse once
	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark lexer only
func BenchmarkLexer(b *testing.B) {
	template := `
{{if .ShowHeader}}
<header><h1>{{.Title}}</h1></header>
{{end}}
<ul>
{{range .Items}}
  <li>{{.}} - {{@index}}</li>
{{end}}
</ul>
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lex := lexer.New(template)
		for {
			tok, _ := lex.NextToken()
			if tok.Type == lexer.TokenEOF {
				break
			}
		}
	}
}

// Benchmark parser only
func BenchmarkParser(b *testing.B) {
	template := `
{{if .ShowHeader}}
<header><h1>{{.Title}}</h1></header>
{{end}}
<ul>
{{range .Items}}
  <li>{{.}} - {{@index}}</li>
{{end}}
</ul>
`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lex := lexer.New(template)
		p := parser.New(lex)
		p.Parse()
	}
}

// mockLoader for compiler benchmark
type mockCompilerLoader struct {
	tmpl *parser.Template
}

func (m *mockCompilerLoader) Load(slug string) (*parser.Template, error) {
	return m.tmpl, nil
}

func (m *mockCompilerLoader) Exists(slug string) bool {
	return true
}

// Benchmark compiler optimization
func BenchmarkCompiler(b *testing.B) {
	// Create test template
	template := `{{if true}}Always shown{{else}}Never shown{{end}}`
	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	// Create mock loader
	ml := &mockCompilerLoader{tmpl: tmpl}
	c := compiler.New(ml)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c.ClearCache()
		c.Compile("test")
	}
}

// Benchmark string builder output
func BenchmarkOutputBuilder(b *testing.B) {
	template := `{{range .Items}}{{.}}{{end}}`

	items := make([]string, 50)
	for i := 0; i < 50; i++ {
		items[i] = "Content "
	}

	data := map[string]interface{}{
		"Items": items,
	}

	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark function calls
func BenchmarkFunctionCalls(b *testing.B) {
	template := `{{upper "hello"}} {{lower "WORLD"}} {{len .Items}}`

	data := map[string]interface{}{
		"Items": []int{1, 2, 3, 4, 5},
	}

	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}

// Benchmark memory allocations
func BenchmarkMemoryAllocations(b *testing.B) {
	template := `Hello, {{.Name}}!`
	data := map[string]interface{}{
		"Name": "World",
	}

	lex := lexer.New(template)
	p := parser.New(lex)
	tmpl, _ := p.Parse()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := runtime.NewContext(data)
		runtime.Execute(tmpl, ctx)
	}
}
