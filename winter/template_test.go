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
			expected: "abc123",
		},
		{
			name:     "SimpleTemplate",
			input:    "{{ add 1 2 }}",
			expected: "3",
		},
		{
			name:     "Template",
			input:    `{{ template "hello_world.tmpl" }}`,
			expected: "Hello, world!\n",
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			src := fmt.Sprintf("src/test/%s", test.name)
			doc := NewTemplateDocument(
				src,
				NewMetadata(src),
				nil,
				nil,
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
