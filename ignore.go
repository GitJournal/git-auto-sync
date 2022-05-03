package main

import (
	"path/filepath"
	"strings"

	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
)

func shouldIgnoreFile(repoPath string, path string) (bool, error) {
	fileName := filepath.Base(path)
	var isTempFile = strings.HasSuffix(fileName, ".swp") || // vim
		strings.HasPrefix(path, "~") || // emacs
		strings.HasSuffix(path, "~") || // kate
		strings.HasPrefix(path, ".") // hidden files

	if isTempFile {
		return true, nil
	}

	// FIXME: Do not automatically ignore all hidden files, make this configurable

	return isFileIgnoredByGit(repoPath, path)
}

func isFileIgnoredByGit(repoPath string, filePath string) (bool, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	patterns, err := gitignore.ReadPatterns(w.Filesystem, nil)
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	patterns = append(patterns, w.Excludes...)
	m := gitignore.NewMatcher(patterns)

	return m.Match([]string{filePath}, false), err
}
