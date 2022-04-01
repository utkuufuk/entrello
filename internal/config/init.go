package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "production" {
		godotenv.Load()
	}

	serializedServices := strings.Split(os.Getenv("SERVICES"), ",")
	if appEnv != "production" && os.Getenv("SERVICES") == "" {
		serializedServices = []string{}
	}

	services := make([]Service, 0, len(serializedServices))

	for _, service := range serializedServices {
		parts := strings.Split(service, "@")
		if len(parts) != 2 {
			panic(fmt.Sprintf("invalid service configuration string: %s", service))
		}
		services = append(services, Service{
			Label:    parts[0],
			Endpoint: parts[1],
		})
	}

	ServerCfg = ServerConfig{
		Port:                     os.Getenv("PORT"),
		Username:                 os.Getenv("USERNAME"),
		Password:                 os.Getenv("PASSWORD"),
		Services:                 services,
		TrelloApiKey:             os.Getenv("TRELLO_API_KEY"),
		TrelloApiToken:           os.Getenv("TRELLO_API_TOKEN"),
		TrelloBoardId:            os.Getenv("TRELLO_BOARD_ID"),
		TrelloSecret:             os.Getenv("TRELLO_SECRET"),
		TrelloWebhookCallbackUrl: os.Getenv("TRELLO_WEBHOOK_CALLBACK_URL"),
	}
}
