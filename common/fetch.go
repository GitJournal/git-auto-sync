package common

import (
	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
)

func fetch(repoConfig RepoConfig) error {
	repoPath := repoConfig.RepoPath
	r, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return tracerr.Wrap(err)
	}

	remotes, err := r.Remotes()
	if err != nil {
		return tracerr.Wrap(err)
	}

	for _, remote := range remotes {
		remoteName := remote.Config().Name

		_, err := GitCommand(repoConfig, []string{"fetch", remoteName})
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	return nil
}
