package cmd

// Builder is a function that builds the src file into the dst directory.
type Builder func(src, dst string, cfg Config) error
