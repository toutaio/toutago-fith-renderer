# Built-in Functions Reference

Complete reference for all built-in functions in Fíth.

## Table of Contents

- [String Functions](#string-functions)
- [Array Functions](#array-functions)
- [Logic Functions](#logic-functions)
- [Encoding Functions](#encoding-functions)
- [Date Functions](#date-functions)

## String Functions

### upper

Convert string to uppercase.

**Signature:** `upper(str string) string`

**Example:**
```
{{upper "hello"}}          → HELLO
{{.Name | upper}}          → ALICE
{{upper .Title}}           → MY PAGE TITLE
```

---

### lower

Convert string to lowercase.

**Signature:** `lower(str string) string`

**Example:**
```
{{lower "HELLO"}}          → hello
{{.Name | lower}}          → alice
{{lower .Title}}           → my page title
```

---

### title

Convert string to title case (capitalize first letter of each word).

**Signature:** `title(str string) string`

**Example:**
```
{{title "hello world"}}    → Hello World
{{.Name | title}}          → Alice Johnson
```

**Note:** Uses Unicode title casing rules.

---

### trim

Remove leading and trailing whitespace.

**Signature:** `trim(str string) string`

**Example:**
```
{{trim "  hello  "}}       → hello
{{.Input | trim}}          → (trimmed value)
```

---

### trimPrefix

Remove a prefix from the beginning of a string.

**Signature:** `trimPrefix(str string, prefix string) string`

**Example:**
```
{{trimPrefix "Hello World" "Hello "}}  → World
{{trimPrefix .Path "/api/"}}           → users/123
```

**Note:** If the string doesn't start with the prefix, returns the original string.

---

### trimSuffix

Remove a suffix from the end of a string.

**Signature:** `trimSuffix(str string, suffix string) string`

**Example:**
```
{{trimSuffix "file.txt" ".txt"}}  → file
{{trimSuffix .Filename ".html"}}  → document
```

**Note:** If the string doesn't end with the suffix, returns the original string.

---

### truncate

Truncate string to specified length and add "..." if truncated.

**Signature:** `truncate(str string, maxLen int) string`

**Example:**
```
{{truncate "Hello World" 5}}              → Hello...
{{truncate .Description 100}}             → (truncated to 100 chars)
{{.LongText | truncate 50}}              → (truncated to 50 chars)
```

**Notes:**
- Counts Unicode characters correctly (not bytes)
- Adds "..." if string is truncated
- Returns original string if shorter than maxLen

---

### replace

Replace all occurrences of a substring with another.

**Signature:** `replace(str string, old string, new string) string`

**Example:**
```
{{replace "Hello World" "World" "Universe"}}  → Hello Universe
{{replace .Text "\n" "<br>"}}                 → (newlines to <br>)
{{.Content | replace "  " " "}}               → (double spaces to single)
```

---

## Array Functions

### join

Join array elements into a string with separator.

**Signature:** `join(array []any, separator string) string`

**Example:**
```
{{join .Tags ", "}}                    → tag1, tag2, tag3
{{join .Items " | "}}                  → item1 | item2 | item3
{{.Categories | join " > "}}           → Home > Products > Electronics
```

**Notes:**
- Works with any slice or array type
- Converts elements to strings automatically

---

### len

Get length of string, array, slice, or map.

**Signature:** `len(value any) int`

**Example:**
```
{{len .Items}}              → 5
{{len "hello"}}             → 5 (character count)
{{len .UserMap}}            → 3 (map key count)
```

**Supports:**
- Strings (returns character count, not byte count)
- Arrays and slices
- Maps (returns number of keys)

---

### first

Get first element of an array or slice.

**Signature:** `first(array []any) any`

**Example:**
```
{{first .Items}}            → (first item)
{{first .Users | .Name}}    → (name of first user)
```

**Error:** Returns error if array is empty.

---

### last

Get last element of an array or slice.

**Signature:** `last(array []any) any`

**Example:**
```
{{last .Items}}             → (last item)
{{last .Users | .Name}}     → (name of last user)
```

**Error:** Returns error if array is empty.

---

## Logic Functions

### default

Return default value if input is "falsy".

**Signature:** `default(value any, defaultValue any) any`

**Falsy values:**
- `nil`
- `false`
- `0` (number)
- `""` (empty string)
- Empty slice/array
- Empty map

**Example:**
```
{{default .Name "Anonymous"}}           → Anonymous (if Name is empty)
{{default .Count 0}}                    → 0 (if Count is nil)
{{.Title | default "Untitled"}}         → Untitled (if Title is empty)
```

**Use cases:**
```
<h1>{{default .PageTitle "Home"}}</h1>
<p>Items: {{default .ItemCount 0}}</p>
```

---

## Encoding Functions

### urlEncode

URL-encode a string (percent encoding).

**Signature:** `urlEncode(str string) string`

**Example:**
```
{{urlEncode "hello world"}}              → hello+world
{{urlEncode "user@example.com"}}         → user%40example.com
{{.SearchQuery | urlEncode}}             → (URL-safe query)
```

**Use case:**
```
<a href="/search?q={{urlEncode .Query}}">Search</a>
```

---

### htmlEscape

Escape HTML special characters to prevent XSS.

**Signature:** `htmlEscape(str string) string`

**Example:**
```
{{htmlEscape "<script>alert('xss')</script>"}}
→ &lt;script&gt;alert(&#39;xss&#39;)&lt;/script&gt;

{{htmlEscape .UserInput}}
→ (safely escaped HTML)
```

**Escapes:**
- `<` → `&lt;`
- `>` → `&gt;`
- `&` → `&amp;`
- `'` → `&#39;`
- `"` → `&#34;`

**Use case:**
```
<div class="comment">{{htmlEscape .Comment}}</div>
```

---

## Date Functions

### date

Format a date/time value.

**Signature:** `date(format string, time time.Time|string) string`

**Format strings:**
```
"2006-01-02"              → 2024-12-27
"Jan 2, 2006"             → Dec 27, 2024
"15:04:05"                → 14:30:45
"3:04 PM"                 → 2:30 PM
"Monday, Jan 2"           → Friday, Dec 27
"2006-01-02 15:04:05"     → 2024-12-27 14:30:45
```

**Examples:**
```
{{date "Jan 2, 2006" .CreatedAt}}
→ Dec 27, 2024

{{date "15:04" .Timestamp}}
→ 14:30

{{.PublishedAt | date "Monday, Jan 2, 2006"}}
→ Friday, Dec 27, 2024
```

**Input types:**
- `time.Time` - Go time object
- `string` - Parsed from common formats (RFC3339, "2006-01-02", etc.)

**Go Date Format Reference:**
```
Year:   2006, 06
Month:  Jan, January, 01, 1
Day:    02, 2, _2
Hour:   15 (24h), 03 (12h), 3
Minute: 04, 4
Second: 05, 5
AM/PM:  PM
Day:    Mon, Monday
TZ:     MST, -0700
```

---

## Function Chaining

All functions can be chained using the pipe operator:

```
{{.Name | upper | trim}}
{{.Description | truncate 100 | htmlEscape}}
{{.Tags | join ", " | upper}}
```

## Custom Functions

You can register custom functions when creating the renderer:

```go
renderer := fith.New(fith.Config{
    TemplateDir: "templates",
    Functions: map[string]runtime.Function{
        "reverse": func(args ...interface{}) (interface{}, error) {
            s := args[0].(string)
            runes := []rune(s)
            for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
                runes[i], runes[j] = runes[j], runes[i]
            }
            return string(runes), nil
        },
    },
})
```

Then use in templates:

```
{{reverse "hello"}}  → olleh
```

See [API Reference](api.md#custom-functions) for more details.

## Error Handling

All functions validate their inputs and return descriptive errors:

```
{{upper 123}}
→ Error: upper: argument must be a string

{{truncate "hello"}}
→ Error: truncate: expected 2 arguments, got 1

{{first .EmptyArray}}
→ Error: first: array is empty
```

## Performance Notes

- String functions create new strings (Go strings are immutable)
- Array functions may iterate entire collections
- Use `len` carefully on large collections
- Cache expensive function results in your data preparation

## Coming Soon

Planned functions for future releases:

### String
- `split` - Split string into array
- `contains` - Check if string contains substring
- `startsWith` / `endsWith` - String prefix/suffix checks

### Array
- `slice` - Extract array subset
- `reverse` - Reverse array order
- `sort` - Sort array
- `filter` - Filter array by condition
- `map` - Transform array elements

### Math
- `add`, `sub`, `mul`, `div` - Arithmetic operations
- `round`, `floor`, `ceil` - Rounding
- `min`, `max` - Min/max values

### Logic
- `coalesce` - First non-nil value
- `ternary` - Inline if-then-else

See the [roadmap](../IMPLEMENTATION_PLAN.md) for details.
