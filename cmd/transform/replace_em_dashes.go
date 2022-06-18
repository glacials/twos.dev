package transform

import (
	"bytes"

	"github.com/glacials/twos.dev/cmd/document"
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
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(d.Body); err != nil {
		return document.Document{}, err
	}

	byts := buf.Bytes()
	for old, nw := range emdashes {
		byts = bytes.ReplaceAll(byts, []byte(old), []byte(nw))
	}
	d.Body = bytes.NewBuffer(byts)

	return d, nil
}
