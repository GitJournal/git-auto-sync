package main

import (
	"fmt"
	"log"
	"sync"

	"os/user"

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
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["UserService"] = true
	options["RunAtLoad"] = true

	svcConfig := &service.Config{
		Name:        "GitAutoSyncDaemon",
		DisplayName: "Git Auto Sync Daemon",
		Description: "Background Process for Auto Syncing Git Repos",
		UserName:    user.Username,

		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	daemon := &Daemon{}
	s, err := service.New(daemon, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(s.Status())
	fmt.Println(s.Install())
	fmt.Println(s.Status())

	err = s.Run()
	if err != nil {
		_ = logger.Error(err)
	}
}

// FIXME: pass some kind of channel which tells this when to close!
func watchForChanges(wg *sync.WaitGroup, repoPath string) {
	defer wg.Done()

	err := common.WatchForChanges(repoPath)
	if err != nil {
		log.Println(err)
	}
}
