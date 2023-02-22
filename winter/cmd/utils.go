package cmd

import (
	"strings"

	"github.com/lithammer/dedent"
	"github.com/mitchellh/go-wordwrap"
)

func wrap(text string) string {
	text = dedent.Dedent(text)
	text = strings.ReplaceAll(text, "\n", " ")
	text = wordwrap.WrapString(text, 80)
	text = strings.TrimSpace(text)
	return text
}
