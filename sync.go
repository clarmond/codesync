package main

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type syncAction int

const (
	actionCopy syncAction = iota
	actionDelete
)

// syncEvent is a single unit of work handed to the worker pool: mirror (or
// remove) one path from the watched source tree in the destination tree.
type syncEvent struct {
	path   string // absolute path under watchDir
	action syncAction
}

type syncer struct {
	watchDir string
	destDir  string
	ignore   *ignoreMatcher
}

// relPath returns ev's path relative to the watched root, using the OS
// separator, e.g. "src/app.js".
func (s *syncer) relPath(absPath string) string {
	rel, err := filepath.Rel(s.watchDir, absPath)
	if err != nil {
		return absPath
	}
	return rel
}

func (s *syncer) destPath(rel string) string {
	return filepath.Join(s.destDir, rel)
}

// process performs the copy or delete for a single event. It's safe to call
// concurrently from multiple worker goroutines since each event operates on
// an independent destination path.
func (s *syncer) process(ev syncEvent) {
	rel := s.relPath(ev.path)

	if s.ignore.ignores(rel) {
		logIgnore("Ignoring %s", rel)
		return
	}

	switch ev.action {
	case actionCopy:
		s.copy(ev.path, rel)
	case actionDelete:
		s.delete(rel)
	}
}

func (s *syncer) copy(src, rel string) {
	dst := s.destPath(rel)

	info, err := os.Lstat(src)
	if err != nil {
		// Path is already gone (e.g. a create immediately followed by a
		// delete) - nothing to sync.
		return
	}

	if info.IsDir() {
		if err := os.MkdirAll(dst, 0o755); err != nil {
			logFail("%s failed to copy: %v", rel, err)
			return
		}
		logSuccess("%s copied %s", rel, time.Now().Format(time.RFC1123))
		return
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		logFail("%s failed to copy: %v", rel, err)
		return
	}
	if err := copyFileContents(src, dst, info.Mode()); err != nil {
		logFail("%s failed to copy: %v", rel, err)
		return
	}
	logSuccess("%s copied %s", rel, time.Now().Format(time.RFC1123))

	cleanupMacSidecar(dst)
}

func (s *syncer) delete(rel string) {
	dst := s.destPath(rel)
	if err := os.RemoveAll(dst); err != nil {
		logFail("Failed to remove %s: %v", rel, err)
		return
	}
	logSuccess("%s removed %s", rel, time.Now().Format(time.RFC1123))
}

func copyFileContents(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// cleanupMacSidecar removes the "._<name>" AppleDouble file that macOS/SMB
// sometimes writes alongside a copied file on a network share. The write
// happens asynchronously after the copy, so this waits briefly before trying,
// mirroring the original Node implementation's 1s delay.
func cleanupMacSidecar(dst string) {
	dir := filepath.Dir(dst)
	sidecar := filepath.Join(dir, "._"+filepath.Base(dst))
	go func() {
		time.Sleep(1 * time.Second)
		os.Remove(sidecar)
	}()
}

// startWorkers launches n goroutines consuming events from the channel and
// bounds how many copies/deletes run concurrently against the destination.
func startWorkers(s *syncer, events <-chan syncEvent, n int) *sync.WaitGroup {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for ev := range events {
				s.process(ev)
			}
		}()
	}
	return &wg
}
