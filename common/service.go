package common

import (
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kardianos/service"
	"github.com/ztrue/tracerr"
)

type Service struct {
	Service service.Service
}

type emptyDaemon struct{}

func (d emptyDaemon) Start(s service.Service) error {
	return nil
}

func (d emptyDaemon) Stop(s service.Service) error {
	return nil
}

func NewServiceWithDaemon(daemon service.Interface) (Service, error) {
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

	s, err := service.New(daemon, svcConfig)
	if err != nil {
		return Service{}, tracerr.Wrap(err)
	}

	return Service{Service: s}, nil
}

func NewService() (Service, error) {
	return NewServiceWithDaemon(emptyDaemon{})
}

func (srv Service) Enable() error {
	// TODO: Uninstall the old one, in case it was running
	//       Also stop the old service?

	err := srv.Service.Install()
	if err != nil {
		if strings.Contains(err.Error(), "Init already exists") {
			return nil
		}
		return tracerr.Wrap(err)
	}

	err = srv.Service.Restart()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (srv Service) Disable() error {
	err := srv.Service.Stop()
	if err != nil {
		return tracerr.Wrap(err)
	}

	err = srv.Service.Uninstall()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
