package main

import (
	"errors"
	"os/exec"

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
		var exerr *exec.ExitError
		if errors.As(err, &exerr) && exerr.ExitCode() == 1 {
			err, _ := GitCommand(repoPath, []string{"rebase", "--abort"})
			if err != nil {
				return tracerr.Wrap(err)
			}
			return nil
		}
		return tracerr.Wrap(err)
	}

	return nil
}

// fixme; FIgure out a way to programatically detect if a rebase is going on
// fixme: See if the exit code is 1 when a rebase can fail?
//        how else can a rebase fail?
// fixme: Return a proper error if a rebase fails!
