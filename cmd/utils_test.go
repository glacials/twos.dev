package cmd

import "testing"

func TestCommonPathAncestor(t *testing.T) {
	for _, test := range []struct {
		a        string
		b        string
		expected string
	}{
		{
			a:        "a",
			b:        "a",
			expected: "a",
		},
		{
			a:        "a/b/c",
			b:        "a/b/c",
			expected: "a/b/c",
		},
		{
			a:        "a/b/c",
			b:        "a/b/c/d",
			expected: "a/b/c",
		},
		{
			a:        "a/b",
			b:        "a/c",
			expected: "a",
		},
		{
			a:        "a/b/c",
			b:        "a/c",
			expected: "a",
		},
		{
			a:        "a/b/c",
			b:        "a/c/d",
			expected: "a",
		},
		{
			a:        "a",
			b:        "b",
			expected: "",
		},
	} {
		actual := commonPathAncestor(test.a, test.b)
		if actual != test.expected {
			t.Errorf("expected common path of %s for %s and %s, but got %s", test.expected, test.a, test.b, actual)
		}
	}
}
