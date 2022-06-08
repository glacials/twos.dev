package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	ignoreFilenames = map[string]struct{}{"README.md": {}, ".DS_Store": {}}
)

func buildTheWorld() error {
	if err := buildStaticDir(staticAssetsDir, dst); err != nil {
		return fmt.Errorf("can't build static assets: %w", err)
	}

	templater, err := NewTemplateBuilder()
	if err != nil {
		return fmt.Errorf("can't make template builder: %w", err)
	}

	if err := filepath.WalkDir("src/warm", func(src string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk src/warm to build `%s`: %w", src, err)
		}
		if d.IsDir() {
			return nil
		}
		if _, ok := ignoreFilenames[d.Name()]; ok {
			return nil
		}

		if ok, err := filepath.Match("*.md", d.Name()); err != nil {
			return fmt.Errorf("can't match `%s` to *.md: %w", src, err)
		} else if ok {
			if err := templater.markdownBuilder(src, dst); err != nil {
				return fmt.Errorf("can't build Markdown in `%s`: %w", src, err)
			}

			return nil
		}

		return fmt.Errorf("don't know what to do with `%s`", src)
	}); err != nil {
		return fmt.Errorf("can't walk src/warm: %w", err)
	}

	if err := filepath.WalkDir("src/cold", func(src string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk src/cold to build `%s`: %w", src, err)
		}
		if d.IsDir() {
			return nil
		}
		if _, ok := ignoreFilenames[d.Name()]; ok {
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

		return fmt.Errorf("don't know what to do with `%s`", src)
	}); err != nil {
		return fmt.Errorf("can't walk src/cold: %w", err)
	}

	if err := filepath.WalkDir("src/img", func(src string, d os.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("can't walk src/img to build `%s`: %w", src, err)
		}
		if d.IsDir() {
			return nil
		}
		if _, ok := ignoreFilenames[d.Name()]; ok {
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

		return fmt.Errorf("don't know what to do with `%s`", src)
	}); err != nil {
		return fmt.Errorf("can't walk src/img: %w", err)
	}

	if err := buildFormatting(dst); err != nil {
		return fmt.Errorf("can't post-process build dir: %w", err)
	}

	return nil
}
