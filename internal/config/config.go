package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type TodoDock struct {
	Email    string `yaml:"email"`
	Password string `yaml:"password"`
}

type Sources struct {
	TodoDock TodoDock `yaml:"tododock"`
}

type Config struct {
	TrelloApiKey   string  `yaml:"trello_api_key"`
	TrelloApiToken string  `yaml:"trello_api_token"`
	Sources        Sources `yaml:"sources"`
}

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
