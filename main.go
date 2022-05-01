package main

import (
	"fmt"
	"time"

	"github.com/ztrue/tracerr"
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

func autoSync(repoPath string) error {
	err := commit(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = fetch(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// -> rebase if possible
	// -> revert if rebase fails
	// -> do a merge
	// -> push the changes

	return nil
}
