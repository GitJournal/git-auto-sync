package main

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	fmt.Println("Hello")
}

type Config struct {
	repoPath     string
	pollInterval time.Duration
}

// what remotes?
// what branches?

func poll() {
	fmt.Println("Poll")
}

func autoSync(repoPath string) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		log.Fatal(err)
	}

	w, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	err = w.AddGlob("*")
	if err != nil {
		log.Fatal(err)
	}

	status, err := w.Status()
	if err != nil {
		log.Fatal(err)
	}

	status.IsClean()

	// commit everything

	// do a fetch
	// -> rebase if possible
	// -> revert if rebase fails
	// -> do a merge
	// -> push the changes

	// write tests for all of this
}

// Commit Tests
// 1. Empty Repo
// 2. One file changes
// 3. Multiple files changed
// 4. A file in a subdirectory changed
