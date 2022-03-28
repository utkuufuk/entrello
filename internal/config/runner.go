package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type RunnerConfig struct {
	TimezoneLocation string   `json:"timezone_location"`
	Trello           Trello   `json:"trello"`
	Sources          []Source `json:"sources"`
}

func ReadRunnerConfig(fileName string) (cfg RunnerConfig, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return cfg, fmt.Errorf("could not open config file: %v", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, err
}
