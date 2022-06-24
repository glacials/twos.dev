package transform

import (
	"bytes"
	"fmt"

	"github.com/glacials/twos.dev/cmd/document"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type heading struct {
	level       atom.Atom
	id          string
	title       string
	subheadings []heading
}

// BuildTOC checks if the document requested a table of contents (via
// frontmatter) and builds one into it if so.
func BuildTOC(d document.Document) (document.Document, error) {
	if !d.TOC {
		return d, nil
	}

	root, err := html.Parse(bytes.NewBuffer(d.Body))
	if err != nil {
		return document.Document{}, err
	}

	headings := findHeadings(root, 2, 4)
	if err != nil {
		return document.Document{}, fmt.Errorf("can't find headings: %w", err)
	}

	ul := buildTOC(headings)

	// Insert TOC at bottom of intro, aka above first H2
	h2, ok := firstOfType(root, atom.H2)
	if !ok {
		return document.Document{}, fmt.Errorf(
			"can't insert TOC into a document with no H2",
		)
	}

	h2.Parent.InsertBefore(ul, h2)

	var buf bytes.Buffer
	if err := html.Render(&buf, root); err != nil {
		return document.Document{}, err
	}

	d.Body = buf.Bytes()

	return d, nil
}

var levels = map[int]atom.Atom{
	1: atom.H1,
	2: atom.H2,
	3: atom.H3,
	4: atom.H4,
	5: atom.H5,
	6: atom.H6,
}

// findHeadings returns a slice of headings from size biggest (e.g. h1) to size
// smallest (e.g. h6) inclusive. Each heading recursively contains its
// subheadings; for example a document that flows like h1 h2 h3 h2 will return
// one h1 heading which has two h2 subheadings, the first of which has an h3
// subheading.
func findHeadings(start *html.Node, biggest, smallest int) []heading {
	if start == nil {
		return []heading{}
	}

	var headings []heading

	for el := start; el != nil; el = el.NextSibling {
		for i, level := range levels {
			if i != 1 && i < biggest {
				if el.DataAtom == level {
					return headings
				}
			}
		}
		if el.DataAtom == levels[biggest] {
			headings = append(
				headings,
				heading{
					id:          id(el),
					level:       levels[biggest],
					title:       el.FirstChild.Data,
					subheadings: findHeadings(nextEl(el), biggest+1, smallest),
				},
			)
		}
		headings = append(
			headings,
			findHeadings(el.FirstChild, biggest, smallest)...)
	}

	return headings
}

func headingToText(heading *html.Node) (string, error) {
	var buf bytes.Buffer
	if err := html.Render(&buf, heading); err != nil {
		return "", fmt.Errorf("can't render heading to text: %w", err)
	}
	return buf.String(), nil
}

func nextEl(el *html.Node) *html.Node {
	if el == nil {
		return nil
	}

	if el.NextSibling != nil {
		return el.NextSibling
	}

	return nextEl(el.Parent)
}

// buildTOC returns a list node containing the headings as <li> item links.
func buildTOC(headings []heading) *html.Node {
	list := html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Ul,
		Data:     "ol",
	}

	for _, heading := range headings {
		li := html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Li,
			Data:     "li",
		}
		a := html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.A,
			Data:     "a",
			Attr: []html.Attribute{
				{Key: "href", Val: fmt.Sprintf("#%s", heading.id)},
			},
		}
		text := html.Node{
			Type: html.TextNode,
			Data: heading.title,
		}

		li.AppendChild(&a)
		a.AppendChild(&text)

		if len(heading.subheadings) > 0 {
			li.AppendChild(buildTOC(heading.subheadings))
		}
		list.AppendChild(&li)
	}

	return &list
}

// id returns the id attribute of the given node, or the empty string if one is
// not present.
func id(node *html.Node) string {
	if node == nil {
		return ""
	}

	for _, attr := range node.Attr {
		if attr.Key == "id" {
			return attr.Val
		}
	}

	return ""
}
