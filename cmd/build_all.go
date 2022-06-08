/*
Copyright © 2022 Benjamin Carlsson

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ignoreFiles = map[string]struct{}{
		"README.md": {},
	}
	ignoreDirectories = map[string]struct{}{
		".git":         {},
		".github":      {},
		"node_modules": {},
	}
)

func buildTheWorld() error {
	if err := os.MkdirAll(dst, 0755); err != nil {
		return fmt.Errorf("can't mkdir `%s`: %w", dst, err)
	}
	if err := filepath.WalkDir(".", func(src string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk towards `%s`: %w", src, err)
		}
		if _, ok := ignoreFiles[d.Name()]; ok {
			return nil
		}
		if _, ok := ignoreDirectories[d.Name()]; ok {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		for pattern, builder := range builders {
			if ok, err := filepath.Match(pattern, src); err != nil {
				return fmt.Errorf("can't match `%s` to %s: %w", src, pattern, err)
			} else if ok {
				if err := builder(src, dst); err != nil {
					return fmt.Errorf("can't build `%s`: %w", src, err)
				}
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("can't walk to build the world: %w", err)
	}

	if err := buildFormatting(dst); err != nil {
		return fmt.Errorf("can't post-process build dir: %w", err)
	}

	return nil
}
