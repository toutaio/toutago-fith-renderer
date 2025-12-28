# CSV Generation Example

This example demonstrates how to use Fíth to generate CSV files for data export and reporting.

## Features Demonstrated

- Basic CSV export
- Date formatting in CSV
- Conditional columns
- CSV value escaping
- Sales reports
- Custom CSV functions

## Running the Example

```bash
go run main.go
```

## Output

The example generates four CSV files:

1. **User Export** - Simple contact list
2. **Sales Report** - Transaction data with dates
3. **Conditional Columns** - Dynamic column selection
4. **Escaped Values** - Proper handling of commas and quotes

## Key Concepts

### 1. Basic CSV Template

```
Name,Email,Age
{{range .Users}}{{.Name}},{{.Email}},{{.Age}}
{{end}}
```

### 2. Date Formatting

Use the `date` function to format dates:

```
{{date "YYYY-MM-DD" .Date}}
```

### 3. Conditional Columns

Include columns based on flags:

```
Name{{if .IncludeEmail}},Email{{end}}
{{range .Items}}{{.Name}}{{if $.IncludeEmail}},{{.Email}}{{end}}
{{end}}
```

### 4. CSV Escaping

Register a custom function to escape CSV values:

```go
engine.RegisterFunction("csvEscape", func(args ...interface{}) (interface{}, error) {
    str := fmt.Sprintf("%v", args[0])
    // Escape logic here
    return escaped, nil
})
```

Use it in templates:

```
{{csvEscape .Description}}
```

## Use Cases

- Data export from databases
- Report generation
- Analytics data extraction
- Bulk data import preparation
- Log file generation
- Configuration exports

## Best Practices

1. **Always escape values** that might contain commas, quotes, or newlines
2. **Use consistent date formats** across your CSV files
3. **Include headers** for clarity
4. **Handle null/empty values** gracefully
5. **Consider using BOM** for Excel compatibility (UTF-8 BOM: `\uFEFF`)

## Production Usage

```go
// Generate CSV and write to file
output, err := engine.RenderString(csvTemplate, data)
if err != nil {
    log.Fatal(err)
}

err = os.WriteFile("export.csv", []byte(output), 0644)
if err != nil {
    log.Fatal(err)
}
```

## Advanced: Streaming Large CSV Files

For large datasets, consider streaming:

```go
// Use RenderBytes for efficiency
bytes, err := engine.RenderBytes(csvTemplate, data)
if err != nil {
    log.Fatal(err)
}

// Write directly to response or file
w.Write(bytes)
```
