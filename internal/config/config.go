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

type SourceConfig struct {
	Name    string `yaml:"name"`
	Enabled bool   `yaml:"enabled"`
	Strict  bool   `yaml:"strict"`
	Label   string `yaml:"label_id"`
	Period  Period `yaml:"period"`
}

type GithubIssues struct {
	SourceConfig SourceConfig `yaml:"source_config"`
	Token        string       `yaml:"personal_access_token"`
}

type TodoDock struct {
	SourceConfig SourceConfig `yaml:"source_config"`
	Email        string       `yaml:"email"`
	Password     string       `yaml:"password"`
}

type Habits struct {
	SourceConfig    SourceConfig `yaml:"source_config"`
	SpreadsheetId   string       `yaml:"spreadsheet_id"`
	CredentialsFile string       `yaml:"credentials_file"`
	TokenFile       string       `yaml:"token_file"`
}

type Sources struct {
	GithubIssues GithubIssues `yaml:"github_issues"`
	TodoDock     TodoDock     `yaml:"tododock"`
	Habits       Habits       `yaml:"habits"`
}

type Trello struct {
	ApiKey      string `yaml:"api_key"`
	ApiToken    string `yaml:"api_token"`
	BoardId     string `yaml:"board_id"`
	TodoListId  string `yaml:"todo_list_id"`
	TodayListId string `yaml:"today_list_id"`
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
