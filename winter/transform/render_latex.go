package transform

import (
	"regexp"

	"github.com/glacials/twos.dev/winter/document"
)

var inlineLaTeX = regexp.MustCompile("\\\\\\(.*?\\\\\\)")
var blockLaTeX = regexp.MustCompile("\\\\\\[.*?\\\\\\]")

// RenderLaTeX converts any LaTeX in the document to HTML.
//
// RenderLaTeX imlpements document.Transformation.
func RenderLaTeX(d document.Document) (document.Document, error) {
	// TODO: Implement. Currently LaTeX is handled by a frontend JS library.

	return d, nil
}
