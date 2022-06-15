package transform

import (
	"github.com/glacials/twos.dev/cmd/document"
	"github.com/glacials/twos.dev/cmd/frontmatter"
)

// StripFrontmatter removes frontmatter indiscriminately from the document and
// returns the result. This transformation should be applied after all useful
// information has been extracted from the frontmatter, and before the body is
// parsed (e.g. as Markdown or HTML).
//
// StripFrontmatter implements document.Transformation.
func StripFrontmatter(d document.Document) (document.Document, error) {
	if _, body, err := frontmatter.Parse(d.Body); err != nil {
		return document.Document{}, err
	} else {
		d.Body = body
	}

	return d, nil
}
