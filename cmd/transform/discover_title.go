package transform

import (
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
	h, err := html.Parse(d.Body)
	if err != nil {
		return document.Document{}, fmt.Errorf("can't parse HTML: %w", err)
	}

	if h1, ok := firstH1(h); ok {
		for child := h1.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.TextNode {
				d.Title = child.Data
				return d, nil
			}
		}
	}

	return d, nil
}

func firstH1(n *html.Node) (*html.Node, bool) {
	if n.DataAtom == atom.H1 {
		return n, true
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if h1, ok := firstH1(child); ok {
			return h1, true
		}
	}

	return nil, false
}
