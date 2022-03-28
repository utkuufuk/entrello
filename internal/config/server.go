package config

import "os"

type ServerConfig struct {
	Port                     string
	Username                 string
	Password                 string
	TrelloWebhookCallbackUrl string
	TrelloSecret             string
}

var ServerCfg ServerConfig

func init() {
	ServerCfg = ServerConfig{
		Port:                     os.Getenv("PORT"),
		Username:                 os.Getenv("USERNAME"),
		Password:                 os.Getenv("PASSWORD"),
		TrelloWebhookCallbackUrl: os.Getenv("TRELLO_WEBHOOK_CALLBACK_URL"),
		TrelloSecret:             os.Getenv("TRELLO_SECRET"),
	}
}
