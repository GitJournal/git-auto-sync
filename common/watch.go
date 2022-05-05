package common

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/ztrue/tracerr"
)

// FIXME: Poll for changes as well?
// FIXME: Watch for suspend / resume

func WatchForChanges(repoPath string) error {
	var err error
	notifyChannel := make(chan notify.EventInfo, 100)

	err = notify.Watch(filepath.Join(repoPath, "..."), notifyChannel, notify.Write, notify.Rename, notify.Remove)
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
