package transform

import (
	"bytes"
	"fmt"

	"github.com/glacials/twos.dev/cmd/document"
)

// AttachBodyTemplate defines the a template called body in d.Template with the
// contents of d.Body, so that it can be included by a page-scoped template like
// essay.html.tmpl.
//
// AttachBodyTemplate implements document.Transformation.
func AttachBodyTemplate(d document.Document) (document.Document, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(d.Body); err != nil {
		return document.Document{}, err
	}

	if _, err := d.Template.New("body").Parse(buf.String()); err != nil {
		return document.Document{}, fmt.Errorf("can't parse template (%s", err)
	}

	return d, nil
}
