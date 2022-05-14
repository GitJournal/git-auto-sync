package main

import (
	"log"
	"sync"

	"github.com/GitJournal/git-auto-sync/common"
	"github.com/kardianos/service"
)

var logger service.Logger

type Daemon struct{}

func (d *Daemon) Start(s service.Service) error {
	go d.run()
	return nil
}

func (d *Daemon) run() {
	config, err := common.ReadConfig()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	for _, repoPath := range config.Repos {
		wg.Add(1)

		logger.Info("Monitoring", repoPath)
		go watchForChanges(&wg, repoPath)
	}

	wg.Wait()
}

func (d *Daemon) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	daemon := Daemon{}
	autoSyncService, err := common.NewServiceWithDaemon(&daemon)
	if err != nil {
		log.Fatal(err)
	}

	s := autoSyncService.Service
	// FIXME: Figure out this logger buillshit!
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		_ = logger.Error(err)
	}
}

// FIXME: pass some kind of channel which tells this when to close!
func watchForChanges(wg *sync.WaitGroup, repoPath string) {
	defer wg.Done()

	cfg := common.NewRepoConfig(repoPath)
	err := common.WatchForChanges(cfg)
	if err != nil {
		log.Println(err)
	}
}

// FIXME: Handle operating system signal which tells it to reload the config
