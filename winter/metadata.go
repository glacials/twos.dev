package winter

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

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
	Category string `yaml:"category"`
	// Kind specifies the type of document this is. In every user-facing
	// context, this is called "type". In Go we cannot use the "type"
	// keyword, so we use "kind" instead.
	Kind kind `yaml:"type"`
	// Filename is the path component of the URL that will point to this document,
	// once rendered.
	// Filename MUST NOT contain any slashes;
	// everything is top-level.
	Filename string `yaml:"filename"`
	// Layout is the path to the source file for the layout this document should be rendered into.
	//
	// If unset, src/templates/text_document.html.tmpl is used.
	Layout string `yaml:"layout"`
	// Parent is the filename component of another document that this one is a child of.
	// Parenthood is a purely semantic relationship;
	// no rendering behavior is inherited.
	//
	// The parenthood relationship can be shown in templates from the child using a template function:
	//
	//   {{ parent }}
	//
	// This retrieves the parent document.
	Parent string `yaml:"parent"`
	// Preview is a sentence-long blurb of the document,
	// to be shown along with its title as a teaser of its contents.
	Preview string `yaml:"preview"`
	// SourcePath is the location on disk of the original file that this document represents.
	// It is relative to the working directory.
	SourcePath string `yaml:"-"`
	// Title is the human-readable title of the document.
	Title string `yaml:"title"`
	// TOC is whether a table of contents should be rendered with the
	// document. If true, the table of contents is rendered immediately
	// above the first non-first-level heading.
	TOC bool `yaml:"toc"`

	// CreatedAt is the time the document was first published.
	CreatedAt time.Time `yaml:"date"`
	// UpdatedAt is the time the document was last meaningfully updated.
	UpdatedAt time.Time `yaml:"updated"`
}

// NewMetadata returns a Metadata with some defaults filled in
// according to the file at path src.
//
// Defaults that depend on parsing the content of the document,
// such as a Preview generated from its content,
// are not filled in.
func NewMetadata(src string) *Metadata {
	filename := filepath.Base(src)
	i := strings.IndexRune(filename, '.')
	noExt := filename[0:i]
	return &Metadata{
		Filename:   fmt.Sprintf("%s.html", noExt),
		Kind:       draft,
		Layout:     "src/templates/text_document.html.tmpl",
		SourcePath: src,
		Title:      noExt,
	}
}
