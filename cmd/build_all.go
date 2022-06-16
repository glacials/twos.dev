package cmd

import (
	"fmt"
	"log"
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
				log.Println("building", src)
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
