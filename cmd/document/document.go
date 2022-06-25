// Package document contains the necessary parsers and renderers to take a
// user-readable source document, such as (optionally templated) Markdown or
// HTML, and render it into a form able to be served statically by a web server
// to the user.
package document

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/glacials/twos.dev/cmd/frontmatter"
)

const htmlExtension = "html.tmpl"

func stripHTMLExtension(filename string) string {
	filename = strings.TrimSuffix(filename, ".html")
	filename = strings.TrimSuffix(filename, ".html.tmpl")
	return filename
}

// BaseVars is the set of variables present when executing the template for any
// document.
type BaseVars struct {
	Parent     string
	Shortname  string
	SourcePath string
	Type       frontmatter.Type

	Now time.Time
}

// TextVars is a set of variables present when executing text-based documents
// like PageType, PostType, or DraftType.
type TextVars struct {
	BaseVars

	Body       template.HTML
	Parent     string
	SourcePath string
	Title      string

	CreatedAt time.Time
	UpdatedAt time.Time
}

// Document is a file that is sent down a pipeline of transformations in order,
// starting from reading a file and ending in writing a static .html file ready
// for distribution. Each transformation is a function that implements
// trasformation.Transformation.
//
// For example, a document that starts as a .md file may first be transformed
// (read: rendered) into a .html file, then by setting its title to the first
// heading in its contents, then by executing it as a template.
type Document struct {
	frontmatter.Matter

	// Body is the current state of the body of the document.
	Body []byte

	// Parent is the name of this document's parent document, if any. For
	// example, tattoo_symbols.html has a parent of tattoo.html. This is a lexical
	// calculation; the parent is not guaranteed to exist.
	Parent string

	// SourcePath is the path to the source file for the document, relative to the
	// repository root.
	SourcePath string

	Shortname string

	// Template is the set of templates this file needs to execute (perhaps more).
	// It starts containing nothing, and must be added to by transformations.
	Template template.Template

	// TemplateVars is the set of variables that will be passed to the
	// template.Execute method when it's eventually called.
	TemplateVars TextVars

	// Title is the title of the document that should be displayed in the <title>
	// tag, any relevant <meta> tags, etc. It may be blank, so a fallback should
	// be used when consuming it.
	Title string

	// Stat holds filesystem information about the original source file for this
	// document.
	Stat os.FileInfo

	transformations []Transformation
	debug           bool
}

// Transformation represents a change to be applied to a document, such as
// translating its body from Markdown to HTML, settings its Title field based on
// body contents, or adding a header or footer. The transformation returns the
// updated document. Note that it is not a pointer.
type Transformation func(Document) (Document, error)

// New returns a document with path as its source file, and with transformations
// to be applied in renderers.
func New(path string, trs []Transformation, debug bool) (Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return Document{}, err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return Document{}, err
	}

	stat, err := f.Stat()
	if err != nil {
		return Document{}, err
	}

	text_document, err := ioutil.ReadFile("src/templates/text_document.html.tmpl")
	if err != nil {
		return Document{}, err
	}

	t, err := template.New("text_document").Parse(string(text_document))
	if err != nil {
		return Document{}, err
	}

	return Document{
		Body:       body,
		Template:   *t,
		SourcePath: path,
		Stat:       stat,

		transformations: trs,
		debug:           debug,
	}, nil
}

// Transform applies all transformations in order to the document body, and
// returns a reader for the result. If any transformation errors, Transform
// stops immediately and returns the error.
//
// The original document is not changed; the transformed copy is returned.
func (d Document) Transform() (Document, error) {
	var err error
	for tindex, transformation := range d.transformations {
		tname := runtime.FuncForPC(reflect.ValueOf(transformation).Pointer()).
			Name()
		_, tshortname, ok := strings.Cut(tname, "transform.")
		if !ok {
			return Document{}, fmt.Errorf(
				"unexpected transformation package name in %s",
				tname,
			)
		}

		// Warning: do not to replace this = with :=! Otherwise we're not inheriting
		// the transformed d from the previous loop iteration.
		d, err = transformation(d)
		if err != nil {
			return Document{}, fmt.Errorf(
				"can't apply transformation %v: %w",
				tshortname,
				err,
			)
		}

		if d.debug {
			path := filepath.Join(
				"dist",
				"debug",
				d.Stat.Name(),
				fmt.Sprintf("%02d_%s.html", tindex, tshortname),
			)
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return Document{}, err
			}

			f, err := os.Create(path)
			if err != nil {
				return Document{}, err
			}
			defer f.Close()

			if err := ioutil.WriteFile(path, d.Body, 0644); err != nil {
				return Document{}, err
			}
		}

	}
	return d, nil
}
