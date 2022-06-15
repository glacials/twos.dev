package transform

import (
	"bytes"

	"github.com/glacials/twos.dev/cmd/document"
)

// ReplaceSmartQuotes returns a reader identical to the given reader but with
// smart quotes replaced with dumb quotes.
//
// ReplaceSmartQuotes implements document.Transformation.
func ReplaceSmartQuotes(d document.Document) (document.Document, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(d.Body); err != nil {
		return document.Document{}, err
	}

	byts := buf.Bytes()
	byts = bytes.ReplaceAll(byts, []byte("“"), []byte("\""))
	byts = bytes.ReplaceAll(byts, []byte("”"), []byte("\""))
	d.Body = bytes.NewBuffer(byts)

	return d, nil
}
