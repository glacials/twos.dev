// Winter is a static website generator that allows for easy extensibility. The
// generator itself is a lightweight framework that only knows how to process
// transformations.
//
// A transformation is a Go function that receives a document as input, applies
// some modification to it, and returns the result. A Winter configuration is
// defined as a list of transformations.
//
// Some examples of transformations are converting Markdown to HTML, scraping a
// piece of frontmatter from the document, or executing the document as a
// template.
package main

import "twos.dev/winter/cmd"

func main() {
	cmd.Execute()
}
