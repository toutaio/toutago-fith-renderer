# Fíth Template Syntax Reference

Complete reference for the Fíth template language.

## Table of Contents

- [Variables](#variables)
- [Comments](#comments)
- [Control Flow](#control-flow)
- [Functions](#functions)
- [Template Composition](#template-composition)
- [Whitespace Control](#whitespace-control)

## Variables

### Basic Output

Output variables using double curly braces:

```
{{.Name}}
{{.User.Email}}
{{.Items[0]}}
```

### Dot Notation

Access nested fields with dot notation:

```
{{.User.Profile.Name}}
{{.Config.Database.Host}}
```

### Map Access

Access map values:

```
{{.Data.key}}
{{.Settings.theme}}
```

### Array Access

Access array elements (note: not yet implemented):

```
{{.Items[0]}}
{{.Users[1].Name}}
```

## Comments

Add comments that won't appear in output:

```
{{# This is a comment #}}
{{# Comments can span
    multiple lines #}}
```

## Control Flow

### If Statements

Conditional rendering:

```
{{if .User.IsLoggedIn}}
  Welcome back, {{.User.Name}}!
{{end}}
```

### If-Else

```
{{if .User.IsAdmin}}
  <a href="/admin">Admin Panel</a>
{{else}}
  <p>Access denied</p>
{{end}}
```

### If-ElseIf-Else

Multiple conditions:

```
{{if .User.IsAdmin}}
  Admin Dashboard
{{elseif .User.IsModerator}}
  Moderator Panel
{{else}}
  User Profile
{{end}}
```

### Truthiness

The following values are considered false:
- `false` (boolean)
- `0` (number)
- `""` (empty string)
- `nil`
- Empty slices/arrays
- Empty maps

All other values are truthy.

### Range Loops

Iterate over arrays and slices:

```
{{range .Items}}
  <li>{{.}}</li>
{{end}}
```

With index:

```
{{range .Items}}
  <li>{{@index}}: {{.}}</li>
{{end}}
```

Loop over structs or maps:

```
{{range .Users}}
  <div>{{.Name}} ({{.Email}})</div>
{{end}}
```

### Range with Else

Provide fallback for empty collections:

```
{{range .Items}}
  <li>{{.}}</li>
{{else}}
  <p>No items found</p>
{{end}}
```

### Loop Variables

Special variables available in range loops:

- `{{@index}}` - Current iteration index (0-based)
- `{{@first}}` - True if first iteration
- `{{@last}}` - True if last iteration
- `{{@odd}}` - True if odd iteration (1st, 3rd, 5th...)
- `{{@even}}` - True if even iteration (2nd, 4th, 6th...)

Example:

```
{{range .Items}}
  <li class="{{if @odd}}odd{{else}}even{{end}}">
    {{if @first}}<strong>{{end}}
    Item {{@index}}: {{.}}
    {{if @first}}</strong>{{end}}
  </li>
{{end}}
```

## Functions

### Function Calls

Call functions with arguments:

```
{{upper .Name}}
{{truncate .Description 100}}
{{replace .Text "old" "new"}}
```

### Filter Pipeline

Chain functions using the pipe operator:

```
{{.Name | upper}}
{{.Name | upper | trim}}
{{.Text | truncate 100 | upper}}
```

### Mixing Styles

You can mix function call and pipeline styles:

```
{{upper .Name | trim}}
{{.Name | upper | truncate 50}}
```

## Template Composition

### Include

Include another template:

```
{{include "header"}}
{{include "partials/navigation"}}
```

### Include with Context

Pass the entire context:

```
{{include "user-card" .User}}
```

### Include with Parameters

Pass specific parameters (note: parameter syntax may vary):

```
{{include "card" .Post}}
```

### Extends

Inherit from a layout template:

```
{{extends "layouts/base"}}
```

Must be the first directive in the template.

### Blocks

Define or override blocks:

```
{{# In layout: layouts/base.html #}}
<!DOCTYPE html>
<html>
<head>
  <title>{{block "title"}}Default Title{{end}}</title>
</head>
<body>
  <header>{{block "header"}}{{end}}</header>
  <main>{{block "content"}}{{end}}</main>
  <footer>{{block "footer"}}Default Footer{{end}}</footer>
</body>
</html>

{{# In page: pages/home.html #}}
{{extends "layouts/base"}}

{{block "title"}}Home Page{{end}}

{{block "content"}}
  <h1>Welcome Home!</h1>
{{end}}
```

Blocks can have default content that is used if not overridden.

## Whitespace Control

### Default Behavior

By default, whitespace around template tags is preserved:

```
<p>
  {{.Name}}
</p>
```

Output:
```
<p>
  Alice
</p>
```

### Trim Whitespace (Future Feature)

Use `-` to trim whitespace:

```
<p>
  {{- .Name -}}
</p>
```

Output:
```
<p>Alice</p>
```

Note: Whitespace control is planned but not yet implemented.

## Expressions

### Comparisons

Use in `if` statements:

```
{{if .Age}}
  Age is set
{{end}}
```

### Operators (Future Feature)

Planned operators:
- Equality: `==`, `!=`
- Comparison: `<`, `>`, `<=`, `>=`
- Logical: `and`, `or`, `not`

Currently, use functions for complex logic.

## Best Practices

### 1. Use Meaningful Variable Names

```
Good: {{.User.FirstName}}
Bad:  {{.u.fn}}
```

### 2. Keep Logic Simple

Move complex logic to your Go code, not templates.

### 3. Comment Complex Sections

```
{{# User authentication section #}}
{{if .User.IsLoggedIn}}
  ...
{{end}}
```

### 4. Use Includes for Reusable Components

```
{{include "partials/header"}}
{{include "partials/footer"}}
```

### 5. Organize Templates by Purpose

```
templates/
  layouts/
    base.html
    admin.html
  pages/
    home.html
    about.html
  partials/
    header.html
    footer.html
```

## Error Handling

Fíth provides detailed error messages with line and column numbers:

```
Error rendering template "home" at line 12, column 5:
  unknown variable: .User.Namee
  Did you mean: .User.Name?
```

## Limitations

Current limitations (planned for future releases):

1. No arithmetic operators in templates
2. No comparison operators (==, !=, <, >, etc.)
3. No logical operators (and, or, not)
4. No array indexing syntax `[0]`
5. No ternary operator
6. No macro definitions

For these features, prepare data in Go code before passing to templates.

## Examples

See the [examples directory](../examples/) for complete working examples:

- `examples/basic/` - Basic variable output and control flow
- `examples/composition_demo.go` - Includes and layouts
- `examples/functions_demo.go` - Using built-in functions
