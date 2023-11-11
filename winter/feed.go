package winter

import (
	"bytes"
	"fmt"
	"html/template"
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

	for _, post := range s.posts() {
		body, err := post.Build()
		if err != nil {
			return err
		}
		dest, err := post.Dest()
		if err != nil {
			return err
		}
		t := template.New(post.Source)
		funcs, err := s.funcmap(post)
		if err != nil {
			return fmt.Errorf("can't generate funcmap: %w", err)
		}
		_ = t.Funcs(funcs)
		_, err = t.Parse(string(body))
		if err != nil {
			return fmt.Errorf("cannot parse page for feed: %w", err)
		}
		if err := loadDeps(t); err != nil {
			return fmt.Errorf("can't load dependency templates for %q: %w", t.Name(), err)
		}
		var buf bytes.Buffer
		if err = t.Execute(&buf, post); err != nil {
			return fmt.Errorf("cannot execute document %q: %w", post.Source, err)
		}

		feed.Items = append(feed.Items, &feeds.Item{
			Id:          dest,
			Title:       post.Title(),
			Author:      feed.Author,
			Content:     string(buf.String()),
			Description: string(body),
			Link: &feeds.Link{
				Href: fmt.Sprintf("%s/%s.html", feed.Link.Href, dest),
			},
			Created: post.CreatedAt(),
			Updated: post.UpdatedAt(),
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
