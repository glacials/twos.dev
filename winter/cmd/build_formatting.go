package main

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
