package transform

import (
	"bytes"
	"html"

	"github.com/glacials/twos.dev/cmd/document"
)

// UnescapeHTML converts any escaped characters back to unescaped.
//
// UnescapeHTML implements document.Transformation.
func UnescapeHTML(d document.Document) (document.Document, error) {
	buf := bytes.NewBuffer(d.Body)
	d.Body = []byte(html.UnescapeString(buf.String()))
	return d, nil
}
