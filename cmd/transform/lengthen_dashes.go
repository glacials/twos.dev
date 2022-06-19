package transform

import (
	"bytes"
	"fmt"

	"github.com/glacials/twos.dev/cmd/document"
)

const (
	styleWrapper = "<span style=\"font-family: sans-serif\">%s</span>"
)

// LengthenDashes updates dashes.html to use a variable-width font for dashes
// so that its examples of en dash, em dash, and hyphen look different enough
// from each other.
//
// LengthenDashes implements document.Transformation.
func LengthenDashes(d document.Document) (document.Document, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(d.Body); err != nil {
		return document.Document{}, err
	}

	byts := buf.Bytes()
	byts = bytes.ReplaceAll(
		byts,
		[]byte("–"),
		[]byte(fmt.Sprintf(styleWrapper, "–")),
	)
	byts = bytes.ReplaceAll(
		byts,
		[]byte("—"),
		[]byte(fmt.Sprintf(styleWrapper, "—")),
	)
	d.Body = bytes.NewBuffer(byts)

	return d, nil
}
