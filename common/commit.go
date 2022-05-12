package common

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/ztrue/tracerr"
)

func commit(repoPath string) error {
	outb, err := GitCommand(repoPath, []string{"status", "--porcelain"})

	if err == nil {
		if len(outb.Bytes()) == 0 {
			return nil
		}
	}

	_, err = GitCommand(repoPath, []string{"add", "--all"})
	if err != nil {
		return tracerr.Wrap(err)
	}

	_, err = GitCommand(repoPath, []string{"commit", "-m", outb.String()})
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func GitCommand(repoPath string, args []string) (bytes.Buffer, error) {
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
