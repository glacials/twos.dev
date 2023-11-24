package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/alecthomas/chroma"
	chromahtml "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

var (
	// styleWrapper is used to surround any text that,
	// if the entire page is in a monospace font,
	// deserves to be in a variable-width font.
	//
	// For example, dashes.html shows the difference between the en and em dashes.
	// In a monospace font, they're nearly identical.
	styleWrapper = []byte("<span style=\"font-family: sans-serif\">$0</span>")
	// earlyReplacements are raw text replacements that will happen before HTML is parsed or rendered.
	earlyReplacements = map[string][]byte{
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

	}
	// lateReplacements are raw text replacements that will happen after HTML has been parsed.
	lateReplacements = map[string][]byte{
		"&#34;":  []byte("\""),
		"&#39;":  []byte("'"),
		"&quot;": []byte("\""),
	}
)

// HTMLDocument represents the HTML for a source file.
//
// The document's source content may or may not be in HTML.
// It may be that the source file was written in another language,
// like Markdown, then converted to HTML via [MarkdownDocument].
//
// Therefore, any page generated from a file in src will at some point be an HTMLDocument.
//
// HTMLDocument implements [Document].
type HTMLDocument struct {
	// deps is a set of paths to source files that,
	// when changed,
	// should cause a rebuild of this document.
	deps   map[string]struct{}
	meta   *Metadata
	next   Document
	result []byte
	// root is the topmost HTML tag in the parsed document,
	// usually <html> or its parent.
	root *html.Node
}

// NewHTMLDocument creates a new document whose original source is at path src.
//
// Nothing is read from disk; src is metadata.
// It may or may not point to a file containing HTML.
// To read and parse HTML, call [Load].
func NewHTMLDocument(src string, meta *Metadata, next Document) *HTMLDocument {
	return &HTMLDocument{
		deps: map[string]struct{}{
			src:                {},
			"public/style.css": {},
		},
		meta: meta,
		next: next,
	}
}

func (doc *HTMLDocument) DependsOn(src string) bool {
	if _, ok := doc.deps[src]; ok {
		return true
	}
	if doc.meta.Layout == src {
		return true
	}
	if doc.meta.ParentFilename == src {
		return true
	}
	if strings.HasSuffix(src, ".css") {
		return true
	}
	if strings.HasPrefix(filepath.Clean(src), "src/templates/") {
		return true
	}
	return doc.next.DependsOn(src)
}

// Load reads HTML from r and loads it into doc.
//
// If called more than once, the last call wins.
func (doc *HTMLDocument) Load(r io.Reader) error {
	body, err := frontmatter.Parse(r, doc.meta)
	if err != nil {
		return fmt.Errorf("can't parse %s: %w", doc.meta.SourcePath, err)
	}
	for old, new := range earlyReplacements {
		re, err := regexp.Compile(old)
		if err != nil {
			return err
		}
		body = re.ReplaceAll(body, new)
	}
	root, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return err
	}
	doc.root = root
	if err := doc.Massage(); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := html.Render(&buf, doc.root); err != nil {
		return fmt.Errorf("cannot render HTML to build %q: %w", doc.meta.WebPath, err)
	}
	byts := buf.Bytes()
	for old, new := range lateReplacements {
		re, err := regexp.Compile(old)
		if err != nil {
			return err
		}
		byts = re.ReplaceAll(byts, new)
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

// Massage messes with loaded content to improve the page when it is ultimately rendered.
//
// Massage performs these tasks:
//
//   - Linkifies and stylizes the first <h1> as a page title
//   - Generates a table of contents, if requested by metadata
//   - Sets target=_blank for all <a> tags pointing to external sites
//   - Generates a preview for the document, if one wasn't manually specified
//   - Syntax-highlights code blocks
func (doc *HTMLDocument) Massage() error {
	if err := doc.setTitle(); err != nil {
		return err
	}
	if err := doc.setPreview(); err != nil {
		return err
	}
	doc.setWebPath()
	if err := doc.insertTOC(); err != nil {
		return err
	}
	if err := doc.highlightCode(); err != nil {
		return err
	}
	if err := doc.openExternalLinksInNewTab(); err != nil {
		return err
	}
	if err := doc.replaceSpecialText(); err != nil {
		return err
	}
	return nil
}

func (doc *HTMLDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *HTMLDocument) Post() bool {
	return doc.meta.Kind == post
}

// Render encodes any loaded content into HTML and writes it to w.
func (doc *HTMLDocument) Render(w io.Writer) error {
	if doc.next == nil {
		if _, err := io.Copy(w, bytes.NewReader(doc.result)); err != nil {
			return fmt.Errorf("cannot render HTML: %w", err)
		}
		return nil
	}
	if err := doc.next.Render(w); err != nil {
		return fmt.Errorf("cannot render from %T to %T: %w", doc, doc.next, err)
	}
	return nil
}

// insertTOC creates and inserts a table of contents into the document,
// if one was requested via metadata.
//
// Specifically, it iterates over the document looking for non-first-level headings (<h2>, <h3>, etc.)
// and inserts an ordered hierarchical list of them immediately before the first <h2>.
func (doc *HTMLDocument) insertTOC() error {
	if !doc.Metadata().TOC {
		return nil
	}
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
				Data:     atom.Span.String(),
				DataAtom: atom.Span,
				Type:     html.ElementNode,
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
	if err := f(doc.root); err != nil {
		return fmt.Errorf("cannot recurse into HTML: %w", err)
	}

	tocPath := "_toc.html.tmpl"
	tocbody, err := os.ReadFile(filepath.Join(doc.meta.TemplateDir, tocPath))
	if err != nil {
		return fmt.Errorf("cannot read toc for %q: %w", doc.meta.SourcePath, err)
	}
	toctmpl, err := template.New(tocPath).Parse(string(tocbody))
	if err != nil {
		return fmt.Errorf("cannot parse toc for %q: %w; %s", doc.meta.SourcePath, err, tocbody)
	}
	subtocPath := "_subtoc.html.tmpl"
	subtocbody, err := os.ReadFile(filepath.Join(doc.meta.TemplateDir, subtocPath))
	if err != nil {
		return fmt.Errorf("cannot read subtoc for %q: %w", doc.meta.SourcePath, err)
	}
	if _, err := toctmpl.New(subtocPath).Parse(string(subtocbody)); err != nil {
		return fmt.Errorf("cannot parse subtoc for %q: %w; %s", doc.meta.SourcePath, err, subtocbody)
	}

	var buf bytes.Buffer
	if err := toctmpl.Execute(&buf, v); err != nil {
		return err
	}
	toc, err := html.Parse(&buf)
	if err != nil {
		return err
	}

	// Insert table of contents before first <h2> (i.e. after the introduction).
	firstH2 := firstTag(doc.root, atom.H2)
	if firstH2 == nil {
		return fmt.Errorf(
			"please add at least one H2 heading to %s in order to provide a table of contents anchor point",
			doc.meta.SourcePath,
		)
	}
	firstH2.Parent.InsertBefore(toc, firstH2)
	return nil
}

func (doc *HTMLDocument) replaceSpecialText() error {
	for old, new := range lateReplacements {
		re, err := regexp.Compile(old)
		if err != nil {
			return err
		}
		for _, node := range allOfNodeTypes(doc.root, map[html.NodeType]struct{}{html.TextNode: {}}) {
			node.Data = re.ReplaceAllString(node.Data, string(new))
		}
	}
	return nil
}

// setPreview extracts information needed for processing from the document's HTML.
func (doc *HTMLDocument) setPreview() error {
	if doc.meta.Preview != "" {
		return nil
	}
	if p := firstTag(doc.root, atom.P); p != nil {
		for child := p.FirstChild; child != nil && doc.meta.Preview == ""; child = child.NextSibling {
			if child.Type == html.TextNode {
				doc.meta.Preview = child.Data
			}
		}
	}
	return nil
}

// setTitle finds the first level-1 heading in the document,
// sets the document title to its contents,
// then removes it.
func (doc *HTMLDocument) setTitle() error {
	h1 := firstTag(doc.root, atom.H1)
	if h1 == nil {
		return nil
	}
	for child := h1.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			doc.meta.Title = child.Data
		}
	}
	if doc.meta.Title == "" {
		return fmt.Errorf("no title found in %s", doc.meta.SourcePath)
	}
	h1.Parent.RemoveChild(h1)
	return nil
}

// setWebPath sets a web path for the document if one was not manually specified,
// and sanitizes any existing web path to remove extraneous extensions.
func (doc *HTMLDocument) setWebPath() {
	if doc.meta.WebPath == "" {
		doc.meta.WebPath = filepath.Base(doc.meta.SourcePath)
	}
	shortname, _, _ := strings.Cut(doc.meta.WebPath, ".")
	doc.meta.WebPath = fmt.Sprintf("/%s.html", shortname)
}

// allOfNodeTypes returns all descendant nodes of n with any of the given types.
// The returned slice is sorted in the same way the document was,
// with parent nodes coming before their children.
func allOfNodeTypes(n *html.Node, t map[html.NodeType]struct{}) (m []*html.Node) {
	if _, ok := t[n.Type]; ok {
		m = append(m, n)
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		m = append(m, allOfNodeTypes(child, t)...)
	}
	return
}

// allOfTypes returns all descendant nodes of n with any of the given types. The
// returned slice is sorted in the same way the document was, with parent nodes
// coming before their children.
func allOfTypes(n *html.Node, t map[atom.Atom]struct{}) (m []*html.Node) {
	if _, ok := t[n.DataAtom]; ok {
		m = append(m, n)
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		m = append(m, allOfTypes(child, t)...)
	}
	return
}

// attr returns the attr attribute of the given element node.
func attr(n *html.Node, attr atom.Atom) string {
	for _, a := range n.Attr {
		if a.Key == attr.String() {
			return a.Val
		}
	}
	return ""
}

// clone returns a deep copy of n.
func clone(n *html.Node) (*html.Node, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		return nil, err
	}
	els, err := html.ParseFragment(&buf, n.Parent)
	if err != nil {
		return nil, err
	}

	return els[0], nil
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

// firstTag returns the first and outermost descendant of n with the given tag.
func firstTag(n *html.Node, t atom.Atom) *html.Node {
	if n.Type == html.ElementNode && n.DataAtom == t {
		return n
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if el := firstTag(child, t); el != nil {
			return el
		}
	}

	return nil
}

func (doc *HTMLDocument) highlightCode() error {
	for codeBlock := range codeBlocks(doc.root) {
		lang := lang(codeBlock)
		formatted, err := syntaxHighlight(lang, codeBlock.FirstChild.Data)
		if err != nil {
			return err
		}

		ancestry := doc.root
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

func lang(code *html.Node) string {
	for _, class := range strings.Fields(attr(code, atom.Class)) {
		if _, l, ok := strings.Cut(class, "language-"); ok {
			return l
		}
	}

	return ""
}

func (doc *HTMLDocument) openExternalLinkInNewTab(a *html.Node) error {
	var href *html.Attribute
	for _, attr := range a.Attr {
		switch attr.Key {
		case "target":
			// Don't overwrite explicitly set targets.
			return nil
		case "href":
			href = &attr
		}
	}
	if href == nil {
		// Probably an <a name="..."> element.
		return nil
	}
	url, err := url.Parse(href.Val)
	if err != nil {
		return fmt.Errorf("cannot parse link %q: %w", href.Val, err)
	}
	if url.Host == "" {
		// Don't make internal links open in new tabs.
		return nil
	}
	a.Attr = append(a.Attr, html.Attribute{
		Key: "target",
		Val: "_blank",
	})
	return nil
}

func (doc *HTMLDocument) openExternalLinksInNewTab() error {
	for _, a := range allOfTypes(doc.root, map[atom.Atom]struct{}{atom.A: {}}) {
		if err := doc.openExternalLinkInNewTab(a); err != nil {
			return err
		}
	}
	return nil
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
