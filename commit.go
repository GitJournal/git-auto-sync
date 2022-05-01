package main

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ztrue/tracerr"
)

func commit(repoPath string) error {
	err, outb := GitCommand(repoPath, []string{"status", "--porcelain"})

	if err == nil {
		if len(outb.Bytes()) == 0 {
			return nil
		}
	}

	err, _ = GitCommand(repoPath, []string{"add", "--all"})
	if err != nil {
		return tracerr.Wrap(err)
	}

	err, _ = GitCommand(repoPath, []string{"commit", "-m", outb.String()})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func GitCommand(repoPath string, args []string) (error, bytes.Buffer) {
	var outb, errb bytes.Buffer

	statusCmd := exec.Command("git", args...)
	statusCmd.Dir = repoPath
	statusCmd.Stdout = &outb
	statusCmd.Stderr = &errb
	err := statusCmd.Run()

	if err != nil {
		fullCmd := "git " + strings.Join(args, " ")
		return tracerr.Errorf("%w: Command: %s\nStdOut: %s\nStdErr: %s", err, fullCmd, outb.String(), errb.String()), outb
	}
	return nil, outb
}
