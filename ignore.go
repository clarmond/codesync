package main

import (
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

const ignoreFileName = ".csignore"

// ignoreMatcher decides whether a path relative to the watched root should be
// skipped: dotfiles/dot-directories are always skipped, plus anything matched
// by the optional .csignore file (gitignore syntax) at the root of the watched
// directory.
type ignoreMatcher struct {
	gi *gitignore.GitIgnore
}

func loadIgnoreMatcher(watchDir string) *ignoreMatcher {
	m := &ignoreMatcher{}

	path := filepath.Join(watchDir, ignoreFileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return m
	}

	m.gi = gitignore.CompileIgnoreLines(strings.Split(string(data), "\n")...)
	return m
}

// ignores reports whether relPath (relative to the watched root, using the
// OS path separator) should be skipped. relPath's own file/directory-ness
// isn't known here, so directory-only patterns (e.g. "node_modules/") are
// matched by also trying the path with a trailing slash - gitignore trailing
// slash patterns only match a literal trailing slash in the tested string.
func (m *ignoreMatcher) ignores(relPath string) bool {
	if hasDotSegment(relPath) {
		return true
	}
	if m.gi == nil {
		return false
	}
	slashPath := filepath.ToSlash(relPath)
	return m.gi.MatchesPath(slashPath) || m.gi.MatchesPath(slashPath+"/")
}

// hasDotSegment reports whether any path segment starts with a dot, matching
// the Node version's blanket dotfile/dot-directory ignore.
func hasDotSegment(relPath string) bool {
	for _, part := range strings.Split(filepath.ToSlash(relPath), "/") {
		if strings.HasPrefix(part, ".") && part != "" {
			return true
		}
	}
	return false
}
