package habits

import (
	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

type HabitsSource struct {
	cfg config.Habits
}

func GetSource(cfg config.Habits) HabitsSource {
	return HabitsSource{cfg}
}

func (h HabitsSource) FetchNewCards() ([]trello.Card, error) {
	// @todo: implement
	return []trello.Card{}, nil
}
