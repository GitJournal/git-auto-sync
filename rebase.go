package main

import (
	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func rebase(repoPath string) error {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	config, err := repo.Config()
	if err != nil {
		return tracerr.Wrap(err)
	}

	ref, err := repo.Reference(plumbing.HEAD, false)
	if err != nil {
		return tracerr.Wrap(err)
	}

	currentBranchName := ref.Target().Short()
	branchConfig := config.Branches[currentBranchName]
	if branchConfig == nil {
		// No tracking branch, nothing to do
		return nil
	}

	remoteName := branchConfig.Remote
	remoteBranchName := branchConfig.Merge.Short()

	err, _ = GitCommand(repoPath, []string{"rebase", remoteName + "/" + remoteBranchName})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
