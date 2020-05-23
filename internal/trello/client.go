package trello

import (
	"fmt"

	"github.com/adlio/trello"
	"github.com/utkuufuk/entrello/internal/config"
)

type Client struct {
	client            *trello.Client
	boardId           string
	listId            string
	existingCardNames map[string][]string
}

func NewClient(c config.Config) (Client, error) {
	if c.BoardId == "" || c.ListId == "" || c.TrelloApiKey == "" || c.TrelloApiToken == "" {
		return Client{}, fmt.Errorf("could not create trello client, missing configuration parameter(s)")
	}
	return Client{
		client:            trello.NewClient(c.TrelloApiKey, c.TrelloApiToken),
		boardId:           c.BoardId,
		listId:            c.ListId,
		existingCardNames: make(map[string][]string),
	}, nil
}

// LoadExistingCardNames retrieves and saves all existing cards from the board that has at least one
// of the given label IDs
func (c Client) LoadExistingCardNames(labels []string) error {
	board, err := c.client.GetBoard(c.boardId, trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not get board data: %w", err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not fetch cards in board: %w", err)
	}

	for _, label := range labels {
		c.existingCardNames[label] = make([]string, 0, len(cards))
	}

	for _, card := range cards {
		for _, label := range card.IDLabels {
			c.existingCardNames[label] = append(c.existingCardNames[label], card.Name)
		}
	}
	return nil
}

// @todo: optionally delete existing cards that do not appear in the new list
// UpdateCards creates the given cards except the ones that already exist
func (c Client) UpdateCards(cards []Card) error {
	for _, card := range cards {
		if contains(c.existingCardNames[card.label], card.name) {
			continue
		}

		if err := c.createCard(card); err != nil {
			return fmt.Errorf("[-] could not create card '%s': %v", card.name, err)
		}

		// @todo: send telegram notification if enabled
	}
	return nil
}

// createCard creates a Trello card
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
