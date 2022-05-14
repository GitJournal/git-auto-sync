package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/GitJournal/git-auto-sync/common"
	cli "github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"
	git "gopkg.in/src-d/go-git.v4"
)

var errRepoPathInvalid = errors.New("Not a valid git repo")

func daemonStatus(ctx *cli.Context) error {
	// FIXME: Implement 'daemonStatus'

	// Print out the configuration
	// Print out uptime
	// Print out if there are any 'rebasing' issues and we are paused

	return nil
}

func daemonList(ctx *cli.Context) error {
	config, err := common.ReadConfig()
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
		// TODO: Check if the parent/ancestor is the repoPath!
	}

	repoPath, err := isValidGitRepo(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	config, err := common.ReadConfig()
	if err != nil {
		return tracerr.Wrap(err)
	}

	contains := false
	for _, rp := range config.Repos {
		if rp == repoPath {
			contains = true
			break
		}
	}

	if !contains {
		config.Repos = append(config.Repos, repoPath)
	}

	err = common.WriteConfig(config)
	if err != nil {
		return tracerr.Wrap(err)
	}

	s, err := common.NewService()
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = s.Enable()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func isValidGitRepo(repoPath string) (string, error) {
	info, err := os.Stat(repoPath)
	if os.IsNotExist(err) {
		return "", tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
	}

	if !info.IsDir() {
		return "", tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
	}

	_, err = git.PlainOpenWithOptions(repoPath, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return "", tracerr.Errorf("Not a valid git repo - %s\n%w", repoPath, err)
	}

	for true {
		_, err := os.Stat(repoPath)
		if err != nil {
			return "", tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
		}

		if repoPath == "." {
			return "", tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
		}

		repoPath = filepath.Dir(repoPath)
	}

	return repoPath, nil
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

	config, err := common.ReadConfig()
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
	err = common.WriteConfig(config)
	if err != nil {
		return tracerr.Wrap(err)
	}

	if len(config.Repos) == 0 {
		s, err := common.NewService()
		if err != nil {
			return tracerr.Wrap(err)
		}

		err = s.Disable()
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	return nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}
