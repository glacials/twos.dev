package transform

import (
	"bytes"
	"strings"
	"time"

	"twos.dev/winter/document"
	"twos.dev/winter/frontmatter"
)

// ExtractFrontmatter slurps the document's frontmatter and saves it to the
// Document. The frontmatter is permanently removed from the document body.
//
// ExtractFrontmatter implements document.Transformation.
func ExtractFrontmatter(d document.Document) (document.Document, error) {
	if matter, body, err := frontmatter.Parse(bytes.NewBuffer(d.Body)); err != nil {
		return document.Document{}, err
	} else {
		d.Matter = matter
		d.Shortname = matterToShortname(matter, d.Stat.Name())

		d.Body = body
	}

	return d, nil
}

func matterToShortname(matter frontmatter.Matter, filename string) string {
	if matter.Shortname == "" {
		parts := strings.Split(filename, ".")
		return parts[0]
	}

	return matter.Shortname
}

func matterToDates(matter frontmatter.Matter) (time.Time, time.Time) {
	if matter.CreatedAt.IsZero() {
		// TODO: Dig through Git history
	}

	return matter.CreatedAt, matter.UpdatedAt
}
