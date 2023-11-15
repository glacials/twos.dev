package winter

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

type testCase struct {
	name     string
	input    string
	expected string
}

func TestMarkdown(t *testing.T) {
	for _, test := range []testCase{
		{
			name:  "Heading",
			input: `# Heading 1`,
			expected: `<h1>Heading 1</h1>
`,
		},
		{
			name:     "SimpleTemplate",
			input:    `{{ add 1 2 }}`,
			expected: "<p>{{ add 1 2 }}</p>\n",
		},
		{
			name:     "Template",
			input:    `{{ template "\_writing.html.tmpl" }}`,
			expected: "<p>{{ template \"_writing.html.tmpl\" }}</p>\n",
		},
		{
			name:     "Bold",
			input:    `**foo**`,
			expected: "<p><strong>foo</strong></p>\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			src := fmt.Sprintf("src/test/%s", test.name)
			doc := NewMarkdownDocument(src, NewMetadata(src), nil)
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
