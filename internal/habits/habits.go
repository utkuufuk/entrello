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

func (h HabitsSource) IsEnabled() bool {
	return h.cfg.Enabled
}

func (h HabitsSource) IsStrict() bool {
	return h.cfg.Strict
}

func (h HabitsSource) GetName() string {
	return "Google Spreadsheet Habits"
}

func (h HabitsSource) GetLabel() string {
	return h.cfg.Label
}

func (h HabitsSource) GetPeriod() config.Period {
	return h.cfg.Period
}

func (h HabitsSource) FetchNewCards() ([]trello.Card, error) {
	// @todo: implement
	return []trello.Card{}, nil
}
