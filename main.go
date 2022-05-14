package main

import (
	"log"
	"os"

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

					return common.WatchForChanges(common.NewRepoConfig(repoPath))
				},
			},
			{
				Name:    "sync",
				Aliases: []string{"s"},
				Usage:   "Sync a repo right now",
				Action: func(ctx *cli.Context) error {
					repoPath, err := os.Getwd()
					if err != nil {
						return tracerr.Wrap(err)
					}

					err = common.AutoSync(repoPath)
					if err != nil {
						return tracerr.Wrap(err)
					}

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
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
