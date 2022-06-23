package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	cli "github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"github.com/GitJournal/git-auto-sync/common"
)

var version = "dev"

func main() {
	app := &cli.App{
		Name:                 "git-auto-sync",
		Version:              version,
		Usage:                "Automatically Sync any Git Repo",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "watch",
				Aliases: []string{"monitor", "w", "m"},
				Usage:   "Monitor a folder for changes",
				Action: func(ctx *cli.Context) error {
					repoPath, err := os.Getwd()
					if err != nil {
						return tracerr.Wrap(err)
					}

					repoPath, err = isValidGitRepo(repoPath)
					if err != nil {
						return tracerr.Wrap(err)
					}

					cfg, err := common.NewRepoConfig(repoPath)
					if err != nil {
						return tracerr.Wrap(err)
					}

					return common.WatchForChanges(cfg)
				},
			},
			{
				Name:    "sync",
				Aliases: []string{"s"},
				Usage:   "Sync a repo right now",
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "env",
						Aliases: []string{"e"},
						Usage:   "Env variables to pass",
					},
				},
				Action: func(ctx *cli.Context) error {
					repoPath, err := os.Getwd()
					if err != nil {
						return tracerr.Wrap(err)
					}

					repoPath, err = isValidGitRepo(repoPath)
					if err != nil {
						return tracerr.Wrap(err)
					}

					cfg, err := common.NewRepoConfig(repoPath)
					if err != nil {
						return tracerr.Wrap(err)
					}

					for _, e := range ctx.StringSlice("env") {
						cfg.Env = append(cfg.Env, e)
					}

					err = common.AutoSync(cfg)
					if err != nil {
						return tracerr.Wrap(err)
					}

					return nil
				},
			},
			{
				Name:  "check",
				Usage: "Check if a file will be ignored",
				Action: func(ctx *cli.Context) error {
					repoPath, err := os.Getwd()
					if err != nil {
						return tracerr.Wrap(err)
					}

					repoPath, err = isValidGitRepo(repoPath)
					if err != nil {
						return tracerr.Wrap(err)
					}

					path := ctx.Args().First()
					path, err = filepath.Abs(path)
					if err != nil {
						return tracerr.Wrap(err)
					}

					ignored, err := common.ShouldIgnoreFile(repoPath, path)
					if err != nil {
						return tracerr.Wrap(err)
					}
					fmt.Println("Ignored:", ignored)

					return nil
				},
			},
			{
				Name:  "version",
				Usage: "Print the version",
				Action: func(ctx *cli.Context) error {
					cli.VersionPrinter(ctx)
					return nil
				},
			},
			{
				Name:    "daemon",
				Aliases: []string{"d"},
				Usage:   "Interact with the background daemon",
				Subcommands: []*cli.Command{
					{
						Name:   "status",
						Usage:  "Show the Daemon's status",
						Action: daemonStatus,
					},
					{
						Name:    "list",
						Aliases: []string{"ls"},
						Usage:   "List of repos being auto-synced",
						Action:  daemonList,
					},
					{
						Name:   "add",
						Usage:  "Add a repo for auto-sync",
						Action: daemonAdd,
					},
					{
						Name:    "remove",
						Aliases: []string{"rm"},
						Usage:   "Remove a repo from auto-sync",
						Action:  daemonRm,
					},
					{
						Name:   "env",
						Usage:  "Set an environment variable",
						Action: daemonEnv,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
