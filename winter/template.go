package winter

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
)

type templateVars struct {
	*Document
	*substructure
}

func loadTemplates(t *template.Template) error {
	partials, err := filepath.Glob("src/templates/*.html.tmpl")
	if err != nil {
		return fmt.Errorf("can't glob for partials: %w", err)
	}

	for _, partial := range partials {
		name := filepath.Base(partial)
		name = strings.TrimPrefix(name, "_")
		name, _, _ = strings.Cut(name, ".") // Trim extensions, even e.g. .html.tmpl

		p, err := ioutil.ReadFile(partial)
		if err != nil {
			return fmt.Errorf(
				"can't read partial `%s`: %w",
				partial,
				err,
			)
		}

		if _, err := t.New(name).Parse(string(p)); err != nil {
			return fmt.Errorf(
				"can't parse partial `%s`: %w",
				partial,
				err,
			)
		}
	}

	return nil
}
