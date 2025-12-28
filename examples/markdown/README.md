# Markdown Generation Example

This example demonstrates how to use Fíth to generate Markdown documentation, blog posts, README files, and changelogs.

## Features Demonstrated

- Blog post generation
- API documentation
- README file generation
- Changelog generation
- Code block formatting
- Tables in Markdown
- Dynamic content

## Running the Example

```bash
go run main.go
```

## Key Concepts

### 1. Blog Posts

Generate blog posts with metadata:

```markdown
# {{.Title}}

**Author:** {{.Author}}  
**Date:** {{date "January 2, 2006" .Date}}

{{.Content}}
```

### 2. API Documentation

Create structured API docs with tables:

```markdown
### {{upper .Method}} {{.Path}}

| Name | Type | Required |
|------|------|----------|
{{range .Params}}| {{.Name}} | {{.Type}} | {{if .Required}}Yes{{else}}No{{end}} |
{{end}}
```

### 3. Code Blocks

Use template literals for code fencing:

```go
codeBlock := "```" + "{{.Language}}\n{{.Code}}\n```"
```

### 4. Dynamic Lists

Generate lists from data:

```markdown
{{range .Items}}• {{.}}
{{end}}
```

## Use Cases

- **Documentation Generation**: Auto-generate docs from code comments
- **Blog Systems**: Generate blog posts from structured data
- **README Files**: Create consistent README templates
- **Changelogs**: Auto-generate from version control data
- **Release Notes**: Format release information
- **Static Sites**: Generate Markdown for static site generators

## Best Practices

1. **Use semantic headings** (# for h1, ## for h2, etc.)
2. **Format code blocks** with language identifiers
3. **Keep line lengths** reasonable (80-100 chars)
4. **Use consistent formatting** for lists and tables
5. **Include metadata** (date, author, version)

## Production Examples

### Auto-generate API Docs

```go
type Endpoint struct {
    Method      string
    Path        string
    Description string
    Params      []Parameter
}

// Load endpoints from code
endpoints := analyzeAPICode()

// Generate markdown
output, _ := engine.Render("api-docs", map[string]interface{}{
    "Endpoints": endpoints,
})

os.WriteFile("API.md", []byte(output), 0644)
```

### Generate Changelog from Git

```go
// Parse git log
commits := parseGitLog()

// Group by version
versions := groupByVersion(commits)

// Generate changelog
output, _ := engine.Render("changelog", map[string]interface{}{
    "Versions": versions,
})

os.WriteFile("CHANGELOG.md", []byte(output), 0644)
```

## Tips

- Use **front matter** for metadata (YAML/TOML)
- **Validate links** in generated markdown
- **Test rendering** with your markdown processor
- **Escape special characters** when needed
- **Use templates** for consistent formatting
