package main

import (
	"fmt"
	"log"
	"os"

	cli "github.com/urfave/cli/v2"
	"github.com/ztrue/tracerr"
)

func main() {
	app := &cli.App{
		Name:  "git-auto-sync",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			cwd, err := os.Getwd()
			if err != nil {
				return tracerr.Wrap(err)
			}

			err = autoSync(cwd)
			if err != nil {
				return tracerr.Wrap(err)
			}
			fmt.Println("Done")
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
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
	err := commit(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = fetch(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// -> rebase if possible
	// -> revert if rebase fails
	// -> do a merge
	// -> push the changes

	return nil
}
