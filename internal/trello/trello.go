package trello

import (
	"fmt"
	"time"

	"github.com/adlio/trello"
	"github.com/utkuufuk/entrello/internal/config"
)

// Client represents a Trello client model
type Client struct {
	// client is the Trello API client
	client *trello.Client

	// boardId is the ID of the board to read & write cards
	boardId string

	// list is the ID of the Trello list to insert new cards
	listId string

	// a map of existing cards in the board, where the key is the label ID and value is the card name
	existingCards map[string][]string
}

// Card represents a Trello card
type Card struct {
	name        string
	label       string
	description string
	dueDate     *time.Time
}

func NewClient(c config.Config) (Client, error) {
	if c.BoardId == "" || c.ListId == "" || c.TrelloApiKey == "" || c.TrelloApiToken == "" {
		return Client{}, fmt.Errorf("could not create trello client, missing configuration parameter(s)")
	}
	return Client{
		client:        trello.NewClient(c.TrelloApiKey, c.TrelloApiToken),
		boardId:       c.BoardId,
		listId:        c.ListId,
		existingCards: make(map[string][]string),
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
	return Card{name, label, description, dueDate}, nil
}

// LoadExistingCards retrieves and saves all existing cards from the board that has at least one
// of the given label IDs
func (c Client) LoadExistingCards(labels []string) error {
	board, err := c.client.GetBoard(c.boardId, trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not get board data: %w", err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not fetch cards in board: %w", err)
	}

	for _, label := range labels {
		c.existingCards[label] = make([]string, 0, len(cards))
	}

	for _, card := range cards {
		for _, label := range card.IDLabels {
			c.existingCards[label] = append(c.existingCards[label], card.Name)
		}
	}
	return nil
}

// @todo: optionally delete existing cards that do not appear in the new list
// UpdateCards creates the given cards except the ones that already exist
func (c Client) UpdateCards(cards []Card) error {
	for _, card := range cards {
		if contains(c.existingCards[card.label], card.name) {
			continue
		}

		if err := c.createCard(card); err != nil {
			return fmt.Errorf("[-] could not create card '%s': %v", card.name, err)
		}

		// @todo: send telegram notification if enabled
	}
	return nil
}

// createCard creates a Trello card using the the API client
func (c Client) createCard(card Card) error {
	return c.client.CreateCard(&trello.Card{
		Name:     card.name,
		Desc:     card.description,
		Due:      card.dueDate,
		IDList:   c.listId,
		IDLabels: []string{card.label},
	}, trello.Defaults())
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
