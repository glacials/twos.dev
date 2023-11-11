package winter

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"text/template/parse"
	"time"

	"twos.dev/winter/graphic"
)

const (
	galTmplPath = "src/templates/imgcontainer.html.tmpl"
	icnTmplPath = "src/templates/_icon.html.tmpl"
	txtTmplPath = "src/templates/text_document.html.tmpl"
)

type archivesVars []archiveVars

func (a archivesVars) Less(i, j int) bool {
	return a[i].Year > a[j].Year
}

func (a archivesVars) Len() int {
	return len(a)
}

func (a archivesVars) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type archiveVars struct {
	Year      int
	Documents documents
}

// CommonVars holds several methods accessible to all templates. Note that this
// is a subset of the methods available to any one template both because the
// Substructure will add some additional methods whose implementations don't
// differ by document type (e.g. func Now() time.Time), and because each
// document type can add methods unique to it.
//
// It is likely that CommonVars will be implemented by the same struct that
// implements Document, but it not necessary.
type CommonVars interface {
	IsDraft() bool
	IsPost() bool

	CreatedAt() time.Time
	UpdatedAt() time.Time
}

// commonVars implements CommonVars via Document, plus adds some additional
// methods whose implementations don't differ by document type.
type commonVars struct {
	*substructureDocument
}

// Now returns the current time.
func (v commonVars) Now() time.Time {
	return time.Now()
}

// Parent returns the document's parent, if any. The returned value implements
// Document.
func (v commonVars) Parent() *substructureDocument {
	return v.substructureDocument.Parent
}

type iconfunc func(graphic.SRC, graphic.Alt) (template.HTML, error)

type iconPartialVars struct {
	commonVars
	Alt graphic.Alt
	SRC graphic.SRC
}

type tocPartialVars struct {
	commonVars
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

// icon returns a function that can be inserted into a template's FuncMap.
// The returned function renders the image at the given path.
// It always renders it 1em tall.
//
// Its arguments are a path relative to the web root, followed by an alt text string.
//
// For example:
//
//	{{ icon "/img/banana.png" "A photo of a banana." }}
//
// When called, the returned function renders the image inline.
func (d *substructureDocument) icon() (iconfunc, error) {
	partial, err := os.ReadFile(icnTmplPath)
	if err != nil {
		return nil, err
	}

	t := template.New(icnTmplPath)
	if _, err := t.Parse(string(partial)); err != nil {
		return nil, err
	}

	return func(src graphic.SRC, alt graphic.Alt) (template.HTML, error) {
		// The template called us, so make sure it recompiles if this template changes.
		deps := d.Dependencies()
		deps["icon"] = struct{}{}

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
				if t.Lookup(name) != nil {
					continue
				}
				b, err := os.ReadFile(name)
				if err != nil {
					return fmt.Errorf("cannot read template file: %w", err)
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
