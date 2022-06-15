package transform

import (
	"strings"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/glacials/twos.dev/cmd/frontmatter"
)

// DiscoverShortname finds the shortname for the document and stores it in the
// document. The shortname may be specified in frontmatter; if it is not, the
// document's source filename minus extension is used.
//
// DiscoverShortname implements document.Transformation.
func DiscoverShortname(d document.Document) (document.Document, error) {
	matter, _, err := frontmatter.Parse(d.Body)
	if err != nil {
		return document.Document{}, err
	}

	d.Shortname = matter.Shortname
	if d.Shortname == "" {
		parts := strings.Split(d.Stat.Name(), ".")
		d.Shortname = parts[0]
	}

	return d, nil
}
