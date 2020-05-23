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

func CreateCard(name, label, description string, dueDate *time.Time) (card Card, err error) {
	if name == "" {
		return card, fmt.Errorf("card name cannot be blank")
	}

	if label == "" {
		return card, fmt.Errorf("label ID cannot be blank")
	}
	return Card{name, label, description, dueDate}, nil
}
