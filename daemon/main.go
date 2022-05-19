package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/GitJournal/git-auto-sync/common"
	cfg "github.com/GitJournal/git-auto-sync/common/config"
	"github.com/kardianos/service"
)

type Daemon struct{}

func (d *Daemon) Start(s service.Service) error {
	go d.run()
	return nil
}

func (d *Daemon) run() {
	config, err := cfg.Read()
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

func main() {
	daemon := Daemon{}
	autoSyncService, err := common.NewServiceWithDaemon(&daemon)
	if err != nil {
		log.Fatal("BuildService", err)
	}

	s := autoSyncService.Service
	logger, err := s.Logger(nil)
	if err != nil {
		log.Fatal("BuildLogger", err)
	}

	err = s.Run()
	if err != nil {
		_ = logger.Error("RunService", err)
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
