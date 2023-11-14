package winter

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

// MarkdownDocument represents a source file written in Markdown,
// with optional Go template syntax embedded in it.
//
// MarkdownDocument implements [Document].
//
// The MarkdownDocument is transitory;
// its only purpose is to create a [TemplateDocument].
type MarkdownDocument struct {
	// Next is a pointer to the incarnation of this document that comes after Markdown rendering is complete.
	Next Document
	// SourcePath is the path on disk to the file this Markdown is read from or generated from.
	// The path is relative to the working directory.
	SourcePath string

	deps map[string]struct{}
	meta *Metadata
}

// NewMarkdownDocument creates a new document whose original source is at path src.
//
// Nothing is read from disk; src is metadata.
// To read and parse Markdown, call [Load].
func NewMarkdownDocument(src string, meta *Metadata) *MarkdownDocument {
	return &MarkdownDocument{
		Next:       NewHTMLDocument(src, meta),
		SourcePath: src,

		deps: map[string]struct{}{
			src:                {},
			"public/style.css": {},
		},
		meta: meta,
	}
}

func (doc *MarkdownDocument) DependsOn(src string) bool {
	if _, ok := doc.deps[src]; ok {
		return true
	}
	return false
}

// Load reads Markdown from r and loads it into doc.
//
// If called more than once, the last call wins.
func (doc *MarkdownDocument) Load(r io.Reader) error {
	// Reset metadata to the zero value.
	// Fields removed from frontmatter shouldn't hold onto previous values.
	body, err := frontmatter.Parse(r, doc.meta)
	if err != nil {
		return fmt.Errorf("can't parse %s: %w", doc.SourcePath, err)
	}

	return doc.Next.Load(bytes.NewBuffer(markdown.ToHTML(
		body,
		parser.NewWithExtensions(
			parser.Attributes|
				parser.Autolink|
				parser.FencedCode|
				parser.Footnotes|
				parser.HeadingIDs|
				parser.MathJax|
				parser.Strikethrough|
				parser.Tables,
		),
		newCustomizedRender(),
	)))
}

func (doc *MarkdownDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *MarkdownDocument) Render(w io.Writer) error {
	return doc.Next.Render(w)
}

var markdownExtensions = map[string]struct{}{
	".md":       {},
	".markdown": {},
}

// isMarkdown returns true if and only if src looks like a path to or filename of a Markdown file.
// The file is not touched or inspected;
// the calculation is purely lexicographic.
func isMarkdown(src string) bool {
	for ext := range markdownExtensions {
		if strings.HasSuffix(strings.ToLower(src), ext) {
			return true
		}
	}
	return false
}

// renderImage overrides the standard Markdown-to-HTML renderer.
// It makes unlinked images clickable for a zoomed / gallery view.
func renderImage(w io.Writer, img *ast.Image, entering bool) error {
	if entering {
		if _, err := io.WriteString(
			w,
			fmt.Sprintf(`
				<label class="gallery-item">
			    <input type="checkbox" />
					<img alt="%s" class="thumbnail" src="%s" title="%s" />
				  <img alt="%s" class="fullsize" src="%s" title="%s" />
				`,
				img.Children[0].AsLeaf().Literal,
				img.Destination,
				img.Title,
				img.Children[0].AsLeaf().Literal,
				img.Destination,
				img.Title,
			),
		); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(w, `</label>`); err != nil {
			return err
		}
	}
	return nil
}

func newCustomizedRender() *mdhtml.Renderer {
	insideLink := false
	opts := mdhtml.RendererOptions{
		Flags: mdhtml.FlagsNone,
		RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
			if img, ok := node.(*ast.Image); ok && !insideLink {
				if err := renderImage(w, img, entering); err != nil {
					panic(err)
				}
				// Alt text is a "child" of ast.Image,
				// but we handle it inside the tag in renderImage.
				return ast.SkipChildren, true
			}
			if _, ok := node.(*ast.Link); ok {
				insideLink = entering
			}
			return ast.GoToNext, false
		},
	}
	return mdhtml.NewRenderer(opts)
}
