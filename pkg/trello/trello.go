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
	existingCards map[string][]Card
}

func NewClient(cfg config.Trello) Client {
	return Client{
		api:           trello.NewClient(cfg.ApiKey, cfg.ApiToken),
		boardId:       cfg.BoardId,
		existingCards: make(map[string][]Card),
	}
}

// NewCard creates a new Trello card model with the given mandatory fields name,
// and the optional description and dueDate fields
func NewCard(name, description string, dueDate *time.Time) (card Card, err error) {
	if name == "" {
		return card, fmt.Errorf("card name cannot be blank")
	}

	return &trello.Card{
		Name: name,
		Desc: description,
		Due:  dueDate,
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
