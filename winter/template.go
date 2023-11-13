package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sort"
	"text/template/parse"
	"time"

	"twos.dev/winter/graphic"
)

const (
	icnTmplPath = "src/templates/_icon.html.tmpl"
	txtTmplPath = "src/templates/text_document.html.tmpl"
)

// TemplateDocument represents a source file written in Go templates.
// The surrounding syntax can be anything.
//
// TemplateDocument implements [Document].
//
// The TemplateDocument is transitory;
// its only purpose is to create an [HTMLDocument].
type TemplateDocument struct {
	// Next is the HTML document generated from this template document.
	Next   *HTMLDocument
	Parent *TemplateDocument

	deps map[string]struct{}
	meta *Metadata
	// posts is a reference to the substructure's set of posts.
	// It should be populated fully before any call to [TemplateDocument.Load],
	// so that those calls can use the posts function in their [html/template.FuncMap] to discover and list posts.
	posts []Document
}

func NewTemplateDocument(src string, collection []Document) *TemplateDocument {
	m := NewMetadata(src)
	m.WebPath = filepath.Base(src)
	return &TemplateDocument{
		Next: &HTMLDocument{meta: m},

		deps: map[string]struct{}{
			src:                {},
			"public/style.css": {},
		},
		meta:  m,
		posts: collection,
	}
}

func (doc *TemplateDocument) DependsOn(src string) bool {
	if _, ok := doc.deps[src]; ok {
		return true
	}
	return false
}

// Load reads a Go template from r and loads it into doc.
// Any templates referenced within are looked for looked for by name,
// relative to the working directory.
//
// To use a template, treat its filepath as a name:
//
//	{{ template "src/templates/_foo.html.tmpl" }}
//
// Any referenced templates will be loaded as well,
// and attached to the main template.
// This operation is recursive.
//
// If called more than once, the last call wins.
func (doc *TemplateDocument) Load(r io.Reader) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("cannot read template document %q: %w", doc.meta.SourcePath, err)
	}
	if doc.meta.Parent != "" {
		doc.Parent = NewTemplateDocument(doc.meta.Parent, doc.posts)
	}
	funcs, err := doc.funcmap()
	if err != nil {
		return fmt.Errorf("cannot generate funcmap for %q: %w", doc.meta.SourcePath, err)
	}
	tmain, err := template.New(doc.meta.SourcePath).Funcs(funcs).Parse(string(raw))
	if err != nil {
		return fmt.Errorf("cannot parse template document %q: %w", doc.meta.SourcePath, err)
	}
	if doc.meta.Layout != "" {
		tmain, err = template.New("body").Funcs(funcs).Parse(string(raw))
		if err != nil {
			return fmt.Errorf("cannot parse template document %q: %w", doc.meta.SourcePath, err)
		}
		layoutBytes, err := os.ReadFile(doc.meta.Layout)
		if err != nil {
			return fmt.Errorf("cannot read %q to execute %q: %w", doc.Metadata().Layout, doc.Metadata().SourcePath, err)
		}
		tlayout, err := tmain.New(doc.meta.Layout).Funcs(funcs).Parse(string(layoutBytes))
		if err != nil {
			return fmt.Errorf("cannot read layout %q to execute %q: %w", doc.meta.Layout, doc.meta.SourcePath, err)
		}
		tmain = tlayout
	}
	if err := loadDeps(tmain); err != nil {
		return fmt.Errorf("cannot load dependencies for %q: %s", doc.meta.SourcePath, err)
	}
	for _, depTmpl := range tmain.Templates() {
		if depTmpl.Name() != "body" && depTmpl.Name() != doc.meta.SourcePath {
			doc.deps[depTmpl.Name()] = struct{}{}
		}
	}

	var buf bytes.Buffer
	if err := tmain.Execute(&buf, doc.meta); err != nil {
		return fmt.Errorf("cannot execute tmain for %q: %w", doc.meta.SourcePath, err)
	}
	return doc.Next.Load(&buf)
}

func (doc *TemplateDocument) Metadata() *Metadata {
	return doc.meta
}

func (doc *TemplateDocument) Render(w io.Writer) error {
	return doc.Next.Render(w)
}

// funcmap returns a [template.FuncMap] for the document.
// It can be used with [html/template.Template.Funcs].
func (doc *TemplateDocument) funcmap() (template.FuncMap, error) {
	iconFunc, err := icon()
	if err != nil {
		return nil, err
	}
	now := time.Now()

	return template.FuncMap{
		"add": add,
		"div": div,
		"mul": mul,
		"sub": sub,

		"now": func() time.Time { return now },

		"icon":   iconFunc,
		"render": render,
		"parent": func() *TemplateDocument {
			return doc.Parent
		},
		"posts":  doc.postsFunc,
		"yearly": yearly,
	}, nil
}

// postsFunc is a function to be used by templates.
// It retrieves a slice of metadatas for all documents of type post.
func (doc *TemplateDocument) postsFunc() []Document {
	return doc.posts
}

func render(doc Document) (template.HTML, error) {
	var buf bytes.Buffer
	if err := doc.Render(&buf); err != nil {
		return template.HTML(""), err
	}
	return template.HTML(buf.String()), nil
}

// yearly returns the given documents grouped by year.
func yearly(docs []Document) years {
	// groups is a map of year to data for that year.
	groups := map[int]*year{}
	for _, doc := range docs {
		if doc.Metadata().Kind != post {
			continue
		}
		y := doc.Metadata().CreatedAt.Year()
		if _, ok := groups[y]; !ok {
			groups[y] = &year{Documents: documents{}, Year: y}
		}
		groups[y].Documents = append(groups[y].Documents, doc)
	}
	yrs := make(years, 0, len(groups))
	for _, year := range groups {
		sort.Sort(year.Documents)
		yrs = append(yrs, year)
	}
	sort.Sort(yrs)
	return yrs
}

// icon returns a function that can be inserted into a template's FuncMap.
// The returned function renders the image at the given path.
// It always renders it 1em tall.
//
// Its arguments are a path relative to the web root,
// followed by an alt text string.
//
// For example:
//
//	{{ icon "/img/banana.png" "A photo of a banana." }}
//
// When called, the returned function renders the image inline.
func icon() (iconfunc, error) {
	partial, err := os.ReadFile(icnTmplPath)
	if err != nil {
		return nil, err
	}

	t := template.New(icnTmplPath)
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(src graphic.SRC, alt graphic.Alt) (template.HTML, error) {
		v := iconPartialVars{
			Alt: alt,
			SRC: src,
		}

		var buf bytes.Buffer
		if err := t.Execute(&buf, v); err != nil {
			return "", fmt.Errorf("can't execute icon template: %w", err)
		}

		return template.HTML(buf.String()), nil
	}, nil
}

type years []*year

func (a years) Less(i, j int) bool {
	return a[i].Year > a[j].Year
}

func (a years) Len() int {
	return len(a)
}

func (a years) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type year struct {
	Year      int
	Documents documents
}

type iconfunc func(graphic.SRC, graphic.Alt) (template.HTML, error)

type iconPartialVars struct {
	Alt graphic.Alt
	SRC graphic.SRC
}

type tocPartialVars struct {
	Items []tocVars
}

type tocVars struct {
	Anchor string
	Items  []tocVars
	HTML   template.HTML
}

func add(a, b int) int {
	return a + b
}

func div(a, b int) int {
	return a / b
}

func mul(a, b int) int {
	return a * b
}

func sub(a, b int) int {
	return a - b
}

// loadDeps searches a parsed template t and its associated templates
// for references to other templates, e.g.
//
//	{{ template "src/templates/_foo.html.tmpl" }}
//
// It loads any referenced files, parses them into templates, and attaches those templates to t.
// It repeats this recursively until t and all associated templates are fully resolved.
//
// No templates are executed.
func loadDeps(t *template.Template) error {
	for _, tmpl := range t.Templates() {
		for _, node := range tmpl.Tree.Root.Nodes {
			if node.Type() == parse.NodeTemplate {
				name := node.(*parse.TemplateNode).Name
				if name == "body" || t.Lookup(name) != nil {
					continue
				}
				b, err := os.ReadFile(name)
				if err != nil {
					return fmt.Errorf("cannot read template file %q: %w", name, err)
				}
				t2, err := tmpl.New(name).Parse(string(b))
				if err != nil {
					return err
				}
				if err := loadDeps(t2); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
