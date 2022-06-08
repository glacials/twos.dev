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
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

// staticFileBuilder returns a builder (it is not one itself) that will build
// src into a rel-relative mirrored directory in dst. For example,
// staticFileBuilder("public") returns a builder func(src, dst string) which
// when called with f("public/a/b/c", "dist") puts the file in dist/a/b/c.
func staticFileBuilder(rel string) func(src, dst string) error {
	return func(src, dst string) error {
		rel, err := filepath.Rel(rel, src)
		if err != nil {
			return fmt.Errorf("can't get `%s` relative to `%s`: %w", src, rel, err)
		}

		if err := os.MkdirAll(filepath.Dir(filepath.Join(dst, rel)), 0755); err != nil {
			return fmt.Errorf("can't make static asset dir `%s`: %w", dst, err)
		}

		s, err := os.Open(src)
		if err != nil {
			log.Printf("can't read static file `%s`: %s", src, err)
			// TODO: This happens sometimes
			return nil
		}
		defer s.Close()

		d, err := os.Create(filepath.Join(dst, rel))
		if err != nil {
			return fmt.Errorf(
				"can't write static file `%s` to `%s`: %w",
				src,
				filepath.Join(dst, rel),
				err,
			)
		}
		defer d.Close()

		if _, err := io.Copy(d, s); err != nil {
			return fmt.Errorf("can't build static asset %s: %w", src, err)
		}

		return nil
	}
}

func buildStaticDir(staticDir, destinationDir string) error {
	if err := copy.Copy(staticDir, destinationDir); err != nil {
		return fmt.Errorf("can't build static assets: %w", err)
	}

	return nil
}
