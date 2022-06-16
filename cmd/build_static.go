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
