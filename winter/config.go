package winter

import (
	"net/url"
)

type Config struct {
	AuthorEmail string
	AuthorName  string
	Debug       bool
	Desc        string
	Dist        string
	Domain      url.URL
	Name        string
	// TODO: Use this for copyright in page footer
	Since int
}
