package transform

import (
	"bytes"

	"github.com/glacials/twos.dev/winter/document"
)

var emdashes = map[string]string{
	" -- ":      "---",
	"–":         "---",
	" – ":       "---",
	"&mdash;":   "---",
	" &mdash; ": "---",
}

// ReplaceEmDashes returns a reader identical to the given reader but with em
// dashes and their common substitutions replaced with monospace-friendly em
// dashes. This must happen before Markdown parsing because go-markdown uses
// LaTeX-style en and em dash syntax: "--" is an en dash while "---" is an em
// dash.
//
// ReplaceEmDashes implements document.Transformation.
func ReplaceEmDashes(d document.Document) (document.Document, error) {
	if d.Shortname == "dashes" {
		// dashes.html is a post specifically about dashes; let it do its thing
		return d, nil
	}

	for old, nw := range emdashes {
		d.Body = bytes.ReplaceAll(d.Body, []byte(old), []byte(nw))
	}

	return d, nil
}
