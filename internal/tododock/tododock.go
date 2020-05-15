package tododock

import (
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
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

func (t TodoDockSource) GetCards() ([]trello.Card, error) {
	return []trello.Card{
		{
			Name:    "name",
			Label:   "home",
			List:    "To-Do",
			DueDate: time.Now(),
		},
	}, nil
}
