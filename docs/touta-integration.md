# Toutā Integration Guide

This guide shows how to integrate Fíth into the Toutā web framework.

## Installation

```bash
go get github.com/toutaio/toutago-fith-renderer
```

## Basic Integration

### 1. Add Fíth to Your Toutā App

```go
package main

import (
    "github.com/toutaio/toutago"
    "github.com/toutaio/toutago-fith-renderer"
)

func main() {
    app := toutago.New()
    
    // Initialize Fíth
    engine, err := fith.New(fith.Config{
        TemplateDir: "templates",
        Extensions:  []string{".html", ".fith"},
        CacheEnabled: true,  // Enable for production
    })
    if err != nil {
        panic(err)
    }
    
    // Set as template engine for Toutā
    app.SetTemplateEngine(engine)
    
    // Your routes
    app.GET("/", homeHandler)
    app.GET("/users/:id", userHandler)
    
    app.Run(":8080")
}
```

### 2. Create Template Directory Structure

```
your-app/
├── main.go
├── handlers/
└── templates/
    ├── layouts/
    │   └── main.html
    ├── pages/
    │   ├── home.html
    │   └── user.html
    └── partials/
        ├── header.html
        └── footer.html
```

### 3. Render Templates in Handlers

```go
func homeHandler(c *toutago.Context) error {
    data := map[string]interface{}{
        "Title": "Home",
        "User":  c.User(),
    }
    
    return c.Render("pages/home", data)
}

func userHandler(c *toutago.Context) error {
    userID := c.Param("id")
    user := getUserByID(userID)
    
    data := map[string]interface{}{
        "Title": user.Name,
        "User":  user,
    }
    
    return c.Render("pages/user", data)
}
```

## Template Examples

### Layout Template

**templates/layouts/main.html:**

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{block "title"}}My App{{end}}</title>
    <link rel="stylesheet" href="/static/css/main.css">
    {{block "head"}}{{end}}
</head>
<body>
    {{include "partials/header.html"}}
    
    <main>
        {{block "content"}}
            <p>No content</p>
        {{end}}
    </main>
    
    {{include "partials/footer.html"}}
    
    <script src="/static/js/main.js"></script>
    {{block "scripts"}}{{end}}
</body>
</html>
```

### Page Template

**templates/pages/home.html:**

```html
{{extends "layouts/main.html"}}

{{block "title"}}Home - My App{{end}}

{{block "content"}}
<div class="hero">
    <h1>Welcome{{if .User}}, {{.User.Name}}{{end}}!</h1>
    <p>This is the Toutā + Fíth demo</p>
</div>

<div class="features">
    {{range .Features}}
    <div class="feature-card">
        <h3>{{.Title}}</h3>
        <p>{{.Description}}</p>
    </div>
    {{end}}
</div>
{{end}}

{{block "scripts"}}
<script src="/static/js/home.js"></script>
{{end}}
```

### Partial Template

**templates/partials/header.html:**

```html
<header>
    <nav>
        <a href="/" class="logo">{{.AppName}}</a>
        <ul>
            <li><a href="/">Home</a></li>
            <li><a href="/about">About</a></li>
            {{if .User}}
                <li><a href="/dashboard">Dashboard</a></li>
                <li><a href="/logout">Logout</a></li>
            {{else}}
                <li><a href="/login">Login</a></li>
            {{end}}
        </ul>
    </nav>
</header>
```

## Advanced Integration

### Custom Context Data

Add data available to all templates:

```go
// Middleware to add common template data
func templateDataMiddleware(next toutago.HandlerFunc) toutago.HandlerFunc {
    return func(c *toutago.Context) error {
        c.Set("AppName", "MyApp")
        c.Set("Version", "1.0.0")
        c.Set("Year", time.Now().Year())
        return next(c)
    }
}

app.Use(templateDataMiddleware)
```

### Custom Functions

Register custom template functions:

```go
engine.RegisterFunction("route", func(args ...interface{}) (interface{}, error) {
    name := args[0].(string)
    return app.RouteURL(name), nil
})

engine.RegisterFunction("asset", func(args ...interface{}) (interface{}, error) {
    path := args[0].(string)
    return "/static/" + path, nil
})
```

Use in templates:

```html
<a href="{{route "user.profile" .User.ID}}">Profile</a>
<img src="{{asset "images/logo.png"}}" alt="Logo">
```

### Flash Messages

```go
// In handler
func loginHandler(c *toutago.Context) error {
    // ... login logic ...
    
    c.Flash("success", "Login successful!")
    return c.Redirect("/dashboard")
}

// Middleware to pass flash to templates
func flashMiddleware(next toutago.HandlerFunc) toutago.HandlerFunc {
    return func(c *toutago.Context) error {
        c.Set("Flash", c.GetFlash())
        return next(c)
    }
}
```

Template usage:

```html
{{if .Flash.success}}
<div class="alert alert-success">{{.Flash.success}}</div>
{{end}}

{{if .Flash.error}}
<div class="alert alert-error">{{.Flash.error}}</div>
{{end}}
```

### Form Helpers

```go
engine.RegisterFunction("csrfToken", func(args ...interface{}) (interface{}, error) {
    c := args[0].(*toutago.Context)
    return c.CSRFToken(), nil
})

engine.RegisterFunction("csrfField", func(args ...interface{}) (interface{}, error) {
    c := args[0].(*toutago.Context)
    return fmt.Sprintf(`<input type="hidden" name="_csrf" value="%s">`, c.CSRFToken()), nil
})
```

Template:

```html
<form method="POST" action="/login">
    {{csrfField .}}
    <input type="email" name="email" required>
    <input type="password" name="password" required>
    <button type="submit">Login</button>
</form>
```

## Production Configuration

```go
func setupTemplates(app *toutago.App, env string) error {
    config := fith.Config{
        TemplateDir:  "templates",
        Extensions:   []string{".html"},
        CacheEnabled: env == "production",
        AutoReload:   env == "development",
    }
    
    engine, err := fith.New(config)
    if err != nil {
        return err
    }
    
    // Add production-specific functions
    if env == "production" {
        engine.RegisterFunction("minify", minifyHTML)
    }
    
    app.SetTemplateEngine(engine)
    return nil
}
```

## Error Handling

### Custom Error Pages

**templates/errors/404.html:**

```html
{{extends "layouts/main.html"}}

{{block "title"}}Page Not Found{{end}}

{{block "content"}}
<div class="error-page">
    <h1>404</h1>
    <p>The page you're looking for doesn't exist.</p>
    <a href="/" class="button">Go Home</a>
</div>
{{end}}
```

Handler:

```go
app.SetErrorHandler(func(err error, c *toutago.Context) error {
    code := c.StatusCode()
    
    data := map[string]interface{}{
        "Error": err.Error(),
        "Code":  code,
    }
    
    template := fmt.Sprintf("errors/%d", code)
    if err := c.Render(template, data); err != nil {
        // Fallback to JSON
        return c.JSON(code, map[string]interface{}{
            "error": err.Error(),
        })
    }
    return nil
})
```

## Performance Tips

1. **Enable Caching in Production**
   ```go
   config.CacheEnabled = true
   ```

2. **Preload Common Templates**
   ```go
   engine.Preload([]string{"layouts/main", "partials/header", "partials/footer"})
   ```

3. **Use Partials Wisely**
   - Keep partials small and focused
   - Cache partial results when possible

4. **Minimize Template Complexity**
   - Move business logic to handlers
   - Keep templates presentational

## Testing

```go
func TestHomeHandler(t *testing.T) {
    app := toutago.New()
    setupTemplates(app, "test")
    
    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    
    c := app.NewContext(req, rec)
    
    err := homeHandler(c)
    require.NoError(t, err)
    
    assert.Equal(t, 200, rec.Code)
    assert.Contains(t, rec.Body.String(), "Welcome")
}
```

## Migration from Other Engines

### From html/template

Fíth syntax is similar but more flexible:

**html/template:**
```html
{{range .Items}}
  {{.Name}}
{{end}}
```

**Fíth:**
```html
{{range .Items}}
  {{.Name}}
{{end}}
```

### From Pongo2

**Pongo2:**
```html
{% extends "base.html" %}
{% block content %}...{% endblock %}
```

**Fíth:**
```html
{{extends "base.html"}}
{{block "content"}}...{{end}}
```

## Complete Example

See [examples/touta-integration/](../examples/touta-integration/) for a complete working example with:

- Full CRUD application
- Authentication
- Form handling
- File uploads
- Error pages
- Multiple layouts

## Resources

- [Fíth Documentation](../docs/)
- [Toutā Documentation](https://touta.io/docs)
- [Example Application](../examples/touta-integration/)
- [API Reference](../docs/api.md)

## Support

For issues or questions:
- GitHub Issues: https://github.com/toutaio/toutago-fith-renderer/issues
- Discussions: https://github.com/toutaio/toutago-fith-renderer/discussions
