package winter

import (
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// firstOfType returns the first descendant node of n with the given type.
func firstOfType(n *html.Node, t atom.Atom) *html.Node {
	if n.DataAtom == t {
		return n
	}

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if el := firstOfType(child, t); el != nil {
			return el
		}
	}

	return nil
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
