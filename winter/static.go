package winter

import (
	"fmt"
	"io"
)

// StaticDocument represents a file on disk that will be copied as-is to the web root.
// The subdirectory of the file relative to the web root will match the relative directory of the source file relative to the ./public directory.
type StaticDocument struct {
	SourcePath string

	deps map[string]struct{}
	meta *Metadata
	r    io.Reader
}

// NewStaticDocument creates a new document whose original source is at path src,
// and whose desired web path is webPath.
//
// Nothing is read from disk; src is metadata.
// To read the static file, call [Load].
func NewStaticDocument(src, webPath string) *StaticDocument {
	m := NewMetadata(src, tmplPath)
	m.WebPath = webPath
	return &StaticDocument{
		SourcePath: src,

		deps: map[string]struct{}{
			src: {},
		},
		meta: m,
	}
}

func (doc *StaticDocument) DependsOn(src string) bool {
	if _, ok := doc.deps[src]; ok {
		return true
	}
	return false
}

func (doc *StaticDocument) Load(r io.Reader) error {
	doc.r = r
	return nil
}

func (doc *StaticDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *StaticDocument) Render(w io.Writer) error {
	if _, err := io.Copy(w, doc.r); err != nil {
		return fmt.Errorf("cannot copy static file %q for render: %w", doc.SourcePath, err)
	}
	return nil
}
