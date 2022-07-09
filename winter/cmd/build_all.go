package main

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
