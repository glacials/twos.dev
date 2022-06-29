package cmd

// Builder represents a function that builds a source file src into a
// destination directory dst.
type Builder func(src, dst string, cfg Config) error
