package trello

import "time"

type Card struct {
	Name    string
	Label   string
	List    string
	DueDate time.Time
}
