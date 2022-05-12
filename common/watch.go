package common

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/ztrue/tracerr"
)

// FIXME: Watch for suspend / resume
// FIXME: Properly log the errors
// FIXME: Properly log each time this performs a sync

func WatchForChanges(repoPath string) error {
	var err error

	err = AutoSync(repoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	notifyFilteredChannel := make(chan notify.EventInfo, 100)
	ticker := time.NewTicker(20 * time.Second)

	// Filtered events
	go func() {
		var events []notify.EventInfo
		for {
			select {
			case ei := <-notifyFilteredChannel:
				events = append(events, ei)
				continue

			case <-ticker.C:
				if len(events) != 0 {
					events = []notify.EventInfo{}

					fmt.Println("Committing")
					err := AutoSync(repoPath)
					if err != nil {
						log.Fatalln(err)
					}
				}
			}
		}
	}()

	//
	// Watch for FS events
	//
	notifyChannel := make(chan notify.EventInfo, 100)

	err = notify.Watch(filepath.Join(repoPath, "..."), notifyChannel, notify.Write, notify.Rename, notify.Remove, notify.Create)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer notify.Stop(notifyChannel)

	for {
		ei := <-notifyChannel
		ignore, err := shouldIgnoreFile(repoPath, ei.Path())
		if err != nil {
			return tracerr.Wrap(err)
		}
		if ignore {
			continue
		}

		notifyFilteredChannel <- ei
	}
}
