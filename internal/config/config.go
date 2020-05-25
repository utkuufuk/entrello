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

type GithubIssues struct {
	Enabled bool   `yaml:"enabled"`
	Period  Period `yaml:"period"`
	Token   string `yaml:"personal_access_token"`
	Label   string `yaml:"label_id"`
}

type TodoDock struct {
	Enabled  bool   `yaml:"enabled"`
	Period   Period `yaml:"period"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	Label    string `yaml:"label_id"`
}

type Sources struct {
	GithubIssues GithubIssues `yaml:"github_issues"`
	TodoDock     TodoDock     `yaml:"tododock"`
}

type Config struct {
	TrelloApiKey   string  `yaml:"trello_api_key"`
	TrelloApiToken string  `yaml:"trello_api_token"`
	BoardId        string  `yaml:"board_id"`
	ListId         string  `yaml:"list_id"`
	Sources        Sources `yaml:"sources"`
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
