package main

import (
	"github.com/utkuufuk/entrello/internal/trello"
)

type Source interface {
	GetCards() ([]trello.Card, error)
	GetName() string
}
