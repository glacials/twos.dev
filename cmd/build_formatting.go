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
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func buildFormatting(sourceDir string) error {
	if err := filepath.WalkDir(sourceDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk path `%s` for HTML templates: %w", path, err)
		}

		if d.IsDir() {
			return nil
		}

		matched, err := filepath.Match("*.html", strings.ToLower(d.Name()))
		if err != nil {
			return fmt.Errorf("can't match against filename `%s` to format HTML: %w", d.Name(), err)
		}
		if !matched {
			return nil
		}

		if strings.Contains(path, "node_modules") {
			return nil
		}

		body, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("can't read HTML file `%s`: %w", path, err)
		}

		re, err := regexp.Compile("(\\s+--\\s+)|(\\s*&mdash;\\s*)|(\\s*—\\s*)")
		if err != nil {
			return fmt.Errorf("can't compile regexp: %w", err)
		}

		body = re.ReplaceAll(body, []byte(" -- "))

		if err := os.WriteFile(path, body, 0644); err != nil {
			return fmt.Errorf("can't write formatting-changed file back to disk: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("can't perform formatting loop: %w", err)
	}
	return nil
}
