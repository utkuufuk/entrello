package tododock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

const (
	BASE_URL = "https://tododock.com/api"
)

type TodoDockSource struct {
	params config.TodoDock
}

func GetSource(cfg config.TodoDock) TodoDockSource {
	return TodoDockSource{cfg}
}

func (t TodoDockSource) GetName() string {
	return "TodoDock"
}

func (t TodoDockSource) GetLabel() string {
	return t.params.Label
}

func (t TodoDockSource) GetNewCards() (cards []trello.Card, err error) {
	id, token, err := t.login()
	if err != nil {
		return cards, nil
	}
	tasks, err := t.fetchTasks(id, token)
	return toCards(tasks, t.params.Label)
}

// login logs-in to TodoDock with the configured user's credentials,
// and returns the user ID and JWT obtained from the HTTP response
func (t TodoDockSource) login() (id int, jwt string, err error) {
	req, err := json.Marshal(map[string]string{
		"email":    t.params.Email,
		"password": t.params.Password,
	})
	if err != nil {
		return -1, "", fmt.Errorf("could not build TodoDock login request body: %w", err)
	}

	url := fmt.Sprintf("%s/login", BASE_URL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return -1, "", fmt.Errorf("could not login to TodoDock: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	return int(data["data"]["id"].(float64)), fmt.Sprintf("%s", data["data"]["token"]), nil
}
