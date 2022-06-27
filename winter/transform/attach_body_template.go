package transform

import (
	"fmt"

	"github.com/glacials/twos.dev/winter/document"
)

// AttachBodyTemplate defines the a template called body in d.Template with the
// contents of d.Body, so that it can be included by a page-scoped template like
// document.html.tmpl.
//
// AttachBodyTemplate implements document.Transformation.
func AttachBodyTemplate(d document.Document) (document.Document, error) {
	if _, err := d.Template.New("body").Parse(string(d.Body)); err != nil {
		return document.Document{}, fmt.Errorf("can't parse template (%s", err)
	}

	return d, nil
}
