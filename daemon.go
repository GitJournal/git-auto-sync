package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/GitJournal/git-auto-sync/common"
	cfg "github.com/GitJournal/git-auto-sync/common/config"
	cli "github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"
	"golang.org/x/exp/slices"
	git "gopkg.in/src-d/go-git.v4"
)

var errRepoPathInvalid = errors.New("Not a valid git repo")

func daemonStatus(ctx *cli.Context) error {
	s, err := common.NewService()
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = s.Status()
	if err != nil {
		return tracerr.Wrap(err)
	}

	config, err := cfg.Read()
	if err != nil {
		return tracerr.Wrap(err)
	}

	fmt.Println("Monitoring - ")
	for _, repoPath := range config.Repos {
		fmt.Println("  ", repoPath)
	}

	// FIXME: Print out if there are any 'rebasing' issues and we are paused

	return nil
}

func daemonList(ctx *cli.Context) error {
	config, err := cfg.Read()
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
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	repoPath, err = isValidGitRepo(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	config, err := cfg.Read()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if slices.Contains(config.Repos, repoPath) {
		fmt.Println("The Daemon is already monitoring " + repoPath)
	} else {
		config.Repos = append(config.Repos, repoPath)
	}

	err = cfg.Write(config)
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

	for {
		info, err := os.Stat(filepath.Join(repoPath, ".git"))
		if err != nil {
			if !os.IsNotExist(err) {
				return "", tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
			}
		}

		if os.IsNotExist(err) {
			repoPath = filepath.Dir(repoPath)
			continue
		}

		if !info.IsDir() {
			return "", tracerr.Errorf("%w - %s", errRepoPathInvalid, repoPath)
		}
		break
	}

	return repoPath, nil
}

func daemonRm(ctx *cli.Context) error {
	repoPath := ctx.Args().First()
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	repoPath, err = isValidGitRepo(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	config, err := cfg.Read()
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
	err = cfg.Write(config)
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

func daemonEnv(ctx *cli.Context) error {
	vars := ctx.Args().Slice()

	for _, v := range vars {
		if !strings.Contains(v, "=") {
			log.Fatalln("Env variables must be in the format 'key=value'")
		}
	}

	config, err := cfg.Read()
	if err != nil {
		return tracerr.Wrap(err)
	}

	envMap := toEnvMap(config.Envs)
	newMap := toEnvMap(vars)

	for k, v := range newMap {
		envMap[k] = v
	}

	config.Envs = toEnvStrings(envMap)
	err = cfg.Write(config)
	if err != nil {
		return tracerr.Wrap(err)
	}

	fmt.Println(strings.Join(config.Envs, "\n"))

	return nil
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func toEnvMap(envs []string) map[string]string {
	m := map[string]string{}
	for _, e := range envs {
		parts := strings.Split(e, "=")
		if len(parts) > 1 {
			m[parts[0]] = strings.Join(parts[1:], "=")
		} else {
			m[e] = ""
		}
	}

	return m
}

func toEnvStrings(m map[string]string) []string {
	vals := []string{}
	for k, v := range m {
		x := fmt.Sprintf("%s=%s", k, v)
		vals = append(vals, x)
	}

	return vals
}
