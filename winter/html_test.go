package winter

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func surround(html string) string {
	return fmt.Sprintf("<html><head></head><body>%s</body></html>", html)
}

func TestHTML(t *testing.T) {
	for _, test := range []testCase{
		{
			name:     "Heading",
			input:    "<h1>Heading 1</h1>",
			expected: surround("<a href=\"Heading\" class=\"post-title\"><h1>Heading 1</h1></a>"),
		},
		{
			name:     "SimpleTemplate",
			input:    `{{ add 1 2 }}`,
			expected: surround("{{ add 1 2 }}"),
		},
		{
			name:     "Template",
			input:    `{{ template "_writing.html.tmpl" }}`,
			expected: surround("{{ template \"_writing.html.tmpl\" }}"),
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			src := fmt.Sprintf("src/test/%s", test.name)
			doc := NewHTMLDocument(
				src,
				NewMetadata(src),
				filepath.Join("testdata", "templates"),
				nil,
			)
			if err := doc.Load(strings.NewReader(test.input)); err != nil {
				t.Errorf("load failed: %s", err)
			}
			var actual bytes.Buffer
			if err := doc.Render(&actual); err != nil {
				t.Errorf("render failed: %s", err)
			}
			assert.Equal(t, test.expected, actual.String())
		})
	}
}
