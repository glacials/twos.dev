package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ignoreSrcDirs = map[string]struct{}{
		"src/asciiart":  {},
		"src/templates": {},
	}
)

func buildTheWorld() error {
	templater, err := NewTemplateBuilder()
	if err != nil {
		return fmt.Errorf("can't make template builder: %w", err)
	}

	if err := filepath.WalkDir("src", func(src string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk src for build `%s`: %w", src, err)
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() == ".DS_Store" {
			return nil
		}
		for path := range ignoreSrcDirs {
			if strings.HasPrefix(src, path) {
				return filepath.SkipDir
			}
		}

		if ok, err := filepath.Match("*.md", d.Name()); err != nil {
			return fmt.Errorf("can't match `%s` to *.md: %w", src, err)
		} else if ok {
			if err := templater.markdownBuilder(src, dst); err != nil {
				return fmt.Errorf("can't build Markdown in `%s`: %w", src, err)
			}

			return nil
		}

		if ok, err := filepath.Match("*.html", d.Name()); err != nil {
			return fmt.Errorf("can't match `%s` to *.html: %w", src, err)
		} else if ok {
			if err := templater.htmlBuilder(src, dst); err != nil {
				return fmt.Errorf("can't build HTML in `%s`: %w", src, err)
			}

			return nil
		}

		if ok, err := filepath.Match("*.[jJ][pP][gG]", d.Name()); err != nil {
			return fmt.Errorf("can't match `%s` to *.jpg: %w", src, err)
		} else if ok {
			if err := imageBuilder(src, dst); err != nil {
				return fmt.Errorf("can't build JPEG `%s`: %w", src, err)
			}

			return nil
		}

		return fmt.Errorf("don't know what to do with file `%s`", src)
	}); err != nil {
		return fmt.Errorf("can't watch file: %w", err)
	}

	if err := buildStaticDir(staticAssetsDir, dst); err != nil {
		return fmt.Errorf("can't build static assets: %w", err)
	}

	if err := buildFormatting(dst); err != nil {
		return fmt.Errorf("can't post-process build dir: %w", err)
	}

	return nil
}
