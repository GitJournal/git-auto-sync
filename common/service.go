package common

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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
	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["UserService"] = true
	options["RunAtLoad"] = true

	ex, err := os.Executable()
	if err != nil {
		return Service{}, tracerr.Wrap(err)
	}
	exDirPath := filepath.Dir(ex)
	executablePath := filepath.Join(exDirPath, "git-auto-sync-daemon")

	deps := []string{}
	if runtime.GOOS == "linux" {
		deps = []string{"After=network-online.target syslog.target"}
	}

	svcConfig := &service.Config{
		Name:        "git-auto-sync-daemon",
		DisplayName: "Git Auto Sync Daemon",
		Description: "Background Process for Auto Syncing Git Repos",

		Executable:   executablePath,
		Dependencies: deps,
		Option:       options,
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
	s := srv.Service

	status, err := s.Status()
	if err != nil {
		if !strings.Contains(err.Error(), "the service is not installed") {
			return tracerr.Wrap(err)
		}
	}

	stopped := false
	if status == service.StatusRunning {
		err := s.Stop()
		if err != nil {
			return tracerr.Wrap(err)
		}
		stopped = true
	}

	err = s.Install()
	if err != nil {
		if strings.Contains(err.Error(), "Init already exists") {
			_ = s.Uninstall()
			_ = s.Install()
		} else {
			return tracerr.Wrap(err)
		}
	} else {
		fmt.Println("Installing git-auto-sync as a daemon")
	}

	if stopped {
		fmt.Println("Restarting git-auto-sync-daemon")
	} else {
		fmt.Println("Starting git-auto-sync-daemon")
	}

	err = s.Start()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (srv Service) Disable() error {
	fmt.Println("Stopping git-auto-sync-daemon")
	err := srv.Service.Stop()
	if err != nil {
		return tracerr.Wrap(err)
	}

	fmt.Println("Uninstalling git-auto-sync as a daemon")
	err = srv.Service.Uninstall()
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (srv Service) Status() error {
	status, err := srv.Service.Status()
	if err != nil {
		return tracerr.Wrap(err)
	}

	switch status {
	case service.StatusRunning:
		fmt.Println("git-auto-sync-daemon is Running!")
	case service.StatusStopped:
		fmt.Println("git-auto-sync-daemon is NOT Running!")
	case service.StatusUnknown:
	default:
		fmt.Println("git-auto-sync-daemon status is Unknown. How mysterious!")
	}

	return nil
}
