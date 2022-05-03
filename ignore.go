package main

import (
	"path/filepath"
	"strings"
)

func shouldIgnoreFile(path string) bool {
	fileName := filepath.Base(path)
	return strings.HasSuffix(fileName, ".swp") || // vim
		strings.HasPrefix(path, "~") || // emacs
		strings.HasSuffix(path, "~") || // kate
		strings.HasPrefix(path, ".") // hidden files

	// Do not automatically ignore all hidden files, make this configurable
	// Also check if ignored by .gitignore
}
