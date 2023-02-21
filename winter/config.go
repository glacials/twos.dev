package winter

import (
	"net/url"
)

// Config is a configuration for the Winter build.
type Config struct {
	// AuthorEmail is the email address of the website author. This is
	// used as metadata for the RSS feed.
	AuthorEmail string
	// AuthorName is the name of the website author. This is used as
	// metadata for the RSS feed.
	AuthorName string
	// Debug is a flag that enables debug mode.
	Debug bool
	// Desc is the description of the website. This is used as metadata
	// for the RSS feed.
	Desc string
	// Dist is the path to the distribution directory. If blank,
	// defaults to "./dist" relative to the working directory.
	Dist string
	// Domain is the domain of the website. This is used as metadata for
	// the RSS feed.
	Domain url.URL
	// Name is the name of the website. This is used in various places
	// in and out of templates.
	Name string
	// Since is the year the website was established, whether through
	// Winter or otherwise. This is used as metadata for the RSS feed.
	//
	// TODO: Use this for copyright in page footer
	Since int
	// SourceDirectories is an additional list of directories to search
	// for source files beyond ./src/cold and ./src/warm.
	SourceDirectories []string
}
