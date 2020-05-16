package trello

import (
	"fmt"
	"time"

	"github.com/adlio/trello"
	"github.com/utkuufuk/entrello/internal/config"
)

// Card represents a Trello card
type Card struct {
	// TODO: probably should include ID as well if we want to auto-reset on archive via webhooks
	Name        string
	Description string
	DueDate     time.Time
}

type Client struct {
	client  *trello.Client
	boardId string
	listId  string
	labelId string
}

func NewClient(cfg config.Config) Client {
	return Client{
		client:  trello.NewClient(cfg.TrelloApiKey, cfg.TrelloApiToken),
		boardId: cfg.BoardId,
		listId:  cfg.ListId,
		labelId: cfg.LabelId,
	}
}

// FetchBoardCards retrieves all cards within the configured board and returns the cards
// with the "TodoDock" label as a string->bool map (so that it's trivial to look up card names)
func (c Client) FetchBoardCards() (map[string]bool, error) {
	board, err := c.client.GetBoard(c.boardId, trello.Defaults())
	if err != nil {
		return nil, fmt.Errorf("could not get board data: %w", err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return nil, fmt.Errorf("could not fetch cards in board: %w", err)
	}

	// add card in the map only if it contains the TodoDock label
	// FIXME: must return all cards & save existing labels, otherwise this won't allow multiple sources
	m := map[string]bool{}
	for _, card := range cards {
		for _, label := range card.IDLabels {
			if label == c.labelId {
				m[card.Name] = true
				break
			}
		}
	}
	return m, nil
}

// AddCard adds the specified card to the configured Trello list
func (c Client) AddCard(card Card) error {
	return c.client.CreateCard(&trello.Card{
		Name:     card.Name,
		Desc:     card.Description,
		Due:      &card.DueDate,
		IDList:   c.listId,
		IDLabels: []string{c.labelId},
	}, trello.Defaults())
}
