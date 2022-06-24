package transform

import (
	"bytes"
	"fmt"

	"github.com/glacials/twos.dev/cmd/document"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// DiscoverTitle finds the title inside the document by looking for the first
// <h1> heading, if any, and stores it in the document.
//
// DiscoverTitle implements document.Transformation.
func DiscoverTitle(d document.Document) (document.Document, error) {
	h, err := html.Parse(bytes.NewBuffer(d.Body))
	if err != nil {
		return document.Document{}, fmt.Errorf("can't parse HTML: %w", err)
	}

	if h1, ok := firstOfType(h, atom.H1); ok {
		for child := h1.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				d.Title = child.Data
				return d, nil
			}
		}
	}

	return d, nil
}

// firstOfType returns the first node inside n with the given type.
func firstOfType(n *html.Node, t atom.Atom) (*html.Node, bool) {
	if n.DataAtom == t {
		return n, true
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if el, ok := firstOfType(child, t); ok {
			return el, true
		}
	}

	return nil, false
}
