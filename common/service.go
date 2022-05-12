package common

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
	"github.com/ztrue/tracerr"
)

type Service struct {
	service service.Service
}

func NewService() (Service, error) {
	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["UserService"] = true
	options["RunAtLoad"] = true

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exDirPath := filepath.Dir(ex)
	executablePath := filepath.Join(exDirPath, "git-auto-sync-daemon")

	svcConfig := &service.Config{
		Name:        "git-auto-sync-daemon",
		DisplayName: "Git Auto Sync Daemon",
		Description: "Background Process for Auto Syncing Git Repos",
		UserName:    user.Username,

		Executable: executablePath,
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target syslog.target"},
		Option: options,
	}

	daemon := &Daemon{}
	s, err := service.New(daemon, svcConfig)
	if err != nil {
		return Service{}, tracerr.Wrap(err)
	}

	return Service{service: s}, nil
}

func (srv Service) Enable() error {
	// TODO: Uninstall the old one, in case it was running
	//       Also stop the old service?

	err := srv.service.Install()
	fmt.Println("Installed")
	if err != nil {
		if strings.Contains(err.Error(), "Init already exists") {
			return nil
		}
		return tracerr.Wrap(err)
	}

	err = srv.service.Restart()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (srv Service) Disable() error {
	err := srv.service.Uninstall()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
