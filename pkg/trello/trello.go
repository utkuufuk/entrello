package trello

import (
	"fmt"
	"time"

	"github.com/adlio/trello"
	"github.com/utkuufuk/entrello/internal/config"
)

type Card *trello.Card

type Client struct {
	api           *trello.Client
	boardId       string
	todoListId    string
	todayListId   string
	existingCards map[string][]Card
}

func NewClient(cfg config.Trello) (client Client, err error) {
	if cfg.BoardId == "" || cfg.TodoListId == "" || cfg.TodayListId == "" || cfg.ApiKey == "" || cfg.ApiToken == "" {
		return client, fmt.Errorf("could not create trello client, missing configuration parameter(s)")
	}

	return Client{
		api:           trello.NewClient(cfg.ApiKey, cfg.ApiToken),
		boardId:       cfg.BoardId,
		todoListId:    cfg.TodoListId,
		todayListId:   cfg.TodayListId,
		existingCards: make(map[string][]Card),
	}, nil
}

// NewCard creates a new Trello card model with the given mandatory fields name, label, description,
// and the optional dueDate field
func NewCard(name, label, description string, dueDate *time.Time) (card Card, err error) {
	if name == "" {
		return card, fmt.Errorf("card name cannot be blank")
	}

	if label == "" {
		return card, fmt.Errorf("label ID cannot be blank")
	}

	return &trello.Card{
		Name:     name,
		Desc:     description,
		Due:      dueDate,
		IDLabels: []string{label},
	}, nil
}

// FilterNewAndStale compares the given cards with the existing cards and returns two arrays;
// one containing new cards and the other containing stale cards.
func (c Client) FilterNewAndStale(cards []Card, label string) (new, stale []Card) {
	m := make(map[string]*trello.Card)
	for _, card := range c.existingCards[label] {
		m[card.Name] = card
	}

	for _, card := range cards {
		_, ok := m[card.Name]
		m[card.Name] = nil
		if ok {
			continue
		}
		new = append(new, card)
	}

	for _, card := range m {
		if card == nil {
			continue
		}
		stale = append(stale, card)
	}

	return new, stale
}
