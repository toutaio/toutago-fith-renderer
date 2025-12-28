# HTML Generation Example

This example demonstrates how to use Fíth to generate complete HTML pages.

## Features Demonstrated

- Basic HTML page structure
- Product listing with conditionals
- Layout inheritance
- Partials/includes
- CSS styling in templates
- Dynamic content rendering

## Running the Example

```bash
go run main.go
```

## Key Concepts

### 1. Simple HTML Pages

Generate complete HTML documents with dynamic content:

```go
engine.Render("simple", map[string]interface{}{
    "Title": "My Page",
    "Content": "Hello, World!",
})
```

### 2. Partials for Reusable Components

Create reusable components like product cards:

```html
<!-- _product-card.html -->
<div class="product-card">
  <h3>{{.Name}}</h3>
  <p class="price">${{.Price}}</p>
</div>
```

Use them with `include`:

```html
{{range .Products}}
  {{include "_product-card.html" .}}
{{end}}
```

### 3. Layout Inheritance

Define a base layout with blocks:

```html
<!-- _layout.html -->
<!DOCTYPE html>
<html>
<head>
  <title>{{block "title"}}Default{{end}}</title>
</head>
<body>
  {{block "content"}}{{end}}
</body>
</html>
```

Extend it in pages:

```html
{{extends "_layout.html"}}
{{block "title"}}My Page{{end}}
{{block "content"}}<p>Content here</p>{{end}}
```

## Use Cases

- Web applications
- Static site generation
- Email HTML templates
- Admin dashboards
- Product catalogs
- Landing pages
