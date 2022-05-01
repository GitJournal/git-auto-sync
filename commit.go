package main

import (
	"bytes"
	"os/exec"

	"github.com/pkg/errors"
)

func commit(repoPath string) error {
	err, outb, _ := GitCommand(repoPath, []string{"status", "--porcelain"})

	if err == nil {
		if len(outb.Bytes()) == 0 {
			return nil
		}
	}

	err, _, _ = GitCommand(repoPath, []string{"add", "--all"})
	if err != nil {
		return errors.Wrap(err, "Commit Failed")
	}

	err, _, _ = GitCommand(repoPath, []string{"commit", "-m", string(outb.Bytes())})
	if err != nil {
		return errors.Wrap(err, "Commit Failed")
	}

	return nil
}

func GitCommand(repoPath string, args []string) (error, bytes.Buffer, bytes.Buffer) {
	var outb, errb bytes.Buffer

	statusCmd := exec.Command("git", args...)
	statusCmd.Dir = repoPath
	statusCmd.Stdout = &outb
	statusCmd.Stderr = &errb
	err := statusCmd.Run()

	return errors.Wrap(err, "Git Command Failed"), outb, errb
}
