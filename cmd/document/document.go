// Package document contains the necessary parsers and renderers to take a
// user-readable source document, such as (optionally templated) Markdown or
// HTML, and render it into a form able to be served statically by a web server
// to the user.
package document

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

const htmlExtension = "html.tmpl"

func stripHTMLExtension(filename string) string {
	filename = strings.TrimSuffix(filename, ".html")
	filename = strings.TrimSuffix(filename, ".html.tmpl")
	return filename
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
	// Body is a reader into the current state of the body of the document.
	Body io.Reader

	// Parent is the name of this document's parent document, if any. For
	// example, tattoo_symbols.html has a parent of tattoo.html. This is a lexical
	// calculation; the parent is not guaranteed to exist.
	Parent string

	// SourceURL is the path to the source file for the document, relative to the
	// repository root.
	SourceURL string

	// Shortname is a human-readable, one-word, all-lowercase tag for the
	// document. The shortname is the part of the filename before the extension. A
	// document's shortname must never change after it is published.
	Shortname string

	// Template is the set of templates this file needs to execute (perhaps more).
	// It starts containing nothing, and must be added to by transformations.
	Template template.Template

	// TemplateVars is the set of variables that will be passed to the
	// template.Execute method when it's eventually called.
	TemplateVars struct {
		Body      template.HTML
		Parent    string
		SourceURL string
		Shortname string
		Title     string

		CreatedAt time.Time
		UpdatedAt time.Time
	}

	// Title is the title of the document that should be displayed in the <title>
	// tag, any relevant <meta> tags, etc. It may be blank, so a fallback should
	// be used when consuming it.
	Title string

	// Stat holds filesystem information about the original source file for this
	// document.
	Stat os.FileInfo

	CreatedAt time.Time
	UpdatedAt time.Time

	transformations []Transformation
}

// Transformation represents a change to be applied to a document, such as
// translating its body from Markdown to HTML, settings its Title field based on
// body contents, or adding a header or footer. The transformation returns the
// updated document. Note that it is not a pointer.
type Transformation func(Document) (Document, error)

// New returns a document with path as its source file, and with transformations
// to be applied in renderers.
func New(path string, transformations []Transformation) (Document, error) {
	f, err := os.Open(path)
	if err != nil {
		return Document{}, err
	}

	stat, err := f.Stat()
	if err != nil {
		return Document{}, err
	}

	essay, err := ioutil.ReadFile("src/templates/essay.html.tmpl")
	if err != nil {
		return Document{}, err
	}

	t, err := template.New("essay").Parse(string(essay))
	if err != nil {
		return Document{}, err
	}

	return Document{
		Body:     f,
		Template: *t,
		SourceURL: fmt.Sprintf(
			"https://github.com/glacials/twos.dev/blob/main/%s",
			path,
		),
		Stat: stat,

		transformations: transformations,
	}, nil
}

// Transform applies all transformations in order to the document body, and
// returns a reader for the result. If any transformation errors, Transform
// stops immediately and returns the error.
//
// The original document is not changed; the transformed copy is returned.
func (d Document) Transform() (Document, error) {
	var err error
	for _, transformation := range d.transformations {
		// The transformation may read data from d.Body, which advances the cursor.
		// If the transformation then replaces d.Body with some transformed data
		// this is fine, but if not it will leave a partially-advanced reader behind
		// for the next transformation. So we'll keep track of any read data and
		// replenish it if the reader is not replaced.
		var readData bytes.Buffer
		tee := io.TeeReader(d.Body, &readData)
		d.Body = tee

		// Warning: do not to replace this = with :=! Otherwise we're not inheriting
		// the transformed d from the previous loop iteration.
		d, err = transformation(d)
		if err != nil {
			return Document{}, fmt.Errorf(
				"can't apply transformation %v: %w",
				runtime.FuncForPC(reflect.ValueOf(transformation).Pointer()).Name(),
				err,
			)
		}

		// Replenish any read data, if the transformation didn't replace the body
		if d.Body == tee {
			d.Body = io.MultiReader(&readData, d.Body)
		}
	}
	return d, nil
}
