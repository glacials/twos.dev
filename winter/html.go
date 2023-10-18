package winter

import (
	"bytes"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// firstTag returns the first element with the given tag which is a descendant
// of n.
func firstTag(n *html.Node, t atom.Atom) *html.Node {
	if n.Type == html.ElementNode && n.DataAtom == t {
		return n
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if el := firstTag(child, t); el != nil {
			return el
		}
	}

	return nil
}

// clone returns a deep copy of n.
func clone(n *html.Node) (*html.Node, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, n); err != nil {
		return nil, err
	}
	els, err := html.ParseFragment(&buf, n.Parent)
	if err != nil {
		return nil, err
	}

	return els[0], nil
}

// allOfTypes returns all descendant nodes of n with any of the given types. The
// returned slice is sorted in the same way the document was, with parent nodes
// coming before their children.
func allOfTypes(n *html.Node, t map[atom.Atom]struct{}) (m []*html.Node) {
	if _, ok := t[n.DataAtom]; ok {
		m = append(m, n)
	}
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		m = append(m, allOfTypes(child, t)...)
	}
	return
}

// attr returns the attr attribute of the given element node.
func attr(n *html.Node, attr atom.Atom) string {
	for _, a := range n.Attr {
		if a.Key == attr.String() {
			return a.Val
		}
	}
	return ""
}
