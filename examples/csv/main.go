// Package main demonstrates CSV generation with Fíth templates.
// Shows how to generate CSV files for data export and reporting.
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/toutaio/toutago-fith-renderer"
)

func main() {
	engine, err := fith.NewWithDefaults()
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Simple CSV export
	fmt.Println("=== Example 1: User Export ===")
	csvTemplate := `Name,Email,Age,Status
{{range .Users}}{{.Name}},{{.Email}},{{.Age}},{{.Status}}
{{end}}`

	users := map[string]interface{}{
		"Users": []map[string]interface{}{
			{"Name": "Alice", "Email": "alice@example.com", "Age": 30, "Status": "Active"},
			{"Name": "Bob", "Email": "bob@example.com", "Age": 25, "Status": "Active"},
			{"Name": "Charlie", "Email": "charlie@example.com", "Age": 35, "Status": "Inactive"},
		},
	}

	output, err := engine.RenderString(csvTemplate, users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// Example 2: Sales report with calculations
	fmt.Println("=== Example 2: Sales Report ===")
	salesTemplate := `Date,Product,Quantity,Price,Total
{{range .Sales}}{{date "YYYY-MM-DD" .Date}},{{.Product}},{{.Quantity}},{{.Price}},{{.Total}}
{{end}}`

	sales := map[string]interface{}{
		"Sales": []map[string]interface{}{
			{
				"Date":     time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
				"Product":  "Laptop",
				"Quantity": 2,
				"Price":    999.99,
				"Total":    1999.98,
			},
			{
				"Date":     time.Date(2024, 12, 2, 0, 0, 0, 0, time.UTC),
				"Product":  "Mouse",
				"Quantity": 5,
				"Price":    29.99,
				"Total":    149.95,
			},
			{
				"Date":     time.Date(2024, 12, 3, 0, 0, 0, 0, time.UTC),
				"Product":  "Keyboard",
				"Quantity": 3,
				"Price":    79.99,
				"Total":    239.97,
			},
		},
	}

	output, err = engine.RenderString(salesTemplate, sales)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// Example 3: Report with conditional columns
	fmt.Println("=== Example 3: Conditional Columns ===")
	conditionalTemplate := `ID,Name,Email{{if .IncludePhone}},Phone{{end}}{{if .IncludeAddress}},Address{{end}}
{{range .Contacts}}{{.ID}},{{.Name}},{{.Email}}{{if $.IncludePhone}},{{.Phone}}{{end}}{{if $.IncludeAddress}},{{.Address}}{{end}}
{{end}}`

	contacts := map[string]interface{}{
		"IncludePhone":   true,
		"IncludeAddress": false,
		"Contacts": []map[string]interface{}{
			{
				"ID":      1,
				"Name":    "Alice Smith",
				"Email":   "alice@example.com",
				"Phone":   "555-0001",
				"Address": "123 Main St",
			},
			{
				"ID":      2,
				"Name":    "Bob Jones",
				"Email":   "bob@example.com",
				"Phone":   "555-0002",
				"Address": "456 Oak Ave",
			},
		},
	}

	output, err = engine.RenderString(conditionalTemplate, contacts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// Example 4: Escaped CSV values
	fmt.Println("=== Example 4: Escaped Values ===")

	// Register custom CSV escape function
	engine.RegisterFunction("csvEscape", func(args ...interface{}) (interface{}, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("csvEscape expects 1 argument")
		}
		str := fmt.Sprintf("%v", args[0])
		// Simple CSV escaping: wrap in quotes if contains comma, newline, or quote
		needsEscape := false
		for _, ch := range str {
			if ch == ',' || ch == '\n' || ch == '"' {
				needsEscape = true
				break
			}
		}
		if needsEscape {
			// Escape quotes by doubling them
			escaped := ""
			for _, ch := range str {
				if ch == '"' {
					escaped += "\"\""
				} else {
					escaped += string(ch)
				}
			}
			return "\"" + escaped + "\"", nil
		}
		return str, nil
	})

	escapedTemplate := `Name,Description,Notes
{{range .Items}}{{csvEscape .Name}},{{csvEscape .Description}},{{csvEscape .Notes}}
{{end}}`

	items := map[string]interface{}{
		"Items": []map[string]interface{}{
			{
				"Name":        "Product A",
				"Description": "A simple product",
				"Notes":       "No issues",
			},
			{
				"Name":        "Product B",
				"Description": "A product with a comma, in its description",
				"Notes":       "Needs attention",
			},
			{
				"Name":        "Product C",
				"Description": "A product with \"quotes\"",
				"Notes":       "Special handling required",
			},
		},
	}

	output, err = engine.RenderString(escapedTemplate, items)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
