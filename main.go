package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	cli "github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"github.com/rjeczalik/notify"
)

func main() {
	app := &cli.App{
		Name:  "git-auto-sync",
		Usage: "fight the loneliness!",
		Commands: []*cli.Command{
			{
				Name:   "watch",
				Usage:  "Watch a folder for changes",
				Action: watchForChanges,
			},
			{
				Name:  "sync",
				Usage: "Sync a repo right now",
				Action: func(ctx *cli.Context) error {
					repoPath, err := os.Getwd()
					if err != nil {
						return tracerr.Wrap(err)
					}

					err = autoSync(repoPath)
					if err != nil {
						return tracerr.Wrap(err)
					}

					return nil
				},
			},
			{
				Name:  "notify",
				Usage: "Sync a repo right now",
				Action: func(ctx *cli.Context) error {
					err := beeep.Alert("Title", "Message body", "assets/warning.png")
					if err != nil {
						panic(err)
					}

					return nil
				},
			},
			{
				Name:  "daemon",
				Usage: "Interact with the background daemon",
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

func watchForChanges(ctx *cli.Context) error {
	repoPath, err := os.Getwd()
	if err != nil {
		return tracerr.Wrap(err)
	}

	notifyChannel := make(chan notify.EventInfo, 100)

	err = notify.Watch("./...", notifyChannel, notify.Write, notify.Rename, notify.Remove)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer notify.Stop(notifyChannel)

	notifyFilteredChannel := make(chan notify.EventInfo, 100)
	ticker := time.NewTicker(20 * time.Second)

	go func() {
		var events []notify.EventInfo
		for {
			select {
			case ei := <-notifyFilteredChannel:
				events = append(events, ei)
				continue

			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				if len(events) != 0 {
					fmt.Println("Committing")
					events = []notify.EventInfo{}
				}
			}
		}
	}()

	// Block until an event is received.
	for {
		ei := <-notifyChannel
		ignore, err := shouldIgnoreFile(repoPath, ei.Path())
		if err != nil {
			return tracerr.Wrap(err)
		}
		if ignore {
			continue
		}

		// Wait for 'x' seconds
		log.Println("Got event:", ei)
		notifyFilteredChannel <- ei
	}
}

// func poll() {
// 	fmt.Println("Poll")
// }
