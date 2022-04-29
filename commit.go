package main

import (
	"bytes"
	"os/exec"
)

func commit(repoPath string) error {
	var outb, errb bytes.Buffer

	statusCmd := exec.Command("git", "status", "--porcelain")
	statusCmd.Dir = repoPath
	statusCmd.Stdout = &outb
	statusCmd.Stderr = &errb
	err := statusCmd.Run()

	if err == nil {
		if len(outb.Bytes()) == 0 {
			return nil
		}
	}

	addCmd := exec.Command("git", "add", "--all")
	addCmd.Dir = repoPath
	err = addCmd.Run()
	if err != nil {
		return err
	}

	commitCmd := exec.Command("git", "commit", "-m", string(outb.Bytes()))
	commitCmd.Dir = repoPath
	err = commitCmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// FIXME: Better errors
//        * Add stacktraces
//        * Wrap errors
