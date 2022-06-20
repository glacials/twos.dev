package transform

import (
	"bytes"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/glacials/twos.dev/cmd/frontmatter"
)

// DiscoverDates finds the creation date and optional last-updated date of the
// document, and stores it in d.
//
// DiscoverDates implements document.Transformation.
func DiscoverDates(d document.Document) (document.Document, error) {
	matter, _, err := frontmatter.Parse(bytes.NewBuffer(d.Body))
	if err != nil {
		return document.Document{}, err
	}

	d.CreatedAt = matter.CreatedAt
	if d.CreatedAt.IsZero() {
		// TODO: Dig through Git history
	}

	d.UpdatedAt = matter.UpdatedAt

	return d, nil
}
