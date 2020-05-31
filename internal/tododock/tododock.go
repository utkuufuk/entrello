package tododock

import (
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

const (
	BASE_URL = "https://tododock.com/api"
)

type TodoDockSource struct {
	cfg config.TodoDock
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

func GetSource(cfg config.TodoDock) TodoDockSource {
	return TodoDockSource{cfg}
}

func (t TodoDockSource) IsEnabled() bool {
	return t.cfg.Enabled
}

func (t TodoDockSource) IsStrict() bool {
	return t.cfg.Strict
}

func (t TodoDockSource) GetName() string {
	return "TodoDock"
}

func (t TodoDockSource) GetLabel() string {
	return t.cfg.Label
}

func (t TodoDockSource) GetPeriod() config.Period {
	return t.cfg.Period
}

func (t TodoDockSource) FetchNewCards() (cards []trello.Card, err error) {
	id, token, err := t.login()
	if err != nil {
		return cards, nil
	}
	tasks, err := t.fetchTasks(id, token)
	return toCards(tasks, t.cfg.Label)
}

// toCards cherry-picks the 'active' and 'due' tasks from a list of tasks,
// then returns a list of cards containing those
func toCards(tasks []task, label string) (cards []trello.Card, err error) {
	cards = make([]trello.Card, 0, len(tasks))
	soon := time.Now().AddDate(0, 0, 2)
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
