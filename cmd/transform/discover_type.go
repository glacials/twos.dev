package transform

import (
	"bytes"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/glacials/twos.dev/cmd/frontmatter"
)

// DiscoverType finds the type of the document and stores it in d.
//
// DiscoverType implements document.Transformation.
func DiscoverType(d document.Document) (document.Document, error) {
	matter, _, err := frontmatter.Parse(bytes.NewBuffer(d.Body))
	if err != nil {
		return document.Document{}, err
	}

	d.Type = matter.Type

	return d, nil
}
