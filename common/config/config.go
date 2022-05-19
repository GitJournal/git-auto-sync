package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kirsle/configdir"
	"github.com/ztrue/tracerr"
)

type Config struct {
	Repos []string `json:"repos"`
}

func Read() (*Config, error) {
	configPath := configdir.LocalConfig("git-auto-sync")
	err := configdir.MakePath(configPath)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	configFile := filepath.Join(configPath, "config.json")
	config := Config{}

	if _, err = os.Stat(configFile); os.IsNotExist(err) {
		return &config, nil
	} else {
		// Load the existing file.
		fh, err := os.Open(configFile)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
		defer fh.Close()

		decoder := json.NewDecoder(fh)
		err = decoder.Decode(&config)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}
	}

	return &config, nil
}

func Write(config *Config) error {
	configPath := configdir.LocalConfig("git-auto-sync")
	err := configdir.MakePath(configPath)
	if err != nil {
		return tracerr.Wrap(err)
	}

	configFile := filepath.Join(configPath, "config.json")

	fh, err := os.Create(configFile)
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer fh.Close()

	encoder := json.NewEncoder(fh)
	err = encoder.Encode(&config)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
