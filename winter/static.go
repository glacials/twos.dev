package winter

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"time"
)

// staticDocument represents a file on disk that will be copied as-is to the web root.
// The subdirectory of the file relative to the web root will match the relative directory of the source file relative to the ./public directory.
type staticDocument struct {
	path string
}

// NewStaticDocument returns a document that represents a static asset. The
// asset must be located in the public directory.
func NewStaticDocument(path string) (*staticDocument, error) {
	return &staticDocument{path: path}, nil
}

// Build returns the raw bytes of the static file.
func (d *staticDocument) Build() ([]byte, error) {
	return os.ReadFile(d.path)
}

// Category returns the empty string.
func (d *staticDocument) Category() string { return "" }

// Dependencies returns an empty set.
func (d *staticDocument) Dependencies() map[string]struct{} {
	return map[string]struct{}{}
}

// Dest returns the final destination for the static document,
// relative to the web root.
func (d *staticDocument) Dest() (string, error) {
	return filepath.Rel("public", d.path)
}

// Execute does nothing.
func (d *staticDocument) Execute(_ io.Writer, _ *template.Template) error {
	return nil
}

// IsDraft returns false.
func (d *staticDocument) IsDraft() bool { return false }

// IsPost returns false.
func (d *staticDocument) IsPost() bool { return false }

// Layout returns the empty string.
func (d *staticDocument) Layout() string { return "" }

// Preview returns the empty string.
func (d *staticDocument) Preview() string { return "" }

// Title returns the empty string.
func (d *staticDocument) Title() string { return "" }

// CreatedAt returns a zero time.
func (d *staticDocument) CreatedAt() time.Time { return time.Time{} }

// UpdatedAt returns a zero time.
func (d *staticDocument) UpdatedAt() time.Time { return time.Time{} }
