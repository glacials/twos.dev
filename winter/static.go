package winter

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
	"time"
)

type staticDocument struct {
	path string
}

// NewStaticDocument returns a document that represents a static asset. The
// asset must be located in the public directory.
func NewStaticDocument(path string) (*staticDocument, error) {
	return &staticDocument{path: path}, nil
}

func (d *staticDocument) Build() ([]byte, error) {
	return ioutil.ReadFile(d.path)
}

func (d *staticDocument) Category() string { return "" }
func (d *staticDocument) Dependencies() map[string]struct{} {
	return map[string]struct{}{}
}

func (d *staticDocument) Dest() (string, error) {
	return filepath.Rel("public", d.path)
}

func (d *staticDocument) Execute(_ io.Writer, _ *template.Template) error {
	return nil
}

func (d *staticDocument) IsDraft() bool {
	return false
}

func (d *staticDocument) IsPost() bool {
	return false
}

func (d *staticDocument) Layout() string {
	return ""
}

func (d *staticDocument) Title() string {
	return ""
}

func (d *staticDocument) CreatedAt() time.Time {
	return time.Time{}
}

func (d *staticDocument) UpdatedAt() time.Time {
	return time.Time{}
}
