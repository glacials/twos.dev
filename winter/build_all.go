package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	ignoreFiles = map[string]struct{}{
		"README.md": {},
		".DS_Store": {},
	}
	ignoreDirectories = map[string]struct{}{
		".git":         {},
		".github":      {},
		"dist":         {},
		"node_modules": {},
	}
)

func buildTheWorld() error {
	seen := map[string]struct{}{}

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
				log.Println("building", src)

				if built(src, seen) {
					return fmt.Errorf(
						"a file like %s was already built from another directory",
						src,
					)
				}

				if err := builder(src, dst); err != nil {
					return fmt.Errorf("can't build `%s`: %w", src, err)
				}
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("can't walk to build the world: %w", err)
	}

	return nil
}

func built(src string, seen map[string]struct{}) bool {
	ext := strings.ToLower(filepath.Ext(src))

	if ext == ".mov" || ext == ".mp4" {
		// Free pass; we have multiple formats of the same video
		return false
	}

	filename := filepath.Base(src)

	if filename == "tattoo.webp" {
		// Historical exception
		return false
	}

	filename, _, _ = strings.Cut(filename, ".")
	if _, ok := seen[filename]; ok {
		return true
	}

	seen[filename] = struct{}{}
	return false
}
