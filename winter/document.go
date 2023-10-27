package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
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
	"github.com/gomarkdown/markdown/parser"
	"github.com/niklasfasching/go-org/org"
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
	encodingOrg

	tocEl     = atom.Ol
	toc       = "<ol id=\"toc\">{{.Entries}}</ol>"
	tocEntry  = "<li><a href=\"#{{.Anchor}}\">{{.Section}}</a></li>"
	tocMax    = 5
	tocMin    = 2
	tocReturn = `
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

	styleWrapper = []byte("<span style=\"font-family: sans-serif\">$0</span>")
	replacements = map[string][]byte{
		// Break some special characters out of monospace homogeneity
		"—":       styleWrapper, // Em dash
		"&mdash;": styleWrapper, // Em dash
		"–":       styleWrapper, // En dash
		"&ndash;": styleWrapper, // En dash
		"ƒ":       styleWrapper, // f-stop symbol
		// Figure dash is included but commented to call out that we do NOT want
		// to change its font as it would violate the whole point of the figure dash to make its font
		// different from surrounding digits.
		// "‒":       styleWrapper, // Figure dash
		"―": styleWrapper, // Horizontal bar
		"⁃": styleWrapper, // Hyphen bullet
		"⁓": styleWrapper, // Swung dash

		"&#34;": []byte("\""),
		"&#39;": []byte("'"),
	}
)

type Document interface {
	// Build returns the document HTML as it will be just before being executed as
	// a template.
	Build() ([]byte, error)
	// Category returns an optional category for the document. This is used
	// by templates for styling and display.
	Category() string
	// Dependencies returns a set of filepaths this document depends on.
	// A dependency is defined as a file that, when changed,
	// should cause any browser displaying this document to refresh.
	Dependencies() map[string]struct{}
	// Dest returns the desired final path of the document.
	// The path is relative to the web root and includes the filename.
	Dest() (string, error)
	// Execute executes the given template in the context of the document
	// (i.e. with whatever variables the template needs to execute successfully).
	// It writes the resulting bytes to the given writer.
	//
	// If the document does not use templates,
	// Execute writes the final document bytes to the given writer directly.
	Execute(w io.Writer, t *template.Template) error
	// IsDraft returns whether the document is of type draft.
	IsDraft() bool
	// IsPost returns whether the document is of type post.
	IsPost() bool
	// Layout returns the path to the base template to use for the document,
	// relative to the project root.
	//
	// This will be the template that gets executed to render the document.
	// This is usually not the document itself,
	// but a base template that embeds it while adding headers, footers, etc.
	//
	// The layout template must embed the document template with:
	//
	//     {{ template "body" }}
	//
	// If Layout returns an empty string,
	// the document is a static asset and will be copied without any template execution.
	Layout() string
	Preview() string
	// Title returns the human-readable title of the document.
	Title() string

	// CreatedAt returns the time the document was first published.
	CreatedAt() time.Time
	// UpdatedAt returns the time the document was last meaningfully updated, or a
	// zero time.Time if never.
	UpdatedAt() time.Time
}

// textDocument is a single HTML or Markdown file that will be compiled into a
// static HTML file.
type textDocument struct {
	metadata

	SrcPath string

	body     []byte
	encoding encoding
	root     *html.Node

	// dependencies is the set of files that building this document depends on,
	// inferred automatically.
	dependencies map[string]struct{}
}

type metadata struct {
	// Category is an optional category for the document. This is used
	// only for a small visual treatment on the index page (if this is
	// of kind post) and on the document page itself.
	//
	// Category MUST be a singular noun that can be pluralized by adding
	// a single "s" at its end, as this is exactly what the visual
	// treatment will do. If this doesn't work for you, go fix that
	// code.
	Category string `yaml:"category"`
	// Kind specifies the type of document this is. In every user-facing
	// context, this is called "type". In Go we cannot use the "type"
	// keyword, so we use "kind" instead.
	Kind kind `yaml:"type"`
	// Preview is a sentence-long blurb of the document,
	// to be shown along with its title as a teaser of its contents.
	Preview string `yaml:"preview"`
	// Shortname is the short name of the document. This is used when
	// picking a filename for the document, among other small and
	// mostly-internal bookkeeping measures. It must never change.
	Shortname string `yaml:"filename"`
	// Title is the human-readable title of the document.
	Title string `yaml:"title"`
	// TOC is whether a table of contents should be rendered with the
	// document. If true, the table of contents is rendered immediately
	// above the first non-first-level heading.
	TOC bool `yaml:"toc"`

	// CreatedAt is the time the document was first published.
	CreatedAt time.Time `yaml:"date"`
	// UpdatedAt is the time the document was last meaningfully updated.
	UpdatedAt time.Time `yaml:"updated"`
}

type encoding int

type kind int

func (k *kind) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	var tmp kind

	tmp, err := parseKind(s)
	if err != nil {
		return err
	}

	*k = tmp

	return nil
}

func parseKind(s string) (kind, error) {
	switch s {
	case "draft", "":
		return draft, nil
	case "post":
		return post, nil
	case "page":
		return page, nil
	case "gallery":
		return gallery, nil
	}
	return -1, fmt.Errorf("unknown kind %q", s)
}

// NewHTMLDocument creates a new document from the HTML file at the given path.
// High-level information about the document is parsed during this call, such as
// frontmatter and structure. Heavier details like template execution are not
// touched until Render is called.
func NewHTMLDocument(src string) (*textDocument, error) {
	d := textDocument{
		SrcPath:      src,
		encoding:     encodingHTML,
		dependencies: map[string]struct{}{},
	}
	if err := d.load(); err != nil {
		return nil, err
	}
	return &d, d.slurpHTML()
}

// NewMarkdownDocument creates a new document from the Markdown file at the
// given path. High-level information about the document is parsed during this
// call, such as frontmatter and structure. Heavier details like template
// execution are not touched until Render is called.
func NewMarkdownDocument(src string) (*textDocument, error) {
	d := textDocument{
		SrcPath:      src,
		encoding:     encodingMarkdown,
		dependencies: map[string]struct{}{},
	}
	if err := d.load(); err != nil {
		return nil, err
	}
	return &d, d.slurpHTML()
}

// NewOrgDocument creates a new document from the Org file at the given
// path. High-level information about the document is parsed during this call,
// such as tags and structure. Heavier details like template execution are not
// touched until Render is called.
func NewOrgDocument(src string) (*textDocument, error) {
	d := textDocument{
		SrcPath:      src,
		encoding:     encodingOrg,
		dependencies: map[string]struct{}{},
	}
	if err := d.load(); err != nil {
		return nil, err
	}
	return &d, d.slurpHTML()
}

func (d *textDocument) Dependencies() map[string]struct{} {
	return d.dependencies
}

func (d *textDocument) load() error {
	f, err := os.Open(d.SrcPath)
	if err != nil {
		return fmt.Errorf("cannot open text document source: %w", err)
	}
	defer f.Close()

	body, err := frontmatter.Parse(f, &d.metadata)
	if err != nil {
		return fmt.Errorf("can't parse %s: %w", d.SrcPath, err)
	}
	d.body = body

	switch d.encoding {
	case encodingHTML:
		return d.parseHTML()
	case encodingMarkdown:
		return d.parseMarkdown()
	case encodingOrg:
		return d.parseOrg()
	default:
		return fmt.Errorf("unknown encoding %d", d.encoding)
	}
}

// slurpHTML extracts information needed for processing from the document's HTML.
func (d *textDocument) slurpHTML() error {
	if h1 := firstTag(d.root, atom.H1); h1 != nil {
		for child := h1.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				d.metadata.Title = child.Data
			}
		}
		if d.metadata.Title == "" {
			return fmt.Errorf("no title found in %s", d.SrcPath)
		}
	}

	if d.Shortname == "" {
		d.Shortname = filepath.Base(d.SrcPath)
	}
	d.Shortname, _, _ = strings.Cut(d.Shortname, ".")
	if d.metadata.Preview == "" {
		if p := firstTag(d.root, atom.P); p != nil {
			for child := p.FirstChild; child != nil && d.metadata.Preview == ""; child = child.NextSibling {
				if child.Type == html.TextNode {
					d.metadata.Preview = child.Data
				}
			}
		}
	}

	return nil
}

func (d *textDocument) parseHTML() error {
	root, err := html.Parse(bytes.NewBuffer(d.body))
	if err != nil {
		return err
	}

	d.root = root
	return nil
}

func (d *textDocument) parseMarkdown() error {
	root, err := html.Parse(
		bytes.NewBuffer(
			markdown.ToHTML(
				d.body,
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
			),
		),
	)
	if err != nil {
		return err
	}

	d.root = root
	return nil
}

func (d *textDocument) parseOrg() error {
	orgparser := org.New()
	orgparser.DefaultSettings["OPTIONS"] = strings.Replace(orgparser.DefaultSettings["OPTIONS"], "toc:t", "toc:nil", 1)
	orgdoc := orgparser.Parse(
		bytes.NewBuffer(d.body),
		d.SrcPath,
	)
	orgdoc.BufferSettings["OPTIONS"] = strings.Replace(
		orgdoc.BufferSettings["OPTIONS"],
		"toc:t",
		"toc:nil",
		1,
	)

	orgwriter := org.NewHTMLWriter()
	orgwriter.TopLevelHLevel = 1

	var err error
	for k, v := range orgdoc.BufferSettings {
		switch strings.ToLower(k) {
		case "category":
			d.metadata.Category = v
		case "date":
			d.metadata.CreatedAt, err = time.Parse("2006-01-02", v)
			if err != nil {
				return err
			}
		case "type":
			var err error
			d.metadata.Kind, err = parseKind(v)
			if err != nil {
				return err
			}
		case "filename":
			d.metadata.Shortname = v
		case "title":
			d.metadata.Title = v
		case "toc":
			d.metadata.TOC = (v == "t" || v == "true")
		case "updated":
			d.metadata.UpdatedAt, err = time.Parse("2006-01-02", v)
			if err != nil {
				return err
			}
		}
	}

	htm, err := orgdoc.Write(orgwriter)
	if err != nil {
		return err
	}

	root, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		return err
	}

	d.root = root
	return nil
}

func (d *textDocument) Build() ([]byte, error) {
	if err := d.load(); err != nil {
		return nil, fmt.Errorf("cannot load %q: %w", d.Shortname, err)
	}
	if err := d.slurpHTML(); err != nil {
		return nil, fmt.Errorf("cannot slurp HTML to build %q: %w", d.Shortname, err)
	}
	if d.TOC {
		if err := d.fillTOC(); err != nil {
			return nil, fmt.Errorf("cannot generate table of contents for %q: %w", d.Shortname, err)
		}
	}
	if err := d.highlightCode(); err != nil {
		return nil, fmt.Errorf("cannot highlight %q: %w", d.Shortname, err)
	}
	var buf bytes.Buffer
	if err := html.Render(&buf, d.root); err != nil {
		return nil, fmt.Errorf("cannot render HTML to build %q: %w", d.Shortname, err)
	}
	b := buf.Bytes()
	for old, new := range replacements {
		re, err := regexp.Compile(old)
		if err != nil {
			return nil, err
		}
		b = re.ReplaceAll(b, new)
	}
	return b, nil
}

func (d *textDocument) Category() string { return d.metadata.Category }
func (d *textDocument) Dest() (string, error) {
	return fmt.Sprintf("%s.html", d.metadata.Shortname), nil
}

func (d *textDocument) Execute(w io.Writer, t *template.Template) error {
	return t.Execute(w, d)
}

func (d *textDocument) Layout() string  { return "src/templates/text_document.html.tmpl" }
func (d *textDocument) IsPost() bool    { return d.Kind == post }
func (d *textDocument) IsDraft() bool   { return d.Kind == draft }
func (d *textDocument) Now() time.Time  { return time.Now() }
func (d *textDocument) Preview() string { return d.metadata.Preview }
func (d *textDocument) Title() string   { return d.metadata.Title }

// CreatedAt returns the time at which the document was published.
// This is not generated automatically; it is up to the author's discretion.
func (d *textDocument) CreatedAt() time.Time { return d.metadata.CreatedAt }

// UpdatedAt returns the time at which the document was most recently meaningfully updated.
// This is not generated automatically; it is up to the author's discretion.
func (d *textDocument) UpdatedAt() time.Time { return d.metadata.UpdatedAt }

// fillTOC iterates over the document looking for non-first-level headings (<h2>, <h3>, etc.)
// and inserts a table of contents for them immediately before the first <h2>.
func (d *textDocument) fillTOC() error {
	var (
		f func(*html.Node) error
		v tocPartialVars
	)
	f = func(n *html.Node) error {
		if hi[n.DataAtom] >= tocMin && hi[n.DataAtom] <= tocMax {
			grp := &v.Items
			for i := tocMin; i < hi[n.DataAtom] && i < tocMax; i += 1 {
				if len(*grp) > 0 {
					grp = &((*grp)[len(*grp)-1].Items)
				}
			}
			// Replace the <h*> tag with a <span>
			wrapper := &html.Node{
				Type:     html.ElementNode,
				Data:     atom.Span.String(),
				DataAtom: atom.Span,
			}
			for child := n.FirstChild; child != nil; child = child.NextSibling {
				child, err := clone(child)
				if err != nil {
					return err
				}
				wrapper.AppendChild(child)
			}
			var buf bytes.Buffer
			if err := html.Render(&buf, wrapper); err != nil {
				return err
			}
			*grp = append(*grp, tocVars{
				Anchor: attr(n, atom.Id),
				HTML:   template.HTML(buf.String()),
			})
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if err := f(c); err != nil {
				return err
			}
		}
		return nil
	}
	if err := f(d.root); err != nil {
		return fmt.Errorf("cannot recurse into HTML: %w", err)
	}

	tocbody, err := os.ReadFile("src/templates/_toc.html.tmpl")
	if err != nil {
		return err
	}
	toctmpl, err := template.New("src/templates/_toc.html.tmpl").Parse(string(tocbody))
	if err != nil {
		return err
	}

	subtocbody, err := os.ReadFile("src/templates/_subtoc.html.tmpl")
	if err != nil {
		return fmt.Errorf("cannot read subtoc template: %w", err)
	}
	subtoctmpl, err := toctmpl.New("src/templates/_subtoc.html.tmpl").Parse(string(subtocbody))
	if err != nil {
		return err
	}

	// Ensure updates to some templates cause rebuilds
	d.dependencies[toctmpl.Name()] = struct{}{}
	d.dependencies[subtoctmpl.Name()] = struct{}{}

	var buf bytes.Buffer
	if err := toctmpl.Execute(&buf, v); err != nil {
		return err
	}
	toc, err := html.Parse(&buf)
	if err != nil {
		return err
	}

	// Insert table of contents before first H2 (i.e. after introduction)
	firstH2 := firstTag(d.root, atom.H2)
	if firstH2 == nil {
		return fmt.Errorf(
			"please add at least one H2 heading to %s in order to provide a table of contents anchor point",
			d.SrcPath,
		)
	}
	firstH2.Parent.InsertBefore(toc, firstH2)
	return nil
}

func (d *textDocument) highlightCode() error {
	for codeBlock := range codeBlocks(d.root) {
		lang := lang(codeBlock)
		formatted, err := syntaxHighlight(lang, codeBlock.FirstChild.Data)
		if err != nil {
			return err
		}

		ancestry := d.root
		if ancestry.Type == html.DocumentNode {
			ancestry = ancestry.FirstChild
		}
		pre, err := html.ParseFragment(strings.NewReader(formatted), ancestry)
		if err != nil {
			return fmt.Errorf("can't parse HTML %q: %w", formatted, err)
		}
		originalPre := codeBlock.Parent
		for _, fragment := range pre {
			if fragment.DataAtom == atom.Head {
				continue
			}
			if fragment.DataAtom == atom.Body {
				f := fragment.FirstChild
				f.Parent = nil
				originalPre.Parent.InsertBefore(f, originalPre)
			}
		}
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
	lexer := lexers.Get(lang)
	if lexer == nil {
		lexer = lexers.Analyse(code)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}
	lexer = chroma.Coalesce(lexer)
	formatter := chromahtml.New(
		chromahtml.Standalone(false),
		chromahtml.TabWidth(2),
		chromahtml.WithClasses(true),
		chromahtml.WithLineNumbers(true),
	)

	// This has ~no effect because we specify colors in style.css manually, and
	// pass chromahtml.WithClasses(true) above, meaning no inline styles get added
	s := styles.Get("dracula")
	it, err := lexer.Tokenise(nil, code)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := formatter.Format(&buf, s, it); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type documents []*substructureDocument

func (d documents) Len() int {
	return len(d)
}

func (d documents) Less(i, j int) bool {
	return d[i].CreatedAt().After(d[j].CreatedAt())
}

func (d documents) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
