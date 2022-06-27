package transform

import (
	"bytes"
	"strings"

	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/glacials/twos.dev/winter/document"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// HighlightSyntax hunts for <code> tags and does its best to syntax highlight
// any code inside.
func HighlightSyntax(d document.Document) (document.Document, error) {
	root, err := html.Parse(bytes.NewBuffer(d.Body))
	if err != nil {
		return document.Document{}, err
	}

	for codeBlock := range codeBlocks(root) {
		lang := lang(codeBlock)
		formatted, err := syntaxHighlight(lang, codeBlock.FirstChild.Data)
		if err != nil {
			return document.Document{}, err
		}

		pre, err := html.Parse(strings.NewReader(formatted))
		if err != nil {
			return document.Document{}, err
		}

		originalPre := codeBlock.Parent

		originalPre.Parent.InsertBefore(pre, originalPre)
		originalPre.Parent.RemoveChild(originalPre)
	}

	var buf bytes.Buffer
	if err := html.Render(&buf, root); err != nil {
		return document.Document{}, err
	}

	d.Body = buf.Bytes()

	return d, nil
}

func codeBlocks(root *html.Node) map[*html.Node]struct{} {
	found := map[*html.Node]struct{}{}

	if root.Type == html.ElementNode && root.DataAtom == atom.Code &&
		root.Parent.DataAtom == atom.Pre {
		return map[*html.Node]struct{}{root: {}}
	}

	for c := root.FirstChild; c != nil; c = c.NextSibling {
		for block := range codeBlocks(c) {
			found[block] = struct{}{}
		}
	}

	return found
}

func lang(code *html.Node) string {
	for _, attr := range code.Attr {
		if attr.Key == "class" {
			for _, class := range strings.Fields(attr.Val) {
				if _, l, ok := strings.Cut(class, "language-"); ok {
					return l
				}
			}
		}
	}

	return ""
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

	f := chromahtml.New(
		chromahtml.WithClasses(true),
		chromahtml.WithLineNumbers(true),
	)

	// This has ~no effect because we specify colors in style.css manually, and
	// pass chromahtml.WithClasses(true) above, meaning no inline styles get added
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
