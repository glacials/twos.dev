package transform

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"github.com/alecthomas/chroma/v2"
	htmlformatter "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/glacials/twos.dev/cmd/document"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
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
		), html.NewRenderer(
			html.RendererOptions{
				RenderNodeHook: func(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
					codeBlock, ok := node.(*ast.CodeBlock)
					if !ok {
						return ast.GoToNext, false
					}

					htm, err := renderCodeBlock(codeBlock)
					if err != nil {
						panic(err)
					}

					if _, err := w.Write([]byte(htm)); err != nil {
						panic(err)
					}

					return ast.GoToNext, true
				},
			},
		),
		),
	)

	return d, nil
}
func renderCodeBlock(codeBlock *ast.CodeBlock) (string, error) {
	s := "<code class=\"block"
	if string(codeBlock.Info) != "" {
		s += fmt.Sprintf(" language-%s", codeBlock.Info)
	}

	code, err := syntaxHighlight(
		string(codeBlock.Info),
		string(codeBlock.Literal),
	)
	if err != nil {
		return "", err
	}

	return code, nil
}

func syntaxHighlight(lang, code string) (string, error) {
	// Determine lexer.
	l := lexers.Get(lang)
	if l == nil {
		l = lexers.Analyse(code)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	// Determine formatter.
	f := htmlformatter.New(
		htmlformatter.WithClasses(true),
		htmlformatter.WithLineNumbers(true),
	)

	// Determine style.
	s := styles.Get("dracula")
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := f.Format(&buf, s, it); err != nil {
		return "", err
	}

	return buf.String(), nil
}
