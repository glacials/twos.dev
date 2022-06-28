package transform

import (
	"bytes"
	"fmt"

	"twos.dev/winter/document"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type heading struct {
	node *html.Node

	level       atom.Atom
	id          string
	title       string
	subheadings []heading
}

// RenderTOC checks if the document requested a table of contents (via
// frontmatter) and builds one into it if so.
func RenderTOC(d document.Document) (document.Document, error) {
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

	for _, heading := range headings {
		if err := addReturnToTopLinks(heading); err != nil {
			return document.Document{}, fmt.Errorf(
				"can't add return to top links: %w",
				err,
			)
		}
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
					node:        el,
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

func addReturnToTopLinks(heading heading) error {
	totop := html.Node{
		Type:     html.ElementNode,
		Data:     atom.A.String(),
		DataAtom: atom.A,
		Attr: []html.Attribute{
			{Key: atom.Href.String(), Val: "#toc"},
			{
				Key: atom.Style.String(),
				Val: "text-decoration: none;",
			},
		},
	}
	totop.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: "↑",
	})

	tohere := html.Node{
		Type:     html.ElementNode,
		Data:     atom.A.String(),
		DataAtom: atom.A,
		Attr: []html.Attribute{
			{Key: atom.Href.String(), Val: fmt.Sprintf("#%s", heading.id)},
			{
				Key: atom.Style.String(),
				Val: "text-decoration: none;",
			},
		},
	}
	tohere.AppendChild(&html.Node{
		Type: html.TextNode,
		Data: "#",
	})

	span := html.Node{
		Type:     html.ElementNode,
		Data:     atom.Span.String(),
		DataAtom: atom.Span,
		Attr: []html.Attribute{
			{Key: atom.Style.String(), Val: "margin-left: 0.5em;"},
		},
	}

	span.AppendChild(&tohere)
	span.AppendChild(&totop)

	heading.node.Parent.InsertBefore(&span, heading.node.NextSibling)

	// Recurse through subheadings
	for _, heading := range heading.subheadings {
		if err := addReturnToTopLinks(heading); err != nil {
			return fmt.Errorf("can't add return to top link: %w", err)
		}
	}

	return nil
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
		Data:     atom.Ol.String(),
		Attr: []html.Attribute{
			{Key: atom.Id.String(), Val: "toc"},
		},
	}

	for _, heading := range headings {
		li := html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.Li,
			Data:     atom.Li.String(),
		}
		a := html.Node{
			Type:     html.ElementNode,
			DataAtom: atom.A,
			Data:     atom.A.String(),
			Attr: []html.Attribute{
				{Key: atom.Href.String(), Val: fmt.Sprintf("#%s", heading.id)},
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
		if attr.Key == atom.Id.String() {
			return attr.Val
		}
	}

	return ""
}
