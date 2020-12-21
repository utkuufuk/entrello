package trello

import (
	"fmt"
	"time"

	"github.com/adlio/trello"
)

// DeleteCard deletes a Trello card using the the Trello API
func (c Client) DeleteCard(card Card) error {
	path := fmt.Sprintf("cards/%s", card.ID)
	return c.api.Delete(path, trello.Defaults(), card)
}

// CreateCard creates a Trello card using the the Trello API
func (c Client) CreateCard(card Card, now time.Time) error {
	if card.Due == nil {
		card.IDList = c.todoListId
		return c.api.CreateCard(card, trello.Defaults())
	}

	dueYear := card.Due.Year()
	dueMonth := card.Due.Month()
	dueDay := card.Due.Day()

	thisYear := now.Year()
	thisMonth := now.Month()
	today := now.Day()

	// insert into the todo list, unless the due date is sometime today
	if dueYear > thisYear {
		card.IDList = c.todoListId
	} else if dueYear == thisYear && dueMonth > thisMonth {
		card.IDList = c.todoListId
	} else if dueYear == thisYear && dueMonth == thisMonth && dueDay > today {
		card.IDList = c.todoListId
	} else {
		card.IDList = c.todayListId
	}

	return c.api.CreateCard(card, trello.Defaults())
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
			if ok := contains(labels, label); !ok {
				continue
			}
			c.existingCards[label] = append(c.existingCards[label], card)
		}
	}
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
