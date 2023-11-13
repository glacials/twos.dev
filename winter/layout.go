package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"
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
	body []byte
	meta *Metadata
	next Document
}

func NewLayoutDocument(src string, meta *Metadata, next Document) *LayoutDocument {
	return &LayoutDocument{
		meta: meta,
		next: next,
	}
}

func (doc *LayoutDocument) DependsOn(src string) bool {
	return strings.HasPrefix(filepath.Clean(src), "src/templates/")
}

func (doc *LayoutDocument) Load(r io.Reader) error {
	body, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("cannot read %q into layout document: %w", doc.meta.SourcePath, err)
	}
	doc.body = body
	return nil
}

func (doc *LayoutDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *LayoutDocument) Render(w io.Writer) error {
	if doc.meta.Layout == "" {
		_, err := io.Copy(w, bytes.NewBuffer(doc.body))
		if err != nil {
			return fmt.Errorf("cannot render non-layout document %q from layout renderer: %w", doc.meta.SourcePath, err)
		}
		return nil
	}
	layoutBytes, err := os.ReadFile(doc.meta.Layout)
	if err != nil {
		return fmt.Errorf("cannot read %q to execute %q: %w", doc.Metadata().Layout, doc.Metadata().SourcePath, err)
	}
	if len(layoutBytes) == 0 {
		return fmt.Errorf("attempt to render %q using layout %q resulted in 0 bytes rendered", doc.meta.SourcePath, doc.meta.Layout)
	}
	funcs := doc.meta.funcmap()
	tlayout, err := template.New(doc.meta.Layout).Funcs(funcs).Parse(string(layoutBytes))
	if err != nil {
		return fmt.Errorf("cannot read layout %q to execute %q: %w", doc.meta.Layout, doc.meta.SourcePath, err)
	}
	if _, err = tlayout.New("body").Funcs(funcs).Parse(string(doc.body)); err != nil {
		return fmt.Errorf("cannot parse template body %q: %w", doc.meta.SourcePath, err)
	}
	if err := loadDeps(tmplPath, tlayout); err != nil {
		return fmt.Errorf("cannot load template dependencies for %q in layout %q: %w", doc.meta.SourcePath, doc.meta.Layout, err)
	}
	if err := tlayout.Execute(w, doc.meta); err != nil {
		return fmt.Errorf("cannot execute layout document %q with layout %q: %w", doc.meta.SourcePath, doc.meta.Layout, err)
	}
	return nil
}
