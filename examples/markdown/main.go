// Package main demonstrates Markdown generation with Fíth templates.
// Shows how to generate documentation, blog posts, and README files.
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

	// Example 1: Blog post
	fmt.Println("=== Example 1: Blog Post ===")
	blogTemplate := `# {{.Title}}

**Author:** {{.Author}}  
**Date:** {{date "January 2, 2006" .Date}}

{{.Content}}

## Tags

{{range .Tags}}• {{.}}
{{end}}

---
*{{.ReadTime}} min read*`

	blog := map[string]interface{}{
		"Title":    "Getting Started with Fíth",
		"Author":   "Alice Developer",
		"Date":     time.Now(),
		"Content":  "Fíth is a powerful template engine for Go...",
		"Tags":     []string{"golang", "templates", "tutorial"},
		"ReadTime": 5,
	}

	output, err := engine.RenderString(blogTemplate, blog)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	fmt.Println()

	// Example 2: API Documentation
	fmt.Println("=== Example 2: API Documentation ===")
	apiTemplate := `# {{.ServiceName}} API

Version: {{.Version}}

## Endpoints

{{range .Endpoints}}### {{upper .Method}} {{.Path}}

{{.Description}}

**Parameters:**

{{if .Params}}| Name | Type | Required | Description |
|------|------|----------|-------------|
{{range .Params}}| {{.Name}} | {{.Type}} | {{if .Required}}Yes{{else}}No{{end}} | {{.Description}} |
{{end}}{{else}}*No parameters*
{{end}}

**Example:**

` + "```" + `{{.Language}}
{{.Example}}
` + "```" + `

{{end}}`

	api := map[string]interface{}{
		"ServiceName": "User Service",
		"Version":     "v1.0.0",
		"Endpoints": []map[string]interface{}{
			{
				"Method":      "GET",
				"Path":        "/users/{id}",
				"Description": "Retrieve a user by ID",
				"Language":    "bash",
				"Example":     "curl https://api.example.com/users/123",
				"Params": []map[string]interface{}{
					{
						"Name":        "id",
						"Type":        "string",
						"Required":    true,
						"Description": "User ID",
					},
				},
			},
			{
				"Method":      "POST",
				"Path":        "/users",
				"Description": "Create a new user",
				"Language":    "json",
				"Example":     `{"name": "Alice", "email": "alice@example.com"}`,
				"Params": []map[string]interface{}{
					{
						"Name":        "name",
						"Type":        "string",
						"Required":    true,
						"Description": "User's full name",
					},
					{
						"Name":        "email",
						"Type":        "string",
						"Required":    true,
						"Description": "User's email address",
					},
				},
			},
		},
	}

	output, err = engine.RenderString(apiTemplate, api)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	fmt.Println()

	// Example 3: README generator
	fmt.Println("=== Example 3: README Generator ===")
	readmeTemplate := `# {{.ProjectName}}

{{.Description}}

## Installation

` + "```" + `bash
{{.InstallCommand}}
` + "```" + `

## Features

{{range .Features}}• {{.}}
{{end}}

## Quick Start

` + "```" + `{{.Language}}
{{.QuickStartCode}}
` + "```" + `

## Documentation

{{range .DocLinks}}- [{{.Name}}]({{.URL}})
{{end}}

## License

{{.License}}
`

	readme := map[string]interface{}{
		"ProjectName":    "awesome-project",
		"Description":    "An awesome project that does amazing things",
		"InstallCommand": "go get github.com/user/awesome-project",
		"Features": []string{
			"Fast and efficient",
			"Easy to use",
			"Well documented",
			"Fully tested",
		},
		"Language":       "go",
		"QuickStartCode": "import \"github.com/user/awesome-project\"\n\nfunc main() {\n    // Your code here\n}",
		"DocLinks": []map[string]interface{}{
			{"Name": "API Reference", "URL": "/docs/api"},
			{"Name": "User Guide", "URL": "/docs/guide"},
			{"Name": "Examples", "URL": "/examples"},
		},
		"License": "MIT",
	}

	output, err = engine.RenderString(readmeTemplate, readme)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
	fmt.Println()

	// Example 4: Changelog generator
	fmt.Println("=== Example 4: Changelog ===")
	changelogTemplate := `# Changelog

{{range .Versions}}## [{{.Version}}] - {{date "2006-01-02" .Date}}

{{if .Added}}### Added
{{range .Added}}- {{.}}
{{end}}
{{end}}{{if .Changed}}### Changed
{{range .Changed}}- {{.}}
{{end}}
{{end}}{{if .Fixed}}### Fixed
{{range .Fixed}}- {{.}}
{{end}}
{{end}}{{if .Deprecated}}### Deprecated
{{range .Deprecated}}- {{.}}
{{end}}
{{end}}
{{end}}`

	changelog := map[string]interface{}{
		"Versions": []map[string]interface{}{
			{
				"Version": "1.1.0",
				"Date":    time.Date(2024, 12, 27, 0, 0, 0, 0, time.UTC),
				"Added": []string{
					"New feature X",
					"Support for Y",
				},
				"Changed": []string{
					"Improved performance of Z",
				},
				"Fixed": []string{
					"Bug in feature A",
					"Memory leak in component B",
				},
			},
			{
				"Version": "1.0.0",
				"Date":    time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
				"Added": []string{
					"Initial release",
					"Core functionality",
				},
			},
		},
	}

	output, err = engine.RenderString(changelogTemplate, changelog)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}
