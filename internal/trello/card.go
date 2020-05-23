package trello

import (
	"fmt"
	"time"
)

// Card represents a Trello card
type Card struct {
	name        string
	label       string
	description string
	dueDate     *time.Time
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
