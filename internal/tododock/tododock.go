package tododock

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/pkg/trello"
)

const (
	BASE_URL = "https://tododock.com/api"
)

type source struct {
	email    string
	password string
}

// task represents the TodoDock task model
type task struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	State         string `json:"state"`
	Color         string `json:"color"`
	Notes         string `json:"notes"`
	Period        int    `json:"period"`
	NextResetDate string `json:"next_reset_date"`
	MuteEmails    int    `json:"mute_reminder_emails"`
}

func GetSource(cfg config.TodoDock) source {
	return source{cfg.Email, cfg.Password}
}

func (s source) FetchNewCards(ctx context.Context, cfg config.SourceConfig, now time.Time) (cards []trello.Card, err error) {
	id, token, err := s.login()
	if err != nil {
		return cards, fmt.Errorf("failed to authenticate with TodoDock: %w", err)
	}
	tasks, err := s.fetchTasks(id, token)
	return toCards(tasks, cfg.Label, now)
}

// toCards cherry-picks the 'active' and 'due' tasks from a list of tasks,
// then returns a list of cards containing those
func toCards(tasks []task, label string, now time.Time) (cards []trello.Card, err error) {
	cards = make([]trello.Card, 0, len(tasks))
	soon := now.AddDate(0, 0, 1)
	for _, t := range tasks {
		d, ok, err := shouldCreateCard(t, soon)
		if !ok {
			if err != nil {
				return cards, err
			}
			continue
		}

		url := fmt.Sprintf("https://tododock.com/home/%d\n%s", t.Id, t.Notes)
		c, err := trello.NewCard(t.Name, label, url, &d)
		if err != nil {
			return cards, fmt.Errorf("could not create card: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}

// shouldCreateCard decides if a Trello card should be created from the given TodoDock task
// by looking at the 'status' & 'next reset date' attributes of the task
func shouldCreateCard(t task, ref time.Time) (d time.Time, ok bool, err error) {
	d, err = time.Parse("2006-01-02 15:04:05", t.NextResetDate)
	if err != nil {
		return ref, false, fmt.Errorf("could not parse next reset date '%s': %w", t.NextResetDate, err)
	}

	// only create a card for active taks that are due soon
	if t.State != "active" || d.After(ref) {
		return ref, false, nil
	}
	return d, true, nil
}
