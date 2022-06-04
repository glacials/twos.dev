package cmd

import (
	"path/filepath"
	"strings"
)

// commonPathAncestor returns the path to the "lowest" directory that both a and
// b are contained within.
func commonPathAncestor(a, b string) string {
	aparts := strings.Split(a, string(filepath.Separator))
	bparts := strings.Split(b, string(filepath.Separator))

	if len(aparts) > len(bparts) {
		return commonPathAncestor(b, a)
	}

	for i := range aparts {
		if aparts[i] != bparts[i] {
			return strings.Join(aparts[0:i], string(filepath.Separator))
		}
	}

	return a
}
