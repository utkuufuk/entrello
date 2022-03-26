package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/utkuufuk/entrello/internal/logger"
)

const (
	PERIOD_TYPE_DEFAULT = "default"
	PERIOD_TYPE_DAY     = "day"
	PERIOD_TYPE_HOUR    = "hour"
	PERIOD_TYPE_MINUTE  = "minute"
)

type Period struct {
	Type     string `json:"type"`
	Interval int    `json:"interval"`
}

type Source struct {
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Strict   bool   `json:"strict"`
	Label    string `json:"label_id"`
	List     string `json:"list_id"`
	Period   Period `json:"period"`
}

type Trello struct {
	ApiKey   string `json:"api_key"`
	ApiToken string `json:"api_token"`
	BoardId  string `json:"board_id"`
}

type Config struct {
	TimezoneLocation string   `json:"timezone_location"`
	Trello           Trello   `json:"trello"`
	Sources          []Source `json:"sources"`
}

func ReadConfig(fileName string) (cfg Config, err error) {
	jsonStr := os.Getenv("ENTRELLO_CONFIG")
	if jsonStr != "" {
		logger.Info("Attempting to read configuration from environment variable 'ENTRELLO_CONFIG'")
		err = json.Unmarshal([]byte(jsonStr), &cfg)
		if err != nil {
			return cfg, fmt.Errorf("could not parse config stored in the environment variable: %s", err)
		}

		return cfg, nil
	}

	logger.Info("Attempting to read configuration from file '%s'", fileName)
	f, err := os.Open(fileName)
	if err != nil {
		return cfg, fmt.Errorf("could not open config file: %v", err)
	}
	defer f.Close()

	decoder := json.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, err
}
