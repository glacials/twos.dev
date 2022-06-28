package transform

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"strings"

	"twos.dev/winter/document"
)

// AttachPartials parses any partial templates the document may need and
// attaches them to the document's template so they can be used.
//
// A partial template is a .html.tmpl file starting with an underscore ("_") and
// is used inside a larger page. For example, a partial may render an image with
// an attached caption.
//
// AttachPartials implements document.Transformation.
func AttachPartials(d document.Document) (document.Document, error) {
	err := LoadPartials(&d.Template)
	if err != nil {
		return document.Document{}, err
	}

	return d, nil
}

func LoadPartials(t *template.Template) error {
	partials, err := filepath.Glob(fmt.Sprintf("src/templates/_*.html.tmpl"))
	if err != nil {
		return fmt.Errorf("can't glob for partials: %w", err)
	}

	for _, partial := range partials {
		name := filepath.Base(partial)
		name = strings.TrimSuffix(name, ".html.tmpl")
		name = strings.TrimPrefix(name, "_")
		p := t.New(name)

		s, err := ioutil.ReadFile(partial)
		if err != nil {
			return fmt.Errorf(
				"can't read partial `%s`: %w",
				partial,
				err,
			)
		}

		if _, err := p.Parse(string(s)); err != nil {
			return fmt.Errorf(
				"can't parse partial `%s`: %w",
				partial,
				err,
			)
		}
	}

	return nil
}
