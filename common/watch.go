package common

import (
	"log"
	"path/filepath"
	"time"

	"github.com/rjeczalik/notify"
	"github.com/ztrue/tracerr"
)

// FIXME: Watch for suspend / resume
// FIXME: Replace the logger with returning an error and retrying after 'x' minutes

type RepoConfig struct {
	RepoPath     string
	PollInterval time.Duration
	FSLag        time.Duration
}

func NewRepoConfig(repoPath string) RepoConfig {
	return RepoConfig{
		RepoPath:     repoPath,
		PollInterval: 10 * time.Minute,
		FSLag:        1 * time.Second,
	}
}

func WatchForChanges(cfg RepoConfig) error {
	repoPath := cfg.RepoPath
	var err error

	err = AutoSync(cfg.RepoPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	notifyFilteredChannel := make(chan notify.EventInfo, 100)
	pollTicker := time.NewTicker(cfg.PollInterval)

	// Filtered events
	go func() {
		for {
			select {
			case <-notifyFilteredChannel:
				// Wait 1 second
				timer1 := time.NewTimer(cfg.FSLag)
				done := make(chan bool)
				go func() {
					<-timer1.C
					done <- true
				}()

				err := AutoSync(repoPath)
				if err != nil {
					log.Fatalln(err)
				}
				continue

			case <-pollTicker.C:
				err := AutoSync(repoPath)
				if err != nil {
					log.Fatalln(err)
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
