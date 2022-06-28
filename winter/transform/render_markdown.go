package transform

import (
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"twos.dev/winter/document"
)

// RenderMarkdown converts the body of the document from Markdown to HTML. Any
// frontmatter must already have been stripped.
//
// RenderMarkdown implements document.Transformation.
func RenderMarkdown(d document.Document) (document.Document, error) {
	if filepath.Ext(d.Stat.Name()) != ".md" {
		return d, nil
	}

	d.Body = markdown.ToHTML(d.Body, parser.NewWithExtensions(
		parser.Tables|
			parser.FencedCode|
			parser.Autolink|
			parser.Strikethrough|
			parser.Footnotes|
			parser.HeadingIDs|
			parser.Footnotes|
			parser.MathJax|
			parser.Attributes,
	), html.NewRenderer(
		html.RendererOptions{Flags: html.FlagsNone},
	),
	)

	return d, nil
}
