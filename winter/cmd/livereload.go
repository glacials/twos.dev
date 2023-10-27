package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

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
	// listeners is a mapping of WebSocket connections browsers have open with us
	// to the last time they were refreshed.
	listeners map[*websocket.Conn]time.Time
	stop      chan struct{}
	watcher   *fsnotify.Watcher
}

// Handler returns a function that handles incoming WebSocket connections which
// implements http.Handler.
func (r *Reloader) Handler() websocket.Handler {
	if r.listeners == nil {
		r.listeners = map[*websocket.Conn]time.Time{}
	}
	return websocket.Handler(func(conn *websocket.Conn) {
		r.listeners[conn] = time.Now()
		<-r.closeSockets
	})
}

// Reload notifies all connected browsers to reload the page.
func (r *Reloader) Reload() {
	for conn, lastRefreshed := range r.listeners {
		if time.Since(lastRefreshed) < 2*time.Second {
			// debounce
			continue
		}
		// If the browser gets our message, it'll close this via refresh. If not,
		// we'll want to prune it anyway.
		delete(r.listeners, conn)
		_, err := conn.Write([]byte("refresh"))
		if err != nil {
			log.Printf("can't refresh browser: %s", err.Error())
			if err := conn.Close(); err != nil {
				log.Printf("can't close browser connection: %s", err.Error())
			}
		}
	}
}

// Watch starts watching the filesystem for changes asynchronously, building any
// changes based on the contents of Builders and then reloading any connected
// browser.
//
// Watch returns once the goroutine has been spun off successfully. Any errors
// enountered while watching the filesystem are printed to stderr.
func (r *Reloader) Watch(paths []string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("cannot initialize fsnotify watcher: %s", err)
	}

	r.watcher = watcher
	r.stop = make(chan struct{})

	for _, path := range paths {
		for pattern := range r.Ignore {
			if ok, err := filepath.Match(pattern, path); err != nil {
				return err
			} else if ok {
				return filepath.SkipDir
			}
		}

		// If we're passed a file, watch it directly
		if stat, err := os.Stat(path); err != nil {
			return fmt.Errorf("cannot stat %s: %w", path, err)
		} else if !stat.IsDir() {
			if err := r.watcher.Add(path); err != nil {
				return fmt.Errorf("cannot add %q to watcher: %w", path, err)
			}
			continue
		}

		// If we're passed a directory, watch it recursively.
		// fsnotify.Watcher only supports watching directories
		// non-recurisvely, so we'll recurse ourselves.
		if err := filepath.WalkDir(
			path,
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
				if !info.IsDir() {
					// We can skip watching files, because we're watching all
					// directories.
					return nil
				}
				if err := r.watcher.Add(path); err != nil {
					return fmt.Errorf("can't watch %s: %w", path, err)
				}
				return nil
			},
		); err != nil {
			return err
		}
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
			if event.Op == fsnotify.Chmod {
				continue
			}
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

// Shutdown stops watching the filesystem for changes and gracefully closes any
// open WebSocket connections.
func (r *Reloader) Shutdown() {
	for conn := range r.listeners {
		conn.Close()
	}
	r.listeners = map[*websocket.Conn]time.Time{}
	r.stop <- struct{}{}
	r.closeSockets <- struct{}{}
}
