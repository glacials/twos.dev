package winter

import (
	"fmt"
	"html/template"
	"io"
	"os"
)

// LayoutDocument represents a document to be assembled by being placed inside a layout-style template.
//
// For example, a blog post document would not itself contain a site header or footer;
// instead, it would be fully rendered then placed inside a layout which includes it:
//
//	{{ template "body" }}
//
// The LayoutDocument's job is to facilitate that embedding.
// It will usually come last in the load/render chain.
type LayoutDocument struct {
	body io.Reader
	meta *Metadata
}

func NewLayoutDocument(src string, meta *Metadata) *LayoutDocument {
	return &LayoutDocument{
		meta: meta,
	}
}

func (doc *LayoutDocument) Load(r io.Reader) error {
	doc.body = r
	return nil
}

func (doc *LayoutDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *LayoutDocument) Render(w io.Writer) error {
	docBytes, err := io.ReadAll(doc.body)
	if err != nil {
		return fmt.Errorf("cannot read body for %q: %w", doc.meta.SourcePath, err)
	}
	layoutBytes, err := os.ReadFile(doc.meta.Layout)
	if err != nil {
		return fmt.Errorf("cannot read %q to execute %q: %w", doc.Metadata().Layout, doc.Metadata().SourcePath, err)
	}
	funcs := doc.meta.funcmap()
	tlayout, err := template.New(doc.meta.Layout).Funcs(funcs).Parse(string(layoutBytes))
	if err != nil {
		return fmt.Errorf("cannot read layout %q to execute %q: %w", doc.meta.Layout, doc.meta.SourcePath, err)
	}
	_, err = tlayout.New("body").Funcs(funcs).Parse(string(docBytes))
	if err != nil {
		return fmt.Errorf("cannot parse template body %q: %w", doc.meta.SourcePath, err)
	}
	return nil
}
