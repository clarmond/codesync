package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

const debounceDelay = 150 * time.Millisecond

// fileWatcher recursively watches watchDir with fsnotify (which, unlike
// chokidar, only watches a single directory level per Add call) and turns
// raw filesystem events into debounced syncEvents on the events channel.
type fileWatcher struct {
	fsw      *fsnotify.Watcher
	watchDir string
	ignore   *ignoreMatcher
	events   chan<- syncEvent

	timersMu sync.Mutex
	timers   map[string]*time.Timer
}

func newFileWatcher(watchDir string, ignore *ignoreMatcher, events chan<- syncEvent) (*fileWatcher, error) {
	fsw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	w := &fileWatcher{
		fsw:      fsw,
		watchDir: watchDir,
		ignore:   ignore,
		events:   events,
		timers:   make(map[string]*time.Timer),
	}

	if err := w.addTree(watchDir); err != nil {
		fsw.Close()
		return nil, err
	}
	return w, nil
}

func (w *fileWatcher) relPath(absPath string) string {
	rel, err := filepath.Rel(w.watchDir, absPath)
	if err != nil {
		return absPath
	}
	return rel
}

// addTree adds an fsnotify watch for root and every non-ignored
// subdirectory beneath it, without descending into ignored subtrees at all
// (so a large ignored directory, e.g. node_modules/, never generates events
// in the first place).
func (w *fileWatcher) addTree(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || !d.IsDir() {
			return nil
		}
		rel := w.relPath(path)
		if rel != "." && w.ignore.ignores(rel) {
			return filepath.SkipDir
		}
		if err := w.fsw.Add(path); err != nil {
			logFail("Failed to watch %s: %v", rel, err)
		}
		return nil
	})
}

// rescanNewDir handles a directory that just appeared: it may already
// contain files (e.g. a git checkout or an extracted archive can create a
// whole subtree faster than individual watches can be attached), so it adds
// watches for every non-ignored subdirectory and enqueues a copy event for
// every non-ignored path found, including dir itself.
func (w *fileWatcher) rescanNewDir(dir string) {
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		rel := w.relPath(path)
		if w.ignore.ignores(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			if err := w.fsw.Add(path); err != nil {
				logFail("Failed to watch %s: %v", rel, err)
			}
		}
		w.enqueue(path, actionCopy)
		return nil
	})
}

func (w *fileWatcher) run() {
	for {
		select {
		case ev, ok := <-w.fsw.Events:
			if !ok {
				return
			}
			w.handle(ev)
		case err, ok := <-w.fsw.Errors:
			if !ok {
				return
			}
			logFail("Watcher error: %v", err)
		}
	}
}

func (w *fileWatcher) close() {
	w.fsw.Close()
}

func (w *fileWatcher) handle(ev fsnotify.Event) {
	switch {
	case ev.Op&(fsnotify.Remove|fsnotify.Rename) != 0:
		w.fsw.Remove(ev.Name) // best-effort; no-op if it wasn't a watched dir
		w.debounce(ev.Name, actionDelete)

	case ev.Op&fsnotify.Create != 0:
		info, err := os.Lstat(ev.Name)
		if err != nil {
			return // already gone
		}
		if info.IsDir() {
			w.rescanNewDir(ev.Name)
			return
		}
		w.debounce(ev.Name, actionCopy)

	case ev.Op&fsnotify.Write != 0:
		w.debounce(ev.Name, actionCopy)
	}
}

// debounce coalesces rapid repeat events for the same path (common with
// editor saves and git operations that touch a file more than once in quick
// succession) into a single enqueued syncEvent.
func (w *fileWatcher) debounce(path string, action syncAction) {
	w.timersMu.Lock()
	defer w.timersMu.Unlock()

	if t, ok := w.timers[path]; ok {
		t.Stop()
	}
	w.timers[path] = time.AfterFunc(debounceDelay, func() {
		w.timersMu.Lock()
		delete(w.timers, path)
		w.timersMu.Unlock()
		w.enqueue(path, action)
	})
}

func (w *fileWatcher) enqueue(path string, action syncAction) {
	w.events <- syncEvent{path: path, action: action}
}
