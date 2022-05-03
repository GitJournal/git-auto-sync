package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	cli "github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"

	"path/filepath"

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
					cwd, err := os.Getwd()
					if err != nil {
						return tracerr.Wrap(err)
					}

					err = autoSync(cwd)
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func watchForChanges(ctx *cli.Context) error {
	notifyChannel := make(chan notify.EventInfo, 100)

	err := notify.Watch("./...", notifyChannel, notify.Write, notify.Rename, notify.Remove)
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
		if shouldIgnoreFile(ei.Path()) {
			continue
		}

		// Wait for 'x' seconds
		log.Println("Got event:", ei)
		notifyFilteredChannel <- ei
	}
}

// type Config struct {
// 	repoPath     string
// 	pollInterval time.Duration
// }

// what remotes?
// what branches?

// func poll() {
// 	fmt.Println("Poll")
// }

func autoSync(repoPath string) error {
	var err error
	err = ensureGitAuthor(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = commit(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = fetch(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = rebase(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = push(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// -> rebase if possible
	// -> revert if rebase fails
	// -> do a merge
	// -> push the changes

	return nil
}

func shouldIgnoreFile(path string) bool {
	fileName := filepath.Base(path)
	return strings.HasSuffix(fileName, ".swp") || // vim
		strings.HasPrefix(path, "~") || // emacs
		strings.HasSuffix(path, "~") || // kate
		strings.HasPrefix(path, ".") // hidden files

	// Do not automatically ignore all hidden files, make this configurable
	// Also check if ignored by .gitignore
}
