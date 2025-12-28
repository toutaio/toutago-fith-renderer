# Migration Guide

Guide for migrating from other template engines to Fíth.

## Table of Contents

- [From html/template](#from-htmltemplate)
- [From text/template](#from-texttemplate)
- [From Jinja2](#from-jinja2-python)
- [From Twig](#from-twig-php)
- [From Handlebars](#from-handlebars)

## From html/template

Fíth provides a more intuitive syntax compared to Go's standard `html/template`.

### Syntax Comparison

**Variables:**

```go
// html/template
{{.User.Name}}

// Fíth
{{.User.Name}}  // Same!
```

**Conditionals:**

```go
// html/template
{{if .IsLoggedIn}}
  Welcome!
{{end}}

// Fíth
{{if .IsLoggedIn}}
  Welcome!
{{end}}  // Same!
```

**Loops:**

```go
// html/template
{{range .Items}}
  {{.}}
{{end}}

// Fíth
{{range .Items}}
  {{.}}
{{end}}  // Same!
```

**Functions:**

```go
// html/template
{{printf "%s" .Name}}
{{.Name | upper}}  // If registered

// Fíth
{{upper .Name}}
{{.Name | upper}}  // Both styles work!
```

### Key Differences

#### 1. Function Calls

```go
// html/template - function names first
{{printf "Hello, %s" .Name}}

// Fíth - cleaner syntax
{{upper .Name}}
{{replace .Text "old" "new"}}
```

#### 2. Pipelines

```go
// html/template
{{.Name | printf "%s"}}

// Fíth - more intuitive
{{.Name | upper | trim}}
```

#### 3. Template Composition

```go
// html/template
{{template "header" .}}
{{define "header"}}...{{end}}

// Fíth - cleaner includes/extends
{{include "header"}}
{{extends "layout"}}
{{block "content"}}...{{end}}
```

#### 4. Comments

```go
// html/template
{{/* Comment */}}

// Fíth
{{# Comment #}}
```

### Migration Steps

**1. Update template delimiters:**

Most syntax is compatible, so basic templates work as-is.

**2. Update function calls:**

```go
// Before (html/template)
funcMap := template.FuncMap{
    "upper": strings.ToUpper,
}

// After (Fíth)
renderer := fith.New(fith.Config{
    Functions: map[string]runtime.Function{
        "upper": func(args ...interface{}) (interface{}, error) {
            return strings.ToUpper(args[0].(string)), nil
        },
    },
})
```

**3. Update composition:**

```go
// Before (html/template)
{{template "header" .}}
{{define "header"}}...{{end}}

// After (Fíth)
{{include "header"}}
// header.html is a separate file
```

**4. Update code:**

```go
// Before (html/template)
tmpl := template.Must(template.ParseFiles("template.html"))
tmpl.Execute(w, data)

// After (Fíth)
renderer := fith.New(fith.Config{TemplateDir: "templates"})
output, err := renderer.Render("template", data)
w.Write([]byte(output))
```

---

## From text/template

Very similar to html/template migration. Main difference: Fíth doesn't auto-escape HTML.

If you need HTML escaping, use the `htmlEscape` function:

```
{{htmlEscape .UserInput}}
```

---

## From Jinja2 (Python)

Fíth is inspired by Jinja2, so migration is straightforward.

### Syntax Comparison

**Variables:**

```python
# Jinja2
{{ user.name }}

# Fíth
{{.User.Name}}  // Note: Capital field names (Go convention)
```

**Conditionals:**

```python
# Jinja2
{% if user.is_admin %}
  Admin
{% endif %}

# Fíth
{{if .User.IsAdmin}}
  Admin
{{end}}
```

**Loops:**

```python
# Jinja2
{% for item in items %}
  {{ item }}
{% endfor %}

# Fíth
{{range .Items}}
  {{.}}
{{end}}
```

**Filters:**

```python
# Jinja2
{{ name|upper }}
{{ text|truncate(100) }}

# Fíth
{{.Name | upper}}
{{.Text | truncate 100}}
```

**Includes:**

```python
# Jinja2
{% include 'header.html' %}

# Fíth
{{include "header"}}
```

**Extends:**

```python
# Jinja2
{% extends "layout.html" %}
{% block content %}...{% endblock %}

# Fíth
{{extends "layout"}}
{{block "content"}}...{{end}}
```

### Key Differences

#### 1. Delimiters

```python
# Jinja2
{% ... %}  - Statements
{{ ... }}  - Expressions
{# ... #}  - Comments

# Fíth
{{...}}   - Everything (statements and expressions)
{{# #}}   - Comments
```

#### 2. Field Access

```python
# Jinja2 (Python)
{{ user.name }}      # Lowercase
{{ user['name'] }}   # Dict access

# Fíth (Go)
{{.User.Name}}       # Uppercase (exported fields)
{{.User.Name}}       # Dot notation only
```

#### 3. Loop Variables

```python
# Jinja2
{% for item in items %}
  {{ loop.index }}
  {{ loop.first }}
{% endfor %}

# Fíth
{{range .Items}}
  {{@index}}
  {{@first}}
{{end}}
```

#### 4. Built-in Tests

```python
# Jinja2
{% if value is defined %}
{% if items is empty %}

# Fíth (use functions or checks)
{{if .Value}}
{{range .Items}}...{{else}}Empty{{end}}
```

### Migration Checklist

- [ ] Convert `{% %}` to `{{...}}`
- [ ] Capitalize field names for Go structs
- [ ] Change `loop.index` to `@index`
- [ ] Change `loop.first` to `@first`
- [ ] Update filter syntax from `|filter()` to `| filter`
- [ ] Register custom filters as Fíth functions

---

## From Twig (PHP)

Very similar to Jinja2 migration.

### Syntax Comparison

**Variables:**

```php
{# Twig #}
{{ user.name }}

{# Fíth #}
{{.User.Name}}
```

**Conditionals:**

```php
{# Twig #}
{% if user.isAdmin %}
  Admin
{% endif %}

{# Fíth #}
{{if .User.IsAdmin}}
  Admin
{{end}}
```

**Loops:**

```php
{# Twig #}
{% for item in items %}
  {{ item }}
{% endfor %}

{# Fíth #}
{{range .Items}}
  {{.}}
{{end}}
```

**Filters:**

```php
{# Twig #}
{{ name|upper }}
{{ text|slice(0, 100) }}

{# Fíth #}
{{.Name | upper}}
{{.Text | truncate 100}}
```

### Key Differences

Similar to Jinja2 differences:
- Delimiters: `{{...}}` for everything
- Field names: Capitalized for Go
- Loop variables: `@index` instead of `loop.index`

---

## From Handlebars

Handlebars has a different philosophy (logic-less templates), but migration is possible.

### Syntax Comparison

**Variables:**

```handlebars
{{! Handlebars }}
{{user.name}}

{{! Fíth }}
{{.User.Name}}
```

**Conditionals:**

```handlebars
{{! Handlebars }}
{{#if isAdmin}}
  Admin
{{/if}}

{{! Fíth }}
{{if .IsAdmin}}
  Admin
{{end}}
```

**Loops:**

```handlebars
{{! Handlebars }}
{{#each items}}
  {{this}}
{{/each}}

{{! Fíth }}
{{range .Items}}
  {{.}}
{{end}}
```

**Helpers:**

```handlebars
{{! Handlebars }}
{{toUpperCase name}}

{{! Fíth }}
{{upper .Name}}
```

**Partials:**

```handlebars
{{! Handlebars }}
{{> header}}

{{! Fíth }}
{{include "header"}}
```

### Key Differences

#### 1. Logic

Handlebars is logic-less; Fíth allows richer logic in templates.

```handlebars
{{! Handlebars - needs helper for complex logic }}
{{#if (and isLoggedIn isAdmin)}}

{{! Fíth - nested ifs }}
{{if .IsLoggedIn}}
  {{if .IsAdmin}}
  {{end}}
{{end}}
```

#### 2. Context

```handlebars
{{! Handlebars }}
{{#each items}}
  {{this}}
  {{@index}}
{{/each}}

{{! Fíth }}
{{range .Items}}
  {{.}}
  {{@index}}
{{end}}
```

---

## General Migration Tips

### 1. Start with Simple Templates

Migrate simple templates first to understand syntax differences.

### 2. Test Incrementally

Test each migrated template individually before moving to the next.

### 3. Use Custom Functions

If Fíth is missing a feature, implement it as a custom function:

```go
renderer := fith.New(fith.Config{
    Functions: map[string]runtime.Function{
        "myFilter": func(args ...interface{}) (interface{}, error) {
            // Your logic here
            return result, nil
        },
    },
})
```

### 4. Prepare Data in Go

Move complex logic from templates to your Go code:

```go
// Instead of complex template logic
data := map[string]interface{}{
    "User": user,
    "IsEligible": checkEligibility(user),
    "FormattedDate": formatDate(user.JoinedAt),
}
```

### 5. Organize Templates

```
templates/
  layouts/
    base.html
  pages/
    home.html
  partials/
    header.html
```

### 6. Use Type-Safe Data

Consider using Go structs instead of maps:

```go
type PageData struct {
    Title string
    User  User
    Items []Item
}

data := PageData{
    Title: "Home",
    User:  user,
    Items: items,
}

output, err := renderer.Render("home", data)
```

---

## Common Patterns

### Loading Indicator

```
{{if .IsLoading}}
  <div class="spinner">Loading...</div>
{{else}}
  {{range .Items}}
    <div>{{.}}</div>
  {{else}}
    <p>No items found</p>
  {{end}}
{{end}}
```

### User Authentication

```
{{if .User}}
  Welcome, {{.User.Name | upper}}!
  {{if .User.IsAdmin}}
    <a href="/admin">Admin Panel</a>
  {{end}}
{{else}}
  <a href="/login">Log In</a>
{{end}}
```

### List with Alternating Styles

```
{{range .Items}}
  <div class="item {{if @odd}}odd{{else}}even{{end}}">
    {{.Name}}
  </div>
{{end}}
```

### Breadcrumbs

```
{{range .Breadcrumbs}}
  {{if not @first}} > {{end}}
  <a href="{{.URL}}">{{.Title}}</a>
{{end}}
```

---

## Need Help?

- See [Syntax Reference](syntax.md) for complete syntax
- See [API Reference](api.md) for Go API details
- See [Built-in Functions](functions.md) for available functions
- Check [examples/](../examples/) for working code

If you encounter migration issues, please open an issue on GitHub!
