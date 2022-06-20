package transform

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/glacials/twos.dev/cmd/document"
	"github.com/gorilla/feeds"
)

var now = time.Now()
var (
	feed = feeds.Feed{
		Title:       "twos.dev",
		Description: "misc thoughts from Benjamin Carlsson",
		Author: &feeds.Author{
			Name:  "Benjamin Carlsson",
			Email: "ben@twos.dev",
		},
		Link: &feeds.Link{
			Href: "https://twos.dev",
		},
		Copyright: fmt.Sprintf("Copyright 2021–%d Benjamin Carlsson", now.Year()),
		Items:     []*feeds.Item{},

		Created: now,
		Updated: now,
	}
)

// UpdateFeeds re-builds RSS and ATOM feeds with the new information from the
// document.
//
// UpdateFeeds implements document.Transformation.
func UpdateFeeds(d document.Document) (document.Document, error) {
	if d.Type != document.PostType {
		return d, nil
	}
	if d.Title == "" {
		return document.Document{}, fmt.Errorf(
			"can't generate RSS for document without title (%s)",
			d.Shortname,
		)
	}

	var item *feeds.Item
	for _, itm := range feed.Items {
		if itm.Id == d.Shortname {
			item = itm
		}
	}
	if item == nil {
		item = &feeds.Item{}
		feed.Items = append([]*feeds.Item{item}, feed.Items...)
	}

	item.Id = d.Shortname
	item.Title = d.Title
	item.Author = feed.Author
	item.Content = string(d.Body)
	item.Description = string(d.Body)
	item.Link = &feeds.Link{
		Href: fmt.Sprintf("%s/%s.html", feed.Link.Href, d.Shortname),
	}
	item.Created = d.CreatedAt
	item.Updated = d.UpdatedAt

	atom, err := feed.ToAtom()
	if err != nil {
		return document.Document{}, err
	}
	if err := ioutil.WriteFile("dist/feed.atom", []byte(atom), 0644); err != nil {
		return document.Document{}, err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return document.Document{}, err
	}
	if err := ioutil.WriteFile("dist/feed.rss", []byte(rss), 0644); err != nil {
		return document.Document{}, err
	}

	return d, nil
}
