package transform

import (
	"bytes"

	"github.com/glacials/twos.dev/cmd/document"
)

var smartquotes = map[string]string{
	"“":       "\"",
	"”":       "\"",
	"&ldquo;": "\"",
	"&rdquo;": "\"",
	"&quot;":  "\"",
	"&#34;":   "\"",
}

// ReplaceSmartQuotes returns a reader identical to the given reader but with
// smart quotes et al. replaced with dumb quotes. This must happen after
// Markdown parsing, as go-markdown replaces dumb quotes with smart quotes even
// when used to indicate template string literals.
//
// ReplaceSmartQuotes implements document.Transformation.
func ReplaceSmartQuotes(d document.Document) (document.Document, error) {
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(d.Body); err != nil {
		return document.Document{}, err
	}

	byts := buf.Bytes()
	for old, new := range smartquotes {
		byts = bytes.ReplaceAll(byts, []byte(old), []byte(new))
	}
	d.Body = bytes.NewBuffer(byts)

	return d, nil
}
