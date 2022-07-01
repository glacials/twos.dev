package winter

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	styleWrapper = "<span style=\"font-family: sans-serif\">%s</span>"
	toctype      = atom.Ol
)

var (
	replacements = map[string]string{
		// TODO: Do we need this now that we have LaTeX?
		//"–": fmt.Sprintf(styleWrapper, "–"),
		//"—": fmt.Sprintf(styleWrapper, "—"),
	}
	tocheadings = map[atom.Atom]struct{}{atom.H2: {}, atom.H3: {}}
)

type Document struct {
	SourcePath string
	root       *html.Node
	incoming   []*Document
	outgoing   []*Document
	meta       metadata
}

type metadata struct {
	shortname string `yaml:"filename"`
	parent    string `yaml:"parent"`
	kind      kind   `yaml:"type"`
	title     string `yaml:"title"`
	toc       bool   `yaml:"toc"`

	createdAt time.Time `yaml:"date"`
	updatedAt time.Time `yaml:"updated"`
}

type kind int

// IsDraft returns whether the document type is DraftType. This function exists
// to be used by templates.
func (t kind) IsDraft() bool { return t == draft }

// IsPost returns whether the document type is PostType. This function exists to
// be used by templates.
func (t kind) IsPost() bool { return t == post }

// IsPage returns whether the document type is PageType. This function exists to
// be used by templates.
func (t kind) IsPage() bool { return t == page }

// IsGallery returns whether the document type is GalleryType. This function
// exists to be used by templates.
func (t kind) IsGallery() bool { return t == gallery }

const (
	draft kind = iota
	post
	page
	gallery
)

func (k kind) UnmarshalYAML(b []byte) error {
	switch string(b) {
	case "draft", "":
		k = draft
	case "post":
		k = post
	case "page":
		k = page
	case "gallery":
		k = gallery
	default:
		return fmt.Errorf("unknown kind %q", string(b))
	}
	return nil
}

func fromHTML(src string) (*Document, error) {
	var d Document
	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	htm, err := frontmatter.Parse(f, &d.meta)
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(bytes.NewBuffer(htm))
	if err != nil {
		return nil, err
	}

	d.root = root
	d.SourcePath = src
	return &d, nil
}

func fromMarkdown(src string) (*Document, error) {
	var d Document

	f, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	md, err := frontmatter.Parse(f, &d)
	if err != nil {
		return nil, err
	}

	root, err := html.Parse(
		bytes.NewBuffer(
			markdown.ToHTML(
				md,
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
				mdhtml.NewRenderer(mdhtml.RendererOptions{Flags: mdhtml.FlagsNone}),
			),
		),
	)
	if err != nil {
		return nil, err
	}

	d.root = root
	d.SourcePath = src
	return &d, nil
}

func (d *Document) render() ([]byte, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, d.root); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *Document) linksout() (hrfs []string, err error) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.DataAtom == atom.A {
			for _, a := range n.Attr {
				if a.Key == "href" {
					uri, err := url.Parse(a.Val)
					if err != nil {
						return
					}

					if uri.Host == "" {
						hrfs = append(hrfs, strings.TrimSuffix(a.Val, filepath.Ext(a.Val)))
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(d.root)
	return
}

func (d *Document) Title() (string, error) {
	if d.meta.title != "" {
		return d.meta.title, nil
	}

	if h1 := firstOfType(d.root, atom.H1); h1 != nil {
		for child := h1.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				d.meta.title = child.Data
				return child.Data, nil
			}
		}
	}

	return "", fmt.Errorf("h1 in %s nonexistent or nontextual", d.Shortname())
}

func (d *Document) Shortname() string {
	if d.meta.shortname != "" {
		return d.meta.shortname
	}

	n := filepath.Base(d.SourcePath)
	n, _, _ = strings.Cut(n, ".")
	return n
}

func (d *Document) Parent() string {
	if d.meta.parent != "" {
		return d.meta.parent
	}

	p, _, ok := strings.Cut(d.Shortname(), "_")
	if !ok {
		return ""
	}

	return p
}

func (d *Document) Kind() kind {
	return d.meta.kind
}

func (d *Document) fillTOC() error {
	toc := html.Node{
		Attr:     []html.Attribute{{Key: "id", Val: "toc"}},
		Data:     toctype.String(),
		DataAtom: toctype,
		Type:     html.ElementNode,
	}
	var f func(*html.Node) error
	f = func(n *html.Node) error {
		if _, ok := tocheadings[n.DataAtom]; ok && n.Type == html.ElementNode {
			a := html.Node{
				Attr:     []html.Attribute{{Key: atom.Href.String(), Val: "#" + id(n)}},
				Data:     atom.A.String(),
				DataAtom: atom.A,
				Type:     html.ElementNode,
			}
			a.AppendChild(&html.Node{Type: html.TextNode, Data: n.FirstChild.Data})
			li := html.Node{
				Data:     atom.Li.String(),
				DataAtom: atom.Li,
				Type:     html.ElementNode,
			}
			toc.AppendChild(&li)

			totop := html.Node{
				Attr: []html.Attribute{
					{Key: atom.Href.String(), Val: "#toc"},
					{Key: atom.Style.String(), Val: "text-decoration:none"},
				},
				Data:     atom.A.String(),
				DataAtom: atom.A,
				Type:     html.ElementNode,
			}
			totop.AppendChild(&html.Node{Type: html.TextNode, Data: "↑"})

			tohere := html.Node{
				Type:     html.ElementNode,
				Data:     atom.A.String(),
				DataAtom: atom.A,
				Attr: []html.Attribute{
					{Key: atom.Href.String(), Val: "#" + id(n)},
					{Key: atom.Style.String(), Val: "text-decoration:none"},
				},
			}
			tohere.AppendChild(&html.Node{Type: html.TextNode, Data: "#"})

			span := html.Node{
				Type:     html.ElementNode,
				Data:     atom.Span.String(),
				DataAtom: atom.Span,
				Attr: []html.Attribute{
					{Key: atom.Style.String(), Val: "margin-left: 0.5em;"},
				},
			}
			span.AppendChild(&tohere)
			span.AppendChild(&totop)

			n.Parent.InsertBefore(&a, n.NextSibling)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}

		return nil
	}
	if err := f(d.root); err != nil {
		return err
	}

	firstH2 := firstOfType(d.root, atom.H2)
	firstH2.Parent.InsertBefore(&toc, firstH2)
	return nil
}
