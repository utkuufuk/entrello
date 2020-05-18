package trello

import (
	"fmt"
	"time"

	"github.com/adlio/trello"
	"github.com/utkuufuk/entrello/internal/config"
)

// Card represents a Trello card
type Card struct {
	name        string
	labelId     string
	description string
	dueDate     *time.Time
}

type Client struct {
	client  *trello.Client
	boardId string
	listId  string
}

func NewClient(cfg config.Config) Client {
	return Client{
		client:  trello.NewClient(cfg.TrelloApiKey, cfg.TrelloApiToken),
		boardId: cfg.BoardId,
		listId:  cfg.ListId,
	}
}

// FetchBoardCards retrieves all cards from the board and returns a map where the keys are card
// names and the values are lists of labels that each card has
func (c Client) FetchBoardCards() (map[string][]string, error) {
	board, err := c.client.GetBoard(c.boardId, trello.Defaults())
	if err != nil {
		return nil, fmt.Errorf("could not get board data: %w", err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		return nil, fmt.Errorf("could not fetch cards in board: %w", err)
	}

	m := make(map[string][]string)
	for _, card := range cards {
		m[card.Name] = card.IDLabels
	}
	return m, nil
}

// AddCard adds the specified card to the configured Trello list
func (c Client) AddCard(card Card) error {
	return c.client.CreateCard(&trello.Card{
		Name:     card.name,
		Desc:     card.description,
		Due:      card.dueDate,
		IDList:   c.listId,
		IDLabels: []string{card.labelId},
	}, trello.Defaults())
}

func CreateCard(name, labelId, description string, dueDate *time.Time) (card Card, err error) {
	if name == "" {
		return card, fmt.Errorf("card name cannot be blank")
	}

	if labelId == "" {
		return card, fmt.Errorf("label ID cannot be blank")
	}
	return Card{name, labelId, description, dueDate}, nil
}

func (c Card) GetName() string {
	return c.name
}

func (c Card) GetLabelId() string {
	return c.labelId
}
