package transform

import (
	"bytes"
	"path/filepath"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

// MarkdownToHTML converts the body of the document from Markdown to HTML. Any
// frontmatter must already have been stripped.
//
// MarkdownToHTML implements document.Transformation.
func MarkdownToHTML(d document.Document) (document.Document, error) {
	if filepath.Ext(d.Stat.Name()) != ".md" {
		return d, nil
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(d.Body); err != nil {
		return document.Document{}, err
	}

	d.Body = bytes.NewBuffer(
		markdown.ToHTML(buf.Bytes(), parser.NewWithExtensions(
			parser.Tables|
				parser.FencedCode|
				parser.Autolink|
				parser.Strikethrough|
				parser.Footnotes|
				parser.HeadingIDs|
				parser.Attributes|
				parser.SuperSubscript,
		), nil),
	)

	return d, nil
}
