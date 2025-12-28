// Package main demonstrates HTML generation with Fíth templates.
// Shows how to build complete HTML pages with partials and layouts.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/toutaio/toutago-fith-renderer"
)

func main() {
	// Create temporary directory for templates
	tmpDir, err := os.MkdirTemp("", "fith-html-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create templates
	createHTMLTemplates(tmpDir)

	// Create engine
	engine, err := fith.New(fith.Config{
		TemplateDir: tmpDir,
		Extensions:  []string{".html"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// Example 1: Simple HTML page
	fmt.Println("=== Example 1: Simple HTML Page ===")
	output, err := engine.Render("simple", map[string]interface{}{
		"Title":   "My Page",
		"Content": "Hello, World!",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// Example 2: Product listing page
	fmt.Println("\n=== Example 2: Product Listing ===")
	output, err = engine.Render("products", map[string]interface{}{
		"Products": []map[string]interface{}{
			{"Name": "Laptop", "Price": 999.99, "InStock": true},
			{"Name": "Mouse", "Price": 29.99, "InStock": true},
			{"Name": "Keyboard", "Price": 79.99, "InStock": false},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)

	// Example 3: User dashboard with layout
	fmt.Println("\n=== Example 3: Dashboard with Layout ===")
	output, err = engine.Render("dashboard", map[string]interface{}{
		"User": map[string]interface{}{
			"Name":   "Alice",
			"Email":  "alice@example.com",
			"Avatar": "/avatars/alice.png",
		},
		"Notifications": []string{
			"New message from Bob",
			"Order #123 shipped",
			"Your subscription expires in 5 days",
		},
		"Stats": map[string]interface{}{
			"Orders":   42,
			"Revenue":  12345.67,
			"Visitors": 1523,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(output)
}

func createHTMLTemplates(dir string) {
	templates := map[string]string{
		// Simple page
		"simple.html": `<!DOCTYPE html>
<html>
<head>
  <title>{{.Title}}</title>
  <meta charset="utf-8">
</head>
<body>
  <h1>{{.Title}}</h1>
  <p>{{.Content}}</p>
</body>
</html>`,

		// Product card partial
		"_product-card.html": `<div class="product-card {{if not .InStock}}out-of-stock{{end}}">
  <h3>{{.Name}}</h3>
  <p class="price">${{.Price}}</p>
  {{if .InStock}}
    <button>Add to Cart</button>
  {{else}}
    <span class="badge">Out of Stock</span>
  {{end}}
</div>`,

		// Products page
		"products.html": `<!DOCTYPE html>
<html>
<head>
  <title>Products</title>
  <style>
    .product-card { border: 1px solid #ddd; padding: 1rem; margin: 0.5rem; }
    .out-of-stock { opacity: 0.6; }
    .price { font-weight: bold; color: #2c5; }
    .badge { background: #f44; color: white; padding: 0.2rem 0.5rem; }
  </style>
</head>
<body>
  <h1>Our Products</h1>
  <div class="products">
    {{range .Products}}
      {{include "_product-card.html" .}}
    {{end}}
  </div>
</body>
</html>`,

		// Base layout
		"_layout.html": `<!DOCTYPE html>
<html>
<head>
  <title>{{block "title"}}App{{end}}</title>
  <meta charset="utf-8">
  <style>
    body { font-family: Arial, sans-serif; margin: 0; }
    .header { background: #333; color: white; padding: 1rem; }
    .content { padding: 2rem; }
    .footer { background: #eee; padding: 1rem; text-align: center; }
  </style>
  {{block "head"}}{{end}}
</head>
<body>
  <div class="header">
    {{block "header"}}
      <h1>My App</h1>
    {{end}}
  </div>
  
  <div class="content">
    {{block "content"}}
      <p>No content</p>
    {{end}}
  </div>
  
  <div class="footer">
    {{block "footer"}}
      <p>&copy; 2024 My Company</p>
    {{end}}
  </div>
</body>
</html>`,

		// Dashboard page
		"dashboard.html": `{{extends "_layout.html"}}

{{block "title"}}Dashboard - {{.User.Name}}{{end}}

{{block "head"}}
<style>
  .user-info { display: flex; align-items: center; gap: 1rem; }
  .avatar { width: 40px; height: 40px; border-radius: 50%; }
  .notifications { background: #fff3cd; padding: 1rem; margin: 1rem 0; }
  .stats { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1rem; }
  .stat { background: white; border: 1px solid #ddd; padding: 1rem; text-align: center; }
  .stat-value { font-size: 2rem; font-weight: bold; color: #2c5; }
</style>
{{end}}

{{block "header"}}
<div class="user-info">
  <img src="{{.User.Avatar}}" alt="{{.User.Name}}" class="avatar">
  <div>
    <h2>{{.User.Name}}</h2>
    <p>{{.User.Email}}</p>
  </div>
</div>
{{end}}

{{block "content"}}
<h1>Dashboard</h1>

{{if .Notifications}}
<div class="notifications">
  <h3>Notifications ({{len .Notifications}})</h3>
  <ul>
    {{range .Notifications}}
      <li>{{.}}</li>
    {{end}}
  </ul>
</div>
{{end}}

<h2>Statistics</h2>
<div class="stats">
  <div class="stat">
    <div class="stat-value">{{.Stats.Orders}}</div>
    <div>Orders</div>
  </div>
  <div class="stat">
    <div class="stat-value">${{.Stats.Revenue}}</div>
    <div>Revenue</div>
  </div>
  <div class="stat">
    <div class="stat-value">{{.Stats.Visitors}}</div>
    <div>Visitors</div>
  </div>
</div>
{{end}}`,
	}

	for name, content := range templates {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			log.Fatal(err)
		}
	}
}
