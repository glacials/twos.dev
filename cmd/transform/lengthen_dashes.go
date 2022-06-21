package transform

import (
	"bytes"
	"fmt"

	"github.com/glacials/twos.dev/cmd/document"
)

const (
	styleWrapper = "<span style=\"font-family: sans-serif\">%s</span>"
)

// LengthenDashes updates dashes to use a variable-width fonts
// en dashes, em dashes, and hyphens look different
// from each other.
//
// LengthenDashes implements document.Transformation.
func LengthenDashes(d document.Document) (document.Document, error) {
	d.Body = bytes.ReplaceAll(
		d.Body,
		[]byte("–"),
		[]byte(fmt.Sprintf(styleWrapper, "–")),
	)
	d.Body = bytes.ReplaceAll(
		d.Body,
		[]byte("—"),
		[]byte(fmt.Sprintf(styleWrapper, "—")),
	)

	return d, nil
}
