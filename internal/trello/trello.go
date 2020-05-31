package trello

import (
	"fmt"
	"time"

	"github.com/adlio/trello"
	"github.com/utkuufuk/entrello/internal/config"
)

type Card *trello.Card

// Client represents a Trello client model
type Client struct {
	// client is the Trello API client
	api *trello.Client

	// boardId is the ID of the board to read & write cards
	boardId string

	// list is the ID of the Trello list to insert new cards
	listId string

	// a map of existing cards in the board, where the key is the label ID and value is the card name
	existingCards map[string][]Card
}

func NewClient(cfg config.Trello) (client Client, err error) {
	if cfg.BoardId == "" || cfg.ListId == "" || cfg.ApiKey == "" || cfg.ApiToken == "" {
		return client, fmt.Errorf("could not create trello client, missing configuration parameter(s)")
	}

	return Client{
		api:           trello.NewClient(cfg.ApiKey, cfg.ApiToken),
		boardId:       cfg.BoardId,
		listId:        cfg.ListId,
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

	if description == "" {
		return card, fmt.Errorf("description cannot be blank")
	}

	return &trello.Card{
		Name:     name,
		Desc:     description,
		Due:      dueDate,
		IDLabels: []string{label},
	}, nil
}

// LoadCards retrieves existing cards from the board that have at least one of the given label IDs
func (c Client) LoadCards(labels []string) error {
	board, err := c.api.GetBoard(c.boardId, trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not get board data: %w", err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not fetch cards in board: %w", err)
	}

	for _, label := range labels {
		c.existingCards[label] = make([]Card, 0, len(cards))
	}

	for _, card := range cards {
		for _, label := range card.IDLabels {
			c.existingCards[label] = append(c.existingCards[label], card)
		}
	}
	return nil
}

// CompareWithExisting compares the given cards with the existing cards and returns two arrays;
// one containing new cards and the other containing stale cards.
func (c Client) CompareWithExisting(cards []Card, label string) (new, stale []Card) {
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

// CreateCard creates a Trello card using the the Trello API
func (c Client) CreateCard(card Card) error {
	card.IDList = c.listId
	return c.api.CreateCard(card, trello.Defaults())
}

// ArchiveCard archives a Trello card using the the Trello API
func (c Client) ArchiveCard(card Card) error {
	return (*trello.Card)(card).Update(trello.Arguments{"closed": "true"})
}

// contains returns true if the list of strings contain the given string
func contains(list []string, item string) bool {
	if item == "" {
		return false
	}

	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
