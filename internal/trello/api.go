package trello

import (
	"fmt"

	"github.com/adlio/trello"
)

// ArchiveCard archives a Trello card using the the Trello API
func (c Client) ArchiveCard(card Card) error {
	return (*trello.Card)(card).Update(trello.Arguments{"closed": "true"})
}

// CreateCard creates a Trello card using the the Trello API
func (c Client) CreateCard(card Card) error {
	card.IDList = c.listId
	return c.api.CreateCard(card, trello.Defaults())
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
			c.existingCards[label] = append(c.existingCards[label], card)
		}
	}
}
