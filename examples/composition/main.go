// Package main demonstrates template composition features in Fíth:
// - Includes with parameters
// - Layout inheritance with blocks
// - Nested composition
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/toutaio/toutago-fith-renderer/loader"
	"github.com/toutaio/toutago-fith-renderer/runtime"
)

func main() {
	// Create temporary directory for templates
	tmpDir, err := os.MkdirTemp("", "fith-composition-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create template files
	createTemplates(tmpDir)

	// Create loader
	ldr := loader.NewFileSystemLoader(tmpDir, []string{".html"})

	// Example 1: Simple Include
	fmt.Println("=== Example 1: Simple Include ===")
	runExample(ldr, "example1.html", map[string]interface{}{
		"siteName": "My Website",
	})

	// Example 2: Include with Parameters
	fmt.Println("\n=== Example 2: Include with Parameters ===")
	runExample(ldr, "example2.html", map[string]interface{}{})

	// Example 3: Layout Inheritance
	fmt.Println("\n=== Example 3: Layout Inheritance ===")
	runExample(ldr, "page.html", map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "Alice",
			"email": "alice@example.com",
		},
		"items": []string{"Item 1", "Item 2", "Item 3"},
	})

	// Example 4: Multiple Blocks
	fmt.Println("\n=== Example 4: Multiple Blocks ===")
	runExample(ldr, "blog-post.html", map[string]interface{}{
		"post": map[string]interface{}{
			"title":   "My First Post",
			"content": "This is the content of my first blog post.",
			"author":  "Bob",
		},
	})
}

func createTemplates(dir string) {
	templates := map[string]string{
		// Partials for includes
		"header.html": `<header>{{.siteName}}</header>`,

		"footer.html": `<footer>© 2024</footer>`,

		"card.html": `<div class="card">
  <h3>{{.title}}</h3>
  <p>{{.content}}</p>
</div>`,

		"user-info.html": `<div class="user">
  <span>{{.name}}</span>
  <span>{{.email}}</span>
</div>`,

		// Example 1: Simple include
		"example1.html": `{{include "header.html"}}
<main>Welcome!</main>
{{include "footer.html"}}`,

		// Example 2: Include with parameters
		"example2.html": `{{include "card.html" title="Hello" content="World"}}
{{include "card.html" title="Goodbye" content="For now"}}`,

		// Layout template
		"layout.html": `<!DOCTYPE html>
<html>
<head>
  <title>{{block "title"}}Default Title{{end}}</title>
</head>
<body>
  <header>
    {{block "header"}}
      <h1>Welcome</h1>
    {{end}}
  </header>
  
  <main>
    {{block "content"}}
      <p>Default content</p>
    {{end}}
  </main>
  
  <footer>
    {{block "footer"}}
      <p>© 2024</p>
    {{end}}
  </footer>
</body>
</html>`,

		// Example 3: Page extending layout
		"page.html": `{{extends "layout.html"}}

{{block "title"}}User Dashboard{{end}}

{{block "header"}}
  <h1>Dashboard</h1>
  {{include "user-info.html" .user}}
{{end}}

{{block "content"}}
  <h2>Your Items:</h2>
  <ul>
  {{range .items}}
    <li>{{.}}</li>
  {{end}}
  </ul>
{{end}}`,

		// Blog layout
		"blog-layout.html": `<!DOCTYPE html>
<html>
<head>
  <title>Blog - {{block "post-title"}}Untitled{{end}}</title>
</head>
<body>
  <nav>Blog Navigation</nav>
  
  <article>
    {{block "article"}}
      <p>No article content</p>
    {{end}}
  </article>
  
  <aside>
    {{block "sidebar"}}
      <p>Sidebar content</p>
    {{end}}
  </aside>
</body>
</html>`,

		// Example 4: Blog post
		"blog-post.html": `{{extends "blog-layout.html"}}

{{block "post-title"}}{{.post.title}}{{end}}

{{block "article"}}
  <h1>{{.post.title}}</h1>
  <p class="author">By {{.post.author}}</p>
  <div class="content">
    {{.post.content}}
  </div>
{{end}}

{{block "sidebar"}}
  <h3>About the Author</h3>
  <p>{{.post.author}} is a writer.</p>
{{end}}`,
	}

	for name, content := range templates {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			log.Fatal(err)
		}
	}
}

func runExample(ldr *loader.FileSystemLoader, templateName string, data map[string]interface{}) {
	tmpl, err := ldr.Load(templateName)
	if err != nil {
		log.Fatalf("Failed to load template %q: %v", templateName, err)
	}

	ctx := runtime.NewContext(data)
	output, err := runtime.ExecuteWithLoader(tmpl, ctx, ldr)
	if err != nil {
		log.Fatalf("Failed to execute template %q: %v", templateName, err)
	}

	fmt.Println(output)
}
