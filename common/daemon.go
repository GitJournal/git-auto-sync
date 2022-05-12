package common

import (
	"fmt"
	"log"
	"sync"

	"github.com/kardianos/service"
)

var logger service.Logger

type Daemon struct{}

func (d *Daemon) Start(s service.Service) error {
	go d.run()
	return nil
}

func (d *Daemon) run() {
	config, err := ReadConfig()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, repoPath := range config.Repos {
		wg.Add(1)

		fmt.Println("Monitoring", repoPath)
		go watchForChanges(&wg, repoPath)
	}

	wg.Wait()
}

func (d *Daemon) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

// FIXME: pass some kind of channel which tells this when to close!
func watchForChanges(wg *sync.WaitGroup, repoPath string) {
	defer wg.Done()

	err := WatchForChanges(repoPath)
	if err != nil {
		log.Println(err)
	}
}
