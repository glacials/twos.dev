package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
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
	tocEl        = atom.Ol
	toc          = "<ol id=\"toc\">{{.Entries}}</ol>"
	tocEntry     = "<li><a href=\"#{{.Anchor}}\">{{.Section}}</a></li>"
	tocReturn    = `
<span style="margin-left:0.5em">
	<a href="#{{.Anchor}}" style="text-decoration:none">#</a>
	<a href="#toc" style="text-decoration:none">&uarr;</a>
</span>
`
)

var (
	replacements = map[string]string{
		// TODO: Do we need this now that we have LaTeX?
		//"–": fmt.Sprintf(styleWrapper, "–"),
		//"—": fmt.Sprintf(styleWrapper, "—"),
		"&#34;": "\"",
		"&#39;": "'",
	}
	tocmin = atom.H2
	tocmax = atom.H3
	hi     = map[atom.Atom]int{
		atom.H1: 1,
		atom.H2: 2,
		atom.H3: 3,
		atom.H4: 4,
		atom.H5: 5,
		atom.H6: 6,
	}
)

type Document struct {
	SourcePath string
	root       *html.Node
	incoming   []*Document
	outgoing   []*Document

	Kind kind `yaml:"type"`
	TOC  bool `yaml:"toc"`

	CreatedAt time.Time `yaml:"date"`
	UpdatedAt time.Time `yaml:"updated"`

	FrontmatterParent    string `yaml:"parent"`
	FrontmatterShortname string `yaml:"filename"`
	FrontmatterTitle     string `yaml:"title"`
}

type kind int

// IsDraft returns whether the document type is DraftType. This function exists
// to be used by templates.
func (k kind) IsDraft() bool { return k == draft }

// IsPost returns whether the document type is PostType. This function exists to
// be used by templates.
func (k kind) IsPost() bool { return k == post }

// IsPage returns whether the document type is PageType. This function exists to
// be used by templates.
func (k kind) IsPage() bool { return k == page }

// IsGallery returns whether the document type is GalleryType. This function
// exists to be used by templates.
func (k kind) IsGallery() bool { return k == gallery }

const (
	draft kind = iota
	post
	page
	gallery
)

func (k *kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	switch s {
	case "draft", "":
		*k = draft
	case "post":
		*k = post
	case "page":
		*k = page
	case "gallery":
		*k = gallery
	default:
		return fmt.Errorf("unknown kind %q", s)
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

	htm, err := frontmatter.Parse(f, &d)
	if err != nil {
		return nil, err
	}

	d.FrontmatterShortname, _, _ = strings.Cut(d.FrontmatterShortname, ".")

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
	if d.TOC {
		if err := d.fillTOC(); err != nil {
			return nil, err
		}
	}
	var buf bytes.Buffer
	if err := html.Render(&buf, d.root); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	for old, new := range replacements {
		b = bytes.ReplaceAll(b, []byte(old), []byte(new))
	}
	return b, nil
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
	if d.FrontmatterTitle != "" {
		return d.FrontmatterTitle, nil
	}

	if h1 := firstOfType(d.root, atom.H1); h1 != nil {
		for child := h1.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				d.FrontmatterTitle = child.Data
				return child.Data, nil
			}
		}
	}

	return "", fmt.Errorf("h1 in %s nonexistent or nontextual", d.Shortname())
}

func (d *Document) Shortname() string {
	if d.FrontmatterShortname != "" {
		s, _, _ := strings.Cut(d.FrontmatterShortname, ".")
		d.FrontmatterShortname = s
		return d.FrontmatterShortname
	}

	n := filepath.Base(d.SourcePath)
	n, _, _ = strings.Cut(n, ".")
	return n
}

func (d *Document) Parent() string {
	if d.FrontmatterParent != "" {
		return d.FrontmatterParent
	}

	p, _, ok := strings.Cut(d.Shortname(), "_")
	if !ok {
		return ""
	}

	return p
}

// fillTOC iterates over the document looking for headings (<h1>, <h2>, etc.)
// and makes a reflective table of contents.
func (d *Document) fillTOC() error {
	var (
		f func(*html.Node)
		v tocPartialVars
	)
	f = func(n *html.Node) {
		if n.DataAtom >= tocmin && n.DataAtom <= tocmax {
			grp := &v.Items
			for i := hi[tocmin]; i < hi[n.DataAtom] && i < hi[tocmax]; i += 1 {
				if len(*grp) > 0 {
					grp = &((*grp)[len(*grp)-1].Items)
				}
			}
			*grp = append(*grp, tocVars{Anchor: id(n), Title: n.FirstChild.Data})
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(d.root)

	toctmpl, err := ioutil.ReadFile("src/templates/_toc.html.tmpl")
	if err != nil {
		return err
	}
	t, err := template.New("toc").Parse(string(toctmpl))
	if err != nil {
		return err
	}

	subtoctmpl, err := ioutil.ReadFile("src/templates/_subtoc.html.tmpl")
	if err != nil {
		return err
	}
	_, err = t.New("subtoc").Parse(string(subtoctmpl))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, v); err != nil {
		return err
	}
	toc, err := html.Parse(&buf)
	if err != nil {
		return err
	}

	firstH2 := firstOfType(d.root, atom.H2)
	if firstH2 == nil {
		return fmt.Errorf(
			"don't know how to build TOC without any H2 headings in %s",
			d.SourcePath,
		)
	}
	firstH2.Parent.InsertBefore(toc, firstH2)
	return nil
}
