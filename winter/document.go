package winter

import (
	"fmt"
	"io"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	draft kind = iota
	post
	page
	static
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

// Document is something that can be built,
// usually from a source file on disk to a destination file on disk.
//
// After a document has been built by calling [Build],
// it can be passed to a template during execution:
//
//	var buf bytes.Buffer
//	t.Execute(&buf, d)
type Document interface {
	// Metadata returns data about the document,
	// which may have been inferred automatically or set by frontmatter.
	Metadata() *Metadata
	// Render generates the final HTML for the document and writes it to w.
	Render(w io.Writer) error
}

// textDocument is a single HTML or Markdown file that will be compiled into a
// static HTML file.
type textDocument struct {
	Metadata

	SrcPath string

	body     []byte
	encoding encoding
	root     *html.Node

	// dependencies is the set of files that building this document depends on,
	// inferred automatically.
	dependencies map[string]struct{}
}

// Metadata holds information about a Document that isn't inside the document itself.
type Metadata struct {
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
	// Filename is the path component of the URL that will point to this document,
	// once rendered.
	// Filename MUST NOT contain any slashes;
	// everything is top-level.
	Filename string `yaml:"filename"`
	// Parent is the filename component of another document that this one is a child of.
	// Parenthood is a purely semantic relationship;
	// no rendering behavior is inherited.
	//
	// The parenthood relationship can be shown in templates from the child using a template function:
	//
	//   {{ parent }}
	//
	// This retrieves the parent document.
	Parent string `yaml:"parent"`
	// Preview is a sentence-long blurb of the document,
	// to be shown along with its title as a teaser of its contents.
	Preview string `yaml:"preview"`
	// SourcePath is the location on disk of the original file that this document represents.
	// It is relative to the working directory.
	SourcePath string `yaml:"-"`
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
		return static, nil
	}
	return -1, fmt.Errorf("unknown kind %q", s)
}

func (d *textDocument) Category() string { return d.Metadata.Category }
func (d *textDocument) Dest() (string, error) {
	if strings.HasSuffix(d.Metadata.Filename, ".html") {
		return d.Metadata.Filename, nil
	}
	return fmt.Sprintf("%s.html", d.Metadata.Filename), nil
}
func (d *textDocument) Layout() string  { return "src/templates/text_document.html.tmpl" }
func (d *textDocument) IsPost() bool    { return d.Kind == post }
func (d *textDocument) IsDraft() bool   { return d.Kind == draft }
func (d *textDocument) Now() time.Time  { return time.Now() }
func (d *textDocument) Preview() string { return d.Metadata.Preview }
func (d *textDocument) Title() string   { return d.Metadata.Title }

// CreatedAt returns the time at which the document was published.
// This is not generated automatically; it is up to the author's discretion.
func (d *textDocument) CreatedAt() time.Time { return d.Metadata.CreatedAt }

// UpdatedAt returns the time at which the document was most recently meaningfully updated.
// This is not generated automatically; it is up to the author's discretion.
func (d *textDocument) UpdatedAt() time.Time { return d.Metadata.UpdatedAt }

type documents []Document

func (d documents) Len() int {
	return len(d)
}

func (d documents) Less(i, j int) bool {
	return d[i].Metadata().CreatedAt.After(d[j].Metadata().CreatedAt)
}

func (d documents) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}
