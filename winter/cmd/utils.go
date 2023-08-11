package cmd

import (
	"regexp"
	"strings"

	"github.com/lithammer/dedent"
	"github.com/mitchellh/go-wordwrap"
)

var (
	linebreak      = regexp.MustCompile(`([^\s])[^\S\r\n]*\n[^\S\r\n]*([^\s])`)
	multilinebreak = regexp.MustCompile(`\n\s*\n+`)
)

// wrap formats a multi-line string for use as a command description. Single newlines
// are removed, common indentation is removed, and the text is wrapped to 80 characters.
func wrap(text string) string {
	text = dedent.Dedent(text)
	text = linebreak.ReplaceAllString(text, "$1 $2")
	text = multilinebreak.ReplaceAllString(text, "\n\n")
	text = wordwrap.WrapString(text, 80)
	text = strings.TrimSpace(text)
	return text
}
