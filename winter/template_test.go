package winter

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestTemplate(t *testing.T) {
	for _, test := range []testCase{
		{
			name:     "NoOp",
			input:    "abc123",
			expected: "abc123\n",
		},
		{
			name:     "SimpleTemplate",
			input:    "{{ add 1 2 }}",
			expected: "3\n",
		},
		{
			name:     "Template",
			input:    `{{ template "hello_world.tmpl" }}`,
			expected: "Hello, world!\n\n",
		},
		{
			name:     "Image",
			input:    `<img src="/path/to/image.jpg" alt="Alt text" />`,
			expected: "<img src=\"/path/to/image.jpg\" alt=\"Alt text\" />\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			src := fmt.Sprintf("src/test/%s", test.name)
			doc := NewTemplateDocument(
				src,
				NewMetadata(src, filepath.Join("testdata", "templates")),
				nil,
				nil,
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
