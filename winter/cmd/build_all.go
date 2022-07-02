package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"twos.dev/winter"
)

var (
	ignoreFiles = map[string]struct{}{
		"README.md": {},
		".DS_Store": {},
	}
	ignoreDirectories = map[string]struct{}{
		".git":         {},
		".github":      {},
		dist:           {},
		"node_modules": {},
	}
)

func buildAll(
	dist string,
	builders map[string]Builder,
	cfg winter.Config,
) error {
	if err := os.RemoveAll(dist); err != nil {
		return err
	}
	if err := os.MkdirAll(dist, 0755); err != nil {
		return fmt.Errorf("can't mkdir `%s`: %w", dist, err)
	}

	seen := map[string]struct{}{}

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

				if err := builder(src, dist, cfg); err != nil {
					return fmt.Errorf("can't build `%s`: %w", src, err)
				}
			}
		}

		return nil
	}); err != nil {
		return fmt.Errorf("can't walk to build the world: %w", err)
	}

	s, err := winter.Discover(cfg)
	if err != nil {
		return err
	}

	if err := s.Execute(dist); err != nil {
		return err
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
