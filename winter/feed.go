package winter

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

const favicon = "/favicon.ico"

func (s *Substructure) writefeed() error {
	now := time.Now()
	siteURL := url.URL{Scheme: "https", Host: s.cfg.Production.URL}

	feed := feeds.Feed{
		Author: &feeds.Author{
			Name:  s.cfg.Author.Name,
			Email: s.cfg.Author.Email,
		},
		Copyright: fmt.Sprintf(
			"Copyright %d–%d %s",
			s.cfg.Since,
			now.Year(),
			s.cfg.Author.Name,
		),
		Created:     time.Date(2022, 4, 7, 0, 0, 0, 0, time.UTC),
		Description: s.cfg.Description,
		Image: &feeds.Image{
			Url: (&url.URL{Scheme: "https", Host: s.cfg.Production.URL, Path: favicon}).String(),
		},
		Items:   []*feeds.Item{},
		Link:    &feeds.Link{Href: siteURL.String()},
		Title:   s.cfg.Name,
		Updated: now,
	}
	atomFeed := (&feeds.Atom{Feed: &feed}).AtomFeed()
	atomFeed.Icon = (&url.URL{Scheme: "https", Host: s.cfg.Production.URL, Path: favicon}).String()

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

		bodyStr = strings.ReplaceAll(bodyStr, "&#34;", "\"")

		itemURL := siteURL
		itemURL.Path = doc.Metadata().WebPath
		feed.Items = append(feed.Items, &feeds.Item{
			Author:      feed.Author,
			Id:          doc.Metadata().WebPath,
			Title:       doc.Metadata().Title,
			Content:     bodyStr,
			Description: bodyStr,
			Link:        &feeds.Link{Href: itemURL.String()},
			Created:     doc.Metadata().CreatedAt,
			Updated:     doc.Metadata().UpdatedAt,
		})
	}

	atom, err := os.Create("dist/feed.atom")
	if err != nil {
		return fmt.Errorf("cannot create Atom feed: %w", err)
	}
	if err := feed.WriteAtom(atom); err != nil {
		return fmt.Errorf("cannot write Atom feed: %w", err)
	}

	rss, err := os.Create("dist/feed.rss")
	if err != nil {
		return fmt.Errorf("cannot create RSS feed: %w", err)
	}
	if err := feed.WriteRss(rss); err != nil {
		return fmt.Errorf("cannot write RSS feed: %w", err)
	}

	return nil
}
