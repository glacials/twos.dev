package transform

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/glacials/twos.dev/winter/document"
)

// DiscoverParent determines whether d has a parent document, and attaches it if
// so. A parent document is a lexical calculation; e.g. tatoo_symbols.html has a
// parent of tattoo.html. The parent document is not guaranteed to exist.
//
// DiscoverParent implements document.Transformation.
func DiscoverParent(d document.Document) (document.Document, error) {
	if strings.Contains(d.Stat.Name(), "_") {
		// An underscore means we're in a "sub-page", e.g. tattoo.html links to
		// tattoo_symbols.html.
		re, err := regexp.Compile("(.+)_(.+)\\.(.+)")
		if err != nil {
			return document.Document{}, fmt.Errorf(
				"can't compile base/sub-page regex: %w",
				err,
			)
		}
		matches := re.FindStringSubmatch(d.Stat.Name())
		if len(matches) > 0 {
			d.Parent = fmt.Sprintf("%s.html", matches[1])
		}
	}

	return d, nil
}
