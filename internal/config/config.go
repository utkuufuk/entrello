package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type GithubIssues struct {
	Enabled bool `yaml:"enabled"`
	Token   string
	LabelId string `yaml:"label_id"`
}

type TodoDock struct {
	Enabled  bool   `yaml:"enabled"`
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
	LabelId  string `yaml:"label_id"`
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
	// open config file
	f, err := os.Open(fileName)
	if err != nil {
		return cfg, fmt.Errorf("could not open config file: %v", err)
	}
	defer f.Close()

	// decode config vars & return as a struct
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	return cfg, err
}
