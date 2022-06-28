package transform

import (
	"bytes"

	"twos.dev/winter/document"
)

var smartquotes = map[string]string{
	"“":       "\"",
	"”":       "\"",
	"&ldquo;": "\"",
	"&rdquo;": "\"",
	"&quot;":  "\"",
	"&#34;":   "\"",
	"&#39;":   "'",
}

// ReplaceSmartQuotes returns a reader identical to the given reader but with
// smart quotes et al. replaced with dumb quotes. This must happen after
// Markdown parsing, as go-markdown replaces dumb quotes with smart quotes even
// when used to indicate template string literals.
//
// ReplaceSmartQuotes implements document.Transformation.
func ReplaceSmartQuotes(d document.Document) (document.Document, error) {
	for old, new := range smartquotes {
		d.Body = bytes.ReplaceAll(d.Body, []byte(old), []byte(new))
	}

	return d, nil
}
