package main

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"golang.org/x/net/websocket"
	"twos.dev/winter"
)

// Reloader watches the filesystem for changes to relevant files so it can
// reload the browser using WebSockets.
type Reloader struct {
	// Builders is a mapping of filepath patterns to functions which can build
	// paths that match those patterns. Deprecated: Use Substructure instead.
	Builders     map[string]Builder
	Ignore       map[string]struct{}
	Substructure *winter.Substructure

	closeSockets chan struct{}
	listeners    []*websocket.Conn
	stop         chan struct{}
	watcher      *fsnotify.Watcher
}

// Handler returns a function that handles incoming WebSocket connections which
// implements http.Handler.
func (r *Reloader) Handler() websocket.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {
		r.listeners = append(r.listeners, conn)
		<-r.closeSockets
	})
}

func (r *Reloader) Reload() {
	for _, conn := range r.listeners {
		conn.Write([]byte("refresh"))
	}
}

// Watch starts watching the filesystem for changes asynchronously, building any
// changes based on the contents of Builders and then reloading any connected
// browser.
//
// Watch returns once the goroutine has been spun off successfully. Any errors
// enountered while watching the filesystem are printed to stderr.
func (r *Reloader) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("cannot initialize fsnotify watcher: %s", err)
	}

	r.watcher = watcher
	r.stop = make(chan struct{})

	if err := filepath.WalkDir(
		".",
		func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			for pattern := range r.Ignore {
				if ok, err := filepath.Match(pattern, path); err != nil {
					return err
				} else if ok {
					return filepath.SkipDir
				}
			}
			if info.IsDir() {
				if err := r.watcher.Add(path); err != nil {
					return fmt.Errorf("can't watch %s: %w", path, err)
				}
				return nil
			}
			return nil
		},
	); err != nil {
		return err
	}

	go r.listen()

	return nil
}

func (r *Reloader) listen() {
	for {
		select {
		case event, ok := <-r.watcher.Events:
			if !ok {
				log.Println("fsnotify watcher closed")
				return
			}
			fmt.Println("event:", event)
			if err := r.Substructure.Rebuild(event.Name, dist); err != nil {
				log.Println(err.Error())
			}
			r.Reload()
		case err := <-r.watcher.Errors:
			if err != nil {
				panic(err)
			}
		case <-r.stop:
			r.watcher.Close()
			return
		}
	}
}

// ShutdownFunc returns a function that, when called, stops watching the
// filesystem for changes and gracefully closes any open WebSockets connections.
func (r *Reloader) ShutdownFunc() func() {
	return func() {
		for _, conn := range r.listeners {
			conn.Close()
		}
		r.listeners = []*websocket.Conn{}
		r.stop <- struct{}{}
		r.closeSockets <- struct{}{}
	}
}
