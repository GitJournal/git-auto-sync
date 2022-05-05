package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cli "github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
)

var errRepoPathInvalid = errors.New("Not a valid git repo")

func daemonStatus(ctx *cli.Context) error {
	return nil
}

func daemonList(ctx *cli.Context) error {
	config, err := readConfig()
	if err != nil {
		return tracerr.Wrap(err)
	}

	for _, repoPath := range config.Repos {
		fmt.Println(repoPath)
	}
	return nil
}

func daemonAdd(ctx *cli.Context) error {
	repoPath := ctx.Args().First()
	if !strings.HasPrefix(repoPath, string(filepath.Separator)) {
		cwd, err := os.Getwd()
		if err != nil {
			return tracerr.Wrap(err)
		}

		repoPath = filepath.Join(cwd, repoPath)
	}

	err := isValidGitRepo(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	config, err := readConfig()
	if err != nil {
		return tracerr.Wrap(err)
	}

	config.Repos = append(config.Repos, repoPath)
	err = writeConfig(config)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func isValidGitRepo(repoPath string) error {
	info, err := os.Stat(repoPath)
	if os.IsNotExist(err) {
		return tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
	}

	if !info.IsDir() {
		return tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
	}

	_, err = git.PlainOpen(repoPath)
	if err != nil {
		return tracerr.Errorf("%w - %s\n%w", errRepoPathInvalid, repoPath, err)
	}

	return nil
}

func daemonRm(ctx *cli.Context) error {
	repoPath := ctx.Args().First()
	if !strings.HasPrefix(repoPath, string(filepath.Separator)) {
		cwd, err := os.Getwd()
		if err != nil {
			return tracerr.Wrap(err)
		}

		repoPath = filepath.Join(cwd, repoPath)
	}

	config, err := readConfig()
	if err != nil {
		return tracerr.Wrap(err)
	}

	pos := -1
	for i, rp := range config.Repos {
		if rp == repoPath {
			pos = i
			break
		}
	}

	if pos == -1 {
		err = errors.New("Repo Not tracked")
		return tracerr.Errorf("%w - %s", err, repoPath)
	}

	config.Repos = remove(config.Repos, pos)
	err = writeConfig(config)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
