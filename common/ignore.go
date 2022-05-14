package common

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/format/gitignore"
)

func ShouldIgnoreFile(repoPath string, fullFilePath string) (bool, error) {
	fileName := filepath.Base(fullFilePath)
	var isTempFile = strings.HasSuffix(fileName, ".swp") || // vim
		strings.HasPrefix(fileName, "~") || // emacs
		strings.HasSuffix(fileName, "~") || // kate
		strings.HasPrefix(fileName, ".") // hidden files

	// FIXME: Do not automatically ignore all hidden files, make this configurable

	if isTempFile {
		return true, nil
	}

	relativePath := fullFilePath[len(repoPath)+1:]
	if strings.HasPrefix(relativePath, ".git/") {
		return true, nil
	}

	empty, err := isEmptyFile(fullFilePath)
	if err != nil {
		return false, tracerr.Wrap(err)
	}
	if empty {
		return true, nil
	}

	return isFileIgnoredByGit(repoPath, fullFilePath)
}

func isFileIgnoredByGit(repoPath string, filePath string) (bool, error) {
	repo, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
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

func isEmptyFile(filePath string) (bool, error) {
	stat, err := os.Stat(filePath)

	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return stat.Size() == 0, nil
}
