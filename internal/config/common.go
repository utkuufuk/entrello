package config

type Trello struct {
	ApiKey   string `json:"api_key"`
	ApiToken string `json:"api_token"`
	BoardId  string `json:"board_id"`
}
