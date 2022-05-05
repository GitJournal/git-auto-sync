package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/GitJournal/git-auto-sync/common"
)

func main() {
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

func watchForChanges(wg *sync.WaitGroup, repoPath string) {
	defer wg.Done()

	err := common.WatchForChanges(repoPath)
	if err != nil {
		log.Println(err)
	}
}
