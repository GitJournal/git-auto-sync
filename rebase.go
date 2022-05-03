package main

import (
	"errors"
	"os"
	"os/exec"
	"path"

	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

func rebase(repoPath string) error {
	bi, err := fetchBranchInfo(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rebaseErr, _ := GitCommand(repoPath, []string{"rebase", bi.UpstreamRemote + "/" + bi.UpstreamBranch})
	if rebaseErr != nil {
		rebaseInProgress, err := isRebasing(repoPath)
		if err != nil {
			return tracerr.Wrap(err)
		}

		var exerr *exec.ExitError
		if errors.As(rebaseErr, &exerr) && exerr.ExitCode() == 1 && rebaseInProgress {
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

// fixme: See if the exit code is 1 when a rebase can fail?
//        how else can a rebase fail?
// fixme: Return a proper error if a rebase fails!

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

type branchInfo struct {
	CurrentBranch  string
	UpstreamRemote string
	UpstreamBranch string
}

func fetchBranchInfo(repoPath string) (branchInfo, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return branchInfo{}, tracerr.Wrap(err)
	}

	config, err := repo.Config()
	if err != nil {
		return branchInfo{}, tracerr.Wrap(err)
	}

	ref, err := repo.Reference(plumbing.HEAD, false)
	if err != nil {
		return branchInfo{}, tracerr.Wrap(err)
	}

	currentBranchName := ref.Target().Short()
	branchConfig := config.Branches[currentBranchName]
	if branchConfig == nil {
		// No tracking branch, nothing to do
		return branchInfo{CurrentBranch: currentBranchName}, nil
	}

	return branchInfo{
		CurrentBranch:  currentBranchName,
		UpstreamRemote: branchConfig.Remote,
		UpstreamBranch: branchConfig.Merge.Short(),
	}, nil
}

func isRebasing(repoPath string) (bool, error) {
	ra, err := exists(path.Join(repoPath, ".git", "rebase-apply"))
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	rm, err := exists(path.Join(repoPath, ".git", "rebase-merge"))
	if err != nil {
		return false, tracerr.Wrap(err)
	}

	return ra || rm, nil
}
