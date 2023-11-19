package winter

import (
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

var textDocExts = map[string]struct{}{
	".htm":      {},
	".html":     {},
	".md":       {},
	".markdown": {},
	".org":      {},
	".tmpl":     {},
}

// Metadata holds information about a Document that isn't inside the document itself.
type Metadata struct {
	// Category is an optional category for the document. This is used
	// only for a small visual treatment on the index page (if this is
	// of kind post) and on the document page itself.
	//
	// Category MUST be a singular noun that can be pluralized by adding
	// a single "s" at its end, as this is exactly what the visual
	// treatment will do. If this doesn't work for you, go fix that
	// code.
	Category string `yaml:"category,omitempty"`
	// CreatedAt is the time the document was first published.
	CreatedAt time.Time `yaml:"date,omitempty"`
	// Kind specifies the type of document this is.
	// In every user-facing context, this is called "type".
	// In Go we cannot use the "type" keyword, so we use "kind" instead.
	Kind kind `yaml:"type,omitempty"`
	// Layout is the path to the source file for the layout this document should be rendered into.
	//
	// If unset, src/templates/text_document.html.tmpl is used.
	Layout string `yaml:"layout,omitempty"`
	// ParentFilename is the filename component of another document that this one is a child of.
	// Parenthood is a purely semantic relationship for the benefit of the user.
	// Templates can access parents to influence rendering.
	ParentFilename string `yaml:"parent,omitempty"`
	// Preview is a sentence-long blurb of the document,
	// to be shown along with its title as a teaser of its contents.
	Preview string `yaml:"preview,omitempty"`
	// SourcePath is the location on disk of the original file that this document represents.
	// It is relative to the working directory.
	SourcePath string `yaml:"-"`
	// TemplateDir is the location on disk of a directory containing any templates that will be used in the document.
	// By default, it is src/templates.
	TemplateDir string `yaml:"-"`
	// Title is the human-readable title of the document.
	Title string `yaml:"title,omitempty"`
	// TOC is whether a table of contents should be rendered with the
	// document. If true, the table of contents is rendered immediately
	// above the first non-first-level heading.
	TOC bool `yaml:"toc,omitempty"`
	// UpdatedAt is the time the document was last meaningfully updated.
	UpdatedAt time.Time `yaml:"updated,omitempty"`
	// WebPath is the path component of the URL that will point to this document,
	// once rendered.
	// WebPath MUST NOT contain any slashes;
	// everything is top-level.
	//
	// WebPath is equivalent to the path to the destination file,
	// relative to dist.
	WebPath string `yaml:"filename,omitempty"`
}

// NewMetadata returns a Metadata with some defaults filled in
// according to the file at path src.
//
// Defaults that depend on parsing the content of the document,
// such as a Preview generated from its content,
// are not filled in.
func NewMetadata(src, tmplDir string) *Metadata {
	filename := filepath.Base(src)
	i := strings.IndexRune(filename, '.')
	if i < 0 {
		i = len(filename)
	}
	noExt := filename[0:i]
	webPath := noExt
	if _, ok := textDocExts[filepath.Ext(src)]; ok {
		webPath = fmt.Sprintf("%s.html", noExt)
	}
	return &Metadata{
		Kind:        draft,
		Layout:      filepath.Join(tmplDir, "text_document.html.tmpl"),
		SourcePath:  src,
		TemplateDir: tmplDir,
		WebPath:     webPath,
	}
}

func (meta *Metadata) IsType(t string) bool {
	k, err := parseKind(t)
	if err != nil {
		return false
	}
	return k == meta.Kind
}

// funcmap returns a [template.FuncMap] for the document.
// It can be used with [html/template.Template.Funcs].
func (meta *Metadata) funcmap() template.FuncMap {
	now := time.Now()

	return template.FuncMap{
		"add": add,
		"div": div,
		"mul": mul,
		"sub": sub,

		"now": func() time.Time { return now },

		"render": render,
		"yearly": yearly,
	}
}
