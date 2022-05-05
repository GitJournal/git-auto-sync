package common

import (
	"errors"

	"github.com/gen2brain/beeep"
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
		if errors.Is(err, errRebaseFailed) {
			err := beeep.Alert("Git Auto Sync - Conflict", "Could not rebase for - "+repoPath, "assets/warning.png")
			if err != nil {
				return tracerr.Wrap(err)
			}
		}
		// How should we continue?
		// - Keep sending the notification each time?
		// - Or something a bit better?
		return tracerr.Wrap(err)
	}

	err = push(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// -> do a merge
	// -> push the changes

	return nil
}
