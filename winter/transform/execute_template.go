package transform

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"

	"github.com/glacials/twos.dev/winter/document"
)

type pageVars struct {
	// SourcePath is the path to the source code for the page being rendered,
	// relative to the repository root.
	SourcePath string

	// Parent, if nonempty, is a path to the page that "owns" this one. For
	// example, tattoo.html owns tattoo_symbols.html. (index.html does not own any
	// pages.)
	Parent string
}

type textDocumentVars struct {
	pageVars

	Body      template.HTML
	Title     string
	Shortname string

	CreatedAt string
	UpdatedAt string
}

// ExecuteTemplate executes the text document template with the document's
// template variables and attached templates as input.
//
// ExecuteTemplate implements document.Transformation.
func ExecuteTemplates(d document.Document) (document.Document, error) {
	var buf bytes.Buffer
	if err := d.Template.Execute(&buf, d.TemplateVars); err != nil {
		return document.Document{}, err
	}

	d.Body = buf.Bytes()
	return d, nil
}

func partialPath(name string) string {
	return filepath.Join(
		"src",
		"templates",
		fmt.Sprintf("_%s.html.tmpl", name),
	)
}
