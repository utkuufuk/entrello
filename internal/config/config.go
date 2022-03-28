package config

import (
	"encoding/json"
	"fmt"
	"os"
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

type RunnerConfig struct {
	TimezoneLocation string   `json:"timezone_location"`
	Trello           Trello   `json:"trello"`
	Sources          []Source `json:"sources"`
}

type ServerConfig struct {
	Port                     string
	Username                 string
	Password                 string
	TrelloApiKey             string
	TrelloApiToken           string
	TrelloBoardId            string
	TrelloSecret             string
	TrelloWebhookCallbackUrl string
}

const (
	PeriodTypeDefault = "default"
	PeriodTypeDay     = "day"
	PeriodTypeHour    = "hour"
	PeriodTypeMinute  = "minute"
)

var ServerCfg ServerConfig

func init() {
	ServerCfg = ServerConfig{
		Port:                     os.Getenv("PORT"),
		Username:                 os.Getenv("USERNAME"),
		Password:                 os.Getenv("PASSWORD"),
		TrelloApiKey:             os.Getenv("TRELLO_API_KEY"),
		TrelloApiToken:           os.Getenv("TRELLO_API_TOKEN"),
		TrelloBoardId:            os.Getenv("TRELLO_BOARD_ID"),
		TrelloSecret:             os.Getenv("TRELLO_SECRET"),
		TrelloWebhookCallbackUrl: os.Getenv("TRELLO_WEBHOOK_CALLBACK_URL"),
	}
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
