package main

import "twos.dev/winter"

// Builder represents a function that builds a source file src into a
// destination directory dst.
type Builder func(src, dst string, cfg winter.Config) error
