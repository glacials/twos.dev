package winter

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var mdrepl = map[string][]byte{
	"&quot;": []byte("\""),
}

var (
	templateStart = []byte("{{")
	templateEnd   = []byte("}}")
)

// MarkdownDocument represents a source file written in Markdown,
// with optional Go template syntax embedded in it.
//
// MarkdownDocument implements [Document].
//
// The MarkdownDocument is transitory;
// its only purpose is to create a [TemplateDocument].
type MarkdownDocument struct {
	deps map[string]struct{}
	meta *Metadata
	// next is a pointer to the incarnation of this document that comes after Markdown rendering is complete.
	next   Document
	result []byte
}

type TemplateNode struct {
	ast.Leaf
	Raw []byte
}

// NewMarkdownDocument creates a new document whose original source is at path src.
//
// Nothing is read from disk; src is metadata.
// To read and parse Markdown, call [Load].
func NewMarkdownDocument(src string, meta *Metadata, next Document) *MarkdownDocument {
	return &MarkdownDocument{
		deps: map[string]struct{}{
			src:                {},
			"public/style.css": {},
		},
		meta: meta,
		next: next,
	}
}

func (doc *MarkdownDocument) DependsOn(src string) bool {
	if _, ok := doc.deps[src]; ok {
		return true
	}
	if doc.meta.WebPath == "/archives.html" || doc.meta.WebPath == "/writing.html" || doc.meta.WebPath == "/index.html" {
		return true
	}
	if strings.HasPrefix(filepath.Clean(src), "src/templates/") {
		return true
	}
	return doc.next.DependsOn(src)
}

// Load reads Markdown from r and loads it into doc.
//
// If called more than once, the last call wins.
func (doc *MarkdownDocument) Load(r io.Reader) error {
	// Reset metadata to the zero value.
	// Fields removed from frontmatter shouldn't hold onto previous values.
	body, err := frontmatter.Parse(r, doc.meta)
	if err != nil {
		return fmt.Errorf("can't parse %s: %w", doc.meta.SourcePath, err)
	}
	p := parser.NewWithExtensions(
		parser.Attributes |
			parser.Autolink |
			parser.FencedCode |
			parser.Footnotes |
			parser.HeadingIDs |
			parser.MathJax |
			parser.Strikethrough |
			parser.Tables,
	)
	p.Opts.ParserHook = parserHook

	byts := markdown.ToHTML(body, p, newRenderer())
	for old, new := range mdrepl {
		byts = bytes.ReplaceAll(byts, []byte(old), new)
	}
	doc.result = byts
	if doc.next == nil {
		return nil
	}
	if err := doc.next.Load(bytes.NewReader(doc.result)); err != nil {
		return fmt.Errorf("cannot load from %T to %T: %w", doc, doc.next, err)
	}
	return nil
}

func (doc *MarkdownDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *MarkdownDocument) Render(w io.Writer) error {
	if doc.next == nil {
		if _, err := io.Copy(w, bytes.NewReader(doc.result)); err != nil {
			return fmt.Errorf("cannot render Markdown: %w", err)
		}
		return nil
	}
	if err := doc.next.Render(w); err != nil {
		return fmt.Errorf("cannot render from %T to %T: %w", doc, doc.next, err)
	}
	return nil
}

func newRenderer() *mdhtml.Renderer {
	opts := mdhtml.RendererOptions{
		Flags: mdhtml.FlagsNone,
	}
	return mdhtml.NewRenderer(opts)
}

func parserHook(data []byte) (ast.Node, []byte, int) {
	if !bytes.HasPrefix(data, templateStart) {
		return nil, nil, 0
	}
	start := bytes.Index(data, templateStart)
	if start < 0 {
		return nil, nil, 0
	}
	end := bytes.Index(data, templateEnd)
	if end < 0 {
		return nil, data, 0
	}
	return &ast.Text{Leaf: ast.Leaf{Literal: data[0 : end+len(templateEnd)]}}, nil, end + len(templateEnd)
}
