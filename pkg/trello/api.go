package trello

import (
	"fmt"

	"github.com/adlio/trello"
	"golang.org/x/exp/slices"
)

// DeleteCard deletes a Trello card
func (c Client) DeleteCard(card Card) error {
	path := fmt.Sprintf("cards/%s", card.ID)
	return c.api.Delete(path, trello.Defaults(), card)
}

// CreateCard creates a Trello card
func (c Client) CreateCard(card Card, label string, listId string) error {
	card.IDLabels = []string{label}
	card.IDList = listId
	return c.api.CreateCard(card, trello.Defaults())
}

// GetCard fetches a Trello card by its ID
func (c Client) GetCard(id string) (Card, error) {
	return c.api.GetCard(id, trello.Defaults())
}

// LoadBoard retrieves existing cards from the board that have at least one of the given label IDs
func (c Client) LoadBoard(labels []string) error {
	board, err := c.api.GetBoard(c.boardId, trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not get board data: %w", err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return fmt.Errorf("could not fetch cards in board: %w", err)
	}

	c.setExistingCards(cards, labels)
	return nil
}

// setExistingCards populates the map within the client from the given cards where the keys
// are labels and the values are card slices
func (c Client) setExistingCards(cards []*trello.Card, labels []string) {
	for _, label := range labels {
		c.existingCards[label] = make([]Card, 0, len(cards))
	}

	for _, card := range cards {
		for _, label := range card.IDLabels {
			if ok := slices.Contains(labels, label); !ok {
				continue
			}
			c.existingCards[label] = append(c.existingCards[label], card)
		}
	}
}
