package main

import (
	"fmt"
	"log"

	"github.com/toutaio/toutago-fith-renderer"
)

func main() {
	// Create a new Fíth engine with default settings
	engine, err := fith.NewWithDefaults()
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Simple variable substitution
	fmt.Println("=== Example 1: Simple Variables ===")
	template1 := "Hello {{.Name}}! You are {{.Age}} years old."
	data1 := map[string]interface{}{
		"Name": "Alice",
		"Age":  30,
	}
	output1, err := engine.RenderString(template1, data1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output1)
	fmt.Println()

	// Example 2: Conditionals
	fmt.Println("=== Example 2: Conditionals ===")
	template2 := `{{if .IsLoggedIn}}Welcome back, {{.Username}}!{{else}}Please log in.{{end}}`
	data2a := map[string]interface{}{
		"IsLoggedIn": true,
		"Username":   "bob",
	}
	output2a, _ := engine.RenderString(template2, data2a)
	fmt.Println("Logged in:", output2a)

	data2b := map[string]interface{}{
		"IsLoggedIn": false,
	}
	output2b, _ := engine.RenderString(template2, data2b)
	fmt.Println("Not logged in:", output2b)
	fmt.Println()

	// Example 3: Loops
	fmt.Println("=== Example 3: Loops ===")
	template3 := `Items: {{range .Items}}{{.}} {{end}}`
	data3 := map[string]interface{}{
		"Items": []string{"apple", "banana", "cherry"},
	}
	output3, _ := engine.RenderString(template3, data3)
	fmt.Println(output3)
	fmt.Println()

	// Example 4: Built-in functions
	fmt.Println("=== Example 4: Built-in Functions ===")
	template4 := `
Upper: {{upper .Text}}
Lower: {{lower .Text}}
Title: {{title .Text}}
Joined: {{join .Words ", "}}
`
	data4 := map[string]interface{}{
		"Text":  "hello world",
		"Words": []string{"foo", "bar", "baz"},
	}
	output4, _ := engine.RenderString(template4, data4)
	fmt.Println(output4)

	// Example 5: Filter pipelines
	fmt.Println("=== Example 5: Filter Pipelines ===")
	template5 := `{{.Text | upper | trim}}`
	data5 := map[string]interface{}{
		"Text": "  hello world  ",
	}
	output5, _ := engine.RenderString(template5, data5)
	fmt.Println("Result:", output5)
	fmt.Println()

	// Example 6: Custom functions
	fmt.Println("=== Example 6: Custom Functions ===")
	engine.RegisterFunction("double", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("double expects 1 argument")
		}
		if n, ok := args[0].(int); ok {
			return n * 2, nil
		}
		return nil, fmt.Errorf("double expects an integer")
	})

	template6 := `{{.Value}} doubled is {{double .Value}}`
	data6 := map[string]interface{}{
		"Value": 21,
	}
	output6, _ := engine.RenderString(template6, data6)
	fmt.Println(output6)
	fmt.Println()

	// Example 7: Nested data structures
	fmt.Println("=== Example 7: Nested Data ===")
	template7 := `
User: {{.User.Name}}
Email: {{.User.Email}}
Address: {{.User.Address.City}}, {{.User.Address.Country}}
`
	data7 := map[string]interface{}{
		"User": map[string]interface{}{
			"Name":  "Charlie",
			"Email": "charlie@example.com",
			"Address": map[string]interface{}{
				"City":    "San Francisco",
				"Country": "USA",
			},
		},
	}
	output7, _ := engine.RenderString(template7, data7)
	fmt.Println(output7)
}
