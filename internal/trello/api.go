package trello

import (
	"fmt"

	"github.com/adlio/trello"
)

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

// CreateCard creates a Trello card using the the Trello API
func (c Client) CreateCard(card Card) error {
	card.IDList = c.listId
	return c.api.CreateCard(card, trello.Defaults())
}

// ArchiveCard archives a Trello card using the the Trello API
func (c Client) ArchiveCard(card Card) error {
	return (*trello.Card)(card).Update(trello.Arguments{"closed": "true"})
}
