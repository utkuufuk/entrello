package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	PERIOD_TYPE_DEFAULT = "default"
	PERIOD_TYPE_DAY     = "day"
	PERIOD_TYPE_HOUR    = "hour"
	PERIOD_TYPE_MINUTE  = "minute"
)

type Period struct {
	Type     string `yaml:"type"`
	Interval int    `yaml:"interval"`
}

type Source struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Strict   bool   `yaml:"strict"`
	Label    string `yaml:"label_id"`
	List     string `yaml:"list_id"`
	Period   Period `yaml:"period"`
}

type Sources struct {
	GithubIssues Source `yaml:"github_issues"`
	TodoDock     Source `yaml:"tododock"`
	Habits       Source `yaml:"habits"`
}

type Trello struct {
	ApiKey   string `yaml:"api_key"`
	ApiToken string `yaml:"api_token"`
	BoardId  string `yaml:"board_id"`
}

type Telegram struct {
	Enabled bool   `yaml:"enabled"`
	Token   string `yaml:"token"`
	ChatId  int64  `yaml:"chat_id"`
}

type Config struct {
	TimezoneLocation string   `yaml:"timezone_location"`
	TimeoutSeconds   int      `yaml:"timeout_secs"`
	Trello           Trello   `yaml:"trello"`
	Sources          Sources  `yaml:"sources"`
	Telegram         Telegram `yaml:"telegram"`
}

// ReadConfig reads the YAML config file & decodes all parameters
func ReadConfig(fileName string) (cfg Config, err error) {
	f, err := os.Open(fileName)
	if err != nil {
		return cfg, fmt.Errorf("could not open config file: %v", err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, err
}
