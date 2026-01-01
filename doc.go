/*
Package fith provides a powerful, flexible template engine for Go.

Fíth (Old Irish: "The art of weaving patterns") is a template engine inspired by
Jinja2 and Twig, designed for generating HTML, text, and other formats from templates.

# Quick Start

Create a new engine and render a template:

	engine, err := fith.New(&fith.Config{
	    TemplateDir: "templates",
	})
	if err != nil {
	    log.Fatal(err)
	}

	data := map[string]interface{}{
	    "Title": "Welcome",
	    "User":  user,
	}

	output, err := engine.Render("home", data)
	if err != nil {
	    log.Fatal(err)
	}
	fmt.Println(output)

# Template Syntax

Variables are accessed using {{.Variable}}:

	Hello {{.Name}}!
	User: {{.User.Profile.Name}}

Conditionals use if/else:

	{{if .IsLoggedIn}}
	    Welcome back!
	{{else}}
	    Please log in
	{{end}}

Loops use range:

	{{range .Items}}
	    Item: {{.}}
	    Index: {{@index}}
	{{end}}

Functions can be called directly or as filters:

	{{upper .Name}}
	{{.Name | upper | trim}}

Templates can include other templates:

	{{include "header"}}
	{{include "card" title="Hello" content=.Message}}

Templates support layouts with blocks:

	{{extends "layouts/main"}}

	{{block "title"}}My Page{{end}}

	{{block "content"}}
	    <p>Page content here</p>
	{{end}}

# Configuration

The engine can be configured with various options:

	config := fith.Config{
	    TemplateDir:     "templates",
	    Extensions:      []string{".html", ".tpl"},
	    CacheEnabled:    true,
	    MaxIncludeDepth: 100,
	}

	engine, err := fith.New(&config)

# Custom Functions

Register custom functions for use in templates:

	engine.RegisterFunction("reverse", func(args ...interface{}) (interface{}, error) {
	    s := fmt.Sprint(args[0])
	    runes := []rune(s)
	    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
	        runes[i], runes[j] = runes[j], runes[i]
	    }
	    return string(runes), nil
	})

Then use it in templates:

	{{reverse .Text}}

# Built-in Functions

String functions: upper, lower, title, trim, truncate, replace
Array functions: join, first, last, len
Logic functions: default
Encoding: urlEncode, htmlEscape
Date: date

# Error Handling

The engine provides detailed error information:

	output, err := engine.Render("template", data)
	if err != nil {
	    if fithErr, ok := err.(*fith.Error); ok {
	        fmt.Printf("Error in %s at line %d: %s\n",
	            fithErr.Slug, fithErr.Line, fithErr.Message)
	    }
	}

# Thread Safety

The engine is safe for concurrent use. Multiple goroutines can render
templates simultaneously using the same engine instance.

# Performance

Templates are compiled and cached for fast rendering. The cache can be
cleared if templates are modified:

	engine.ClearCache()

Benchmarks show rendering performance of ~2-3μs for typical templates
with caching enabled.
*/
package fith
