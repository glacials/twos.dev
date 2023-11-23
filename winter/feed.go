package winter

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/feeds"
)

func (s *Substructure) writefeed() error {
	now := time.Now()
	feed := feeds.Feed{
		Title:       s.cfg.Name,
		Description: s.cfg.Description,
		Author: &feeds.Author{
			Name:  s.cfg.Author.Name,
			Email: s.cfg.Author.Email,
		},
		Link: &feeds.Link{Href: (&url.URL{Scheme: "https", Host: s.cfg.Hostname}).String()},
		Copyright: fmt.Sprintf(
			"Copyright %d–%d %s",
			s.cfg.Since,
			now.Year(),
			s.cfg.Author.Name,
		),
		Items: []*feeds.Item{},

		Created: now,
		Updated: now,
	}

	for _, doc := range s.docs {
		if doc.Metadata().Kind != post {
			continue
		}
		var buf bytes.Buffer
		layout := doc.Metadata().Layout
		doc.Metadata().Layout = ""
		if err := doc.Render(&buf); err != nil {
			return err
		}
		doc.Metadata().Layout = layout

		bodyStr := buf.String()
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          doc.Metadata().WebPath,
			Title:       doc.Metadata().Title,
			Author:      feed.Author,
			Content:     bodyStr,
			Description: bodyStr,
			Link: &feeds.Link{
				Href: fmt.Sprintf("%s/%s.html", feed.Link.Href, doc.Metadata().WebPath),
			},
			Created: doc.Metadata().CreatedAt,
			Updated: doc.Metadata().UpdatedAt,
		})
	}

	atom, err := feed.ToAtom()
	if err != nil {
		return err
	}
	if err := os.WriteFile("dist/feed.atom", []byte(atom), 0o644); err != nil {
		return err
	}

	rss, err := feed.ToRss()
	if err != nil {
		return err
	}
	if err := os.WriteFile("dist/feed.rss", []byte(rss), 0o644); err != nil {
		return err
	}

	return nil
}
