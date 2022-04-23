package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()

	services, err := parseServices(os.Getenv("SERVICES"))
	if err != nil {
		fmt.Println("Could not parse the environment variable 'SERVICES':", err)
		os.Exit(1)
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
