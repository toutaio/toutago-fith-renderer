// Package loader provides template loading capabilities for the Fíth template engine.
//
// The loader package implements various strategies for loading templates from different sources:
//   - FileSystemLoader: Load templates from the filesystem
//   - EmbedLoader: Load templates from embedded filesystems (embed.FS)
//   - Template caching for performance
//
// # Basic Usage
//
// Create a filesystem loader and load templates:
//
//	loader := loader.NewFileSystemLoader("templates", []string{".html", ".tpl"})
//	tmpl, err := loader.Load("layouts/main")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Template Resolution
//
// Templates are resolved using slugs (logical names):
//   - "user/profile" → "templates/user/profile.html"
//   - "layouts/main" → "templates/layouts/main.html"
//
// The loader tries each configured extension in order until a file is found.
//
// # Caching
//
// All loaders implement automatic caching:
//   - Templates are parsed once and cached
//   - Subsequent loads return the cached parsed template
//   - Cache can be cleared manually with ClearCache()
//
// # Embed Support
//
// For production deployments, use embedded templates:
//
//	//go:embed templates
//	var templateFS embed.FS
//
//	loader := loader.NewEmbedLoader(templateFS, "templates", []string{".html"})
//
// # Custom Loaders
//
// Implement the Loader interface to create custom loading strategies:
//
//	type Loader interface {
//	    Load(slug string) (*parser.Template, error)
//	    Exists(slug string) bool
//	}
package loader
