package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/gomarkdown/markdown"
	mdhtml "github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	draft kind = iota
	post
	page
	gallery
)

const (
	encodingHTML encoding = iota
	encodingMarkdown

	styleWrapper = "<span style=\"font-family: sans-serif\">$0</span>"
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
	hi = map[atom.Atom]int{
		atom.H1: 1,
		atom.H2: 2,
		atom.H3: 3,
		atom.H4: 4,
		atom.H5: 5,
		atom.H6: 6,
	}
	replacements = map[string]string{
		// Break dashes out of monospace homogeneity
		"–": styleWrapper, // En dash
		"—": styleWrapper, // Em dash
		"⁓": styleWrapper, // Swung dash
		"―": styleWrapper, // Horizontal bar
		"⁃": styleWrapper, // Hyphen bullet

		"&#34;": "\"",
		"&#39;": "'",
	}
	tocmin = atom.H2
	tocmax = atom.H3
)

// document is a single HTML or Markdown file that will be compiled into a
// static HTML file.
type document struct {
	SrcPath string

	Kind      kind   `yaml:"type"`
	Parent    string `yaml:"parent"`
	Shortname string `yaml:"filename"`
	Title     string `yaml:"title"`
	TOC       bool   `yaml:"toc"`

	CreatedAt time.Time `yaml:"date"`
	UpdatedAt time.Time `yaml:"updated"`

	body     []byte
	encoding encoding
	root     *html.Node
	incoming []*document
	outgoing []*document
}

type encoding int

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

const ()

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

// NewHTMLDocument creates a new document from the HTML file at the given path.
// High-level information about the document is parsed during this call, such as
// frontmatter and structure. Heavier details like template execution are not
// touched until Render is called.
func NewHTMLDocument(src string) (*document, error) {
	d := document{SrcPath: src}
	if err := d.load(); err != nil {
		return nil, err
	}
	d.encoding = encodingHTML
	if err := d.parseHTML(d.body); err != nil {
		return nil, err
	}
	return &d, d.slurpHTML()
}

// NewMarkdownDocument creates a new document from the Markdown file at the
// given path. High-level information about the document is parsed during this
// call, such as frontmatter and structure. Heavier details like template
// execution are not touched until Render is called.
func NewMarkdownDocument(src string) (*document, error) {
	d := document{SrcPath: src}
	if err := d.load(); err != nil {
		return nil, err
	}
	d.encoding = encodingMarkdown
	if err := d.parseMarkdown(d.body); err != nil {
		return nil, err
	}
	return &d, d.slurpHTML()
}

func (d *document) load() error {
	f, err := os.Open(d.SrcPath)
	if err != nil {
		return err
	}
	defer f.Close()

	body, err := frontmatter.Parse(f, &d)
	if err != nil {
		return fmt.Errorf("can't parse %s: %w", d.SrcPath, err)
	}
	d.body = body
	return nil
}

func (d *document) slurpHTML() error {
	if d.Title == "" {
		if h1 := firstOfType(d.root, atom.H1); h1 != nil {
			for child := h1.FirstChild; child != nil; child = child.NextSibling {
				if child.Type == html.TextNode {
					d.Title = child.Data
				}
			}
			if d.Title == "" {
				return fmt.Errorf("no title found in %s", d.SrcPath)
			}
		}
	}

	if d.Shortname == "" {
		d.Shortname = filepath.Base(d.SrcPath)
	}
	d.Shortname, _, _ = strings.Cut(d.Shortname, ".")

	if d.Parent == "" && d.Kind != gallery {
		if p, _, ok := strings.Cut(d.Shortname, "_"); ok {
			d.Parent = p
		}
	}

	return nil
}

func (d *document) parse() error {
	switch d.encoding {
	case encodingHTML:
		return d.parseHTML(d.body)
	case encodingMarkdown:
		return d.parseMarkdown(d.body)
	default:
		return fmt.Errorf("unknown encoding %d", d.encoding)
	}
}

func (d *document) parseHTML(body []byte) error {
	root, err := html.Parse(bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	d.root = root
	return nil
}

func (d *document) parseMarkdown(body []byte) error {
	root, err := html.Parse(
		bytes.NewBuffer(
			markdown.ToHTML(
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
				mdhtml.NewRenderer(mdhtml.RendererOptions{Flags: mdhtml.FlagsNone}),
			),
		),
	)
	if err != nil {
		return err
	}

	d.root = root
	return nil
}

func (d *document) build() ([]byte, error) {
	if err := d.load(); err != nil {
		return nil, err
	}
	if err := d.slurpHTML(); err != nil {
		return nil, err
	}
	if err := d.parse(); err != nil {
		return nil, err
	}
	if d.TOC {
		if err := d.fillTOC(); err != nil {
			return nil, err
		}
	}
	if err := d.highlightCode(); err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if err := html.Render(&buf, d.root); err != nil {
		return nil, err
	}
	b := buf.Bytes()
	for old, new := range replacements {
		re := regexp.MustCompile(old)
		b = re.ReplaceAll(b, []byte(new))
	}
	return b, nil
}

func (d *document) linksout() (hrfs []string, err error) {
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

// fillTOC iterates over the document looking for headings (<h1>, <h2>, etc.)
// and makes a reflective table of contents.
func (d *document) fillTOC() error {
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
			*grp = append(*grp, tocVars{
				Anchor: attr(n, atom.Id),
				Title:  n.FirstChild.Data,
			})
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
			d.SrcPath,
		)
	}
	firstH2.Parent.InsertBefore(toc, firstH2)
	return nil
}

func (d *document) highlightCode() error {
	for codeBlock := range codeBlocks(d.root) {
		lang := lang(codeBlock)
		formatted, err := syntaxHighlight(lang, codeBlock.FirstChild.Data)
		if err != nil {
			return err
		}

		pre, err := html.Parse(strings.NewReader(formatted))
		if err != nil {
			return err
		}

		originalPre := codeBlock.Parent
		originalPre.Parent.InsertBefore(pre, originalPre)
		originalPre.Parent.RemoveChild(originalPre)
	}

	return nil
}

// codeBlocks returns all code blocks in the document. A code block is defined
// as a <code> tag which is a directy child of a <pre> tag.
func codeBlocks(root *html.Node) map[*html.Node]struct{} {
	blocks := map[*html.Node]struct{}{}
	for _, node := range allOfTypes(root, map[atom.Atom]struct{}{atom.Code: {}}) {
		if node.Parent.DataAtom == atom.Pre {
			blocks[node] = struct{}{}
		}
	}
	return blocks
}

func lang(code *html.Node) string {
	for _, class := range strings.Fields(attr(code, atom.Class)) {
		if _, l, ok := strings.Cut(class, "language-"); ok {
			return l
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

type documents []*document

func (d documents) Len() int {
	return len(d)
}

func (d documents) Less(i, j int) bool {
	// Index must be rendered after others in order for all writing to show on
	// index. TODO: Fix, maybe by having posts() lazily evaluate the rest.
	return d[i].CreatedAt.After(d[j].CreatedAt)
}

func (d documents) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
