package common

import (
	"errors"

	"github.com/gen2brain/beeep"
	"github.com/ztrue/tracerr"
)

// FIXME: Add logs for when we commit, pull, and push
func AutoSync(repoConfig RepoConfig) error {
	var err error
	err = ensureGitAuthor(repoConfig)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = commit(repoConfig)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = fetch(repoConfig)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = rebase(repoConfig)
	if err != nil {
		if errors.Is(err, errRebaseFailed) {
			repoPath := repoConfig.RepoPath
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

	err = push(repoConfig)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// -> do a merge
	// -> push the changes

	return nil
}
