package common

import (
	"bytes"
	"os/exec"
	"sort"
	"strings"

	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
)

func commit(repoConfig RepoConfig) error {
	repoPath := repoConfig.RepoPath
	repo, err := git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return tracerr.Wrap(err)
	}

	w, err := repo.Worktree()
	if err != nil {
		return tracerr.Wrap(err)
	}

	status, err := w.Status()
	if err != nil {
		return tracerr.Wrap(err)
	}

	hasChanges := false
	commitMsg := []string{}
	for filePath, fileStatus := range status {
		if fileStatus.Worktree == git.Unmodified && fileStatus.Staging == git.Unmodified {
			continue
		}

		ignore, err := ShouldIgnoreFile(repoPath, filePath)
		if err != nil {
			return tracerr.Wrap(err)
		}

		if ignore {
			continue
		}

		hasChanges = true
		_, err = w.Add(filePath)
		if err != nil {
			return tracerr.Wrap(err)
		}

		msg := ""
		if fileStatus.Worktree == git.Untracked && fileStatus.Staging == git.Untracked {
			msg += "?? "
		} else {
			msg += " " + string(fileStatus.Worktree) + " "
		}
		msg += filePath
		commitMsg = append(commitMsg, msg)
	}

	sort.Strings(commitMsg)
	msg := strings.Join(commitMsg, "\n")

	if !hasChanges {
		return nil
	}

	_, err = GitCommand(repoConfig, []string{"commit", "-m", msg})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func GitCommand(repoConfig RepoConfig, args []string) (bytes.Buffer, error) {
	repoPath := repoConfig.RepoPath

	var outb, errb bytes.Buffer

	statusCmd := exec.Command("git", args...)
	statusCmd.Dir = repoPath
	statusCmd.Stdout = &outb
	statusCmd.Stderr = &errb
	err := statusCmd.Run()

	if err != nil {
		fullCmd := "git " + strings.Join(args, " ")
		err := tracerr.Errorf("%w: Command: %s\nStdOut: %s\nStdErr: %s", err, fullCmd, outb.String(), errb.String())
		return outb, err
	}
	return outb, nil
}
