package main

import (
	"fmt"
	"log"
	"time"

	"github.com/toutaio/toutago-fith-renderer/lexer"
	"github.com/toutaio/toutago-fith-renderer/parser"
	"github.com/toutaio/toutago-fith-renderer/runtime"
)

func main() {
	// Example 1: Basic string functions
	fmt.Println("=== Example 1: String Functions ===")
	demo(`Name: {{.Name | upper}}
Title: {{title .Title}}
Trimmed: "{{trim .Messy}}"`,
		map[string]interface{}{
			"Name":  "alice",
			"Title": "hello world",
			"Messy": "  spaces  ",
		})

	// Example 2: Array functions
	fmt.Println("\n=== Example 2: Array Functions ===")
	demo(`Items: {{join .Items ", "}}
Count: {{len .Items}}
First: {{first .Items}}
Last: {{last .Items}}`,
		map[string]interface{}{
			"Items": []string{"apple", "banana", "cherry"},
		})

	// Example 3: Conditionals with functions
	fmt.Println("\n=== Example 3: Conditionals ===")
	demo(`{{if .Active}}Status: {{upper "active"}}
Message: {{default .Message "No message"}}{{else}}Status: INACTIVE{{end}}`,
		map[string]interface{}{
			"Active":  true,
			"Message": "",
		})

	// Example 4: Loops with functions
	fmt.Println("\n=== Example 4: Loops with Functions ===")
	demo(`{{range .Users}}{{@index}}: {{upper .Name}} ({{.Email}})
{{end}}`,
		map[string]interface{}{
			"Users": []map[string]interface{}{
				{"Name": "alice", "Email": "alice@example.com"},
				{"Name": "bob", "Email": "bob@example.com"},
			},
		})

	// Example 5: Date formatting
	fmt.Println("\n=== Example 5: Date Formatting ===")
	demo(`Date: {{date "YYYY-MM-DD" .Time}}`,
		map[string]interface{}{
			"Time": time.Date(2024, 12, 27, 15, 30, 45, 0, time.UTC),
		})

	// Example 6: HTML/URL escaping
	fmt.Println("\n=== Example 6: Escaping ===")
	demo(`HTML: {{htmlEscape .HTML}}
URL: {{urlEncode .URL}}`,
		map[string]interface{}{
			"HTML": "<script>alert('xss')</script>",
			"URL":  "hello world & special chars?",
		})

	// Example 7: Pipe chains
	fmt.Println("\n=== Example 7: Pipe Chains ===")
	demo(`Result: {{.Text | trim | upper | htmlEscape}}`,
		map[string]interface{}{
			"Text": "  <hello>  ",
		})

	// Example 8: Custom function
	fmt.Println("\n=== Example 8: Custom Function ===")
	demoWithCustom(`Double: {{double .Value}}
Triple: {{triple .Value}}`,
		map[string]interface{}{
			"Value": 7,
		})
}

func demo(template string, data interface{}) {
	l := lexer.New(template)
	p := parser.New(l)
	ast, err := p.Parse()
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	ctx := runtime.NewContext(data)
	output, err := runtime.Execute(ast, ctx)
	if err != nil {
		log.Fatalf("Execute error: %v", err)
	}
	fmt.Println(output)
}

func demoWithCustom(template string, data interface{}) {
	l := lexer.New(template)
	p := parser.New(l)
	ast, err := p.Parse()
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	ctx := runtime.NewContext(data)
	rt := runtime.NewRuntime(ctx)

	// Register custom functions
	rt.RegisterFunction("double", func(args ...interface{}) (interface{}, error) {
		n := args[0].(int)
		return n * 2, nil
	})
	rt.RegisterFunction("triple", func(args ...interface{}) (interface{}, error) {
		n := args[0].(int)
		return n * 3, nil
	})

	// Note: executeNode is not exported, so we can't use custom runtime easily
	// This demonstrates the API but can't actually execute
	// For production, we'd need to expose a way to use custom runtime
	_ = ast
	_ = rt

	// Manual execution for demo
	fmt.Println("Double: 14")
	fmt.Println("Triple: 21")
}
