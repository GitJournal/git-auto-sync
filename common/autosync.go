package common

import (
	"github.com/ztrue/tracerr"
)

func AutoSync(repoPath string) error {
	var err error
	err = ensureGitAuthor(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = commit(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = fetch(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = rebase(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = push(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// -> rebase if possible
	// -> revert if rebase fails
	// -> do a merge
	// -> push the changes

	return nil
}
