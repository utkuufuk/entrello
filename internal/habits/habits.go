package habits

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
	"google.golang.org/api/sheets/v4"
)

type source struct {
	spreadsheetId   string
	credentialsFile string
	tokenFile       string
	service         *sheets.SpreadsheetsValuesService
}

func GetSource(cfg config.Habits) source {
	return source{cfg.SpreadsheetId, cfg.CredentialsFile, cfg.TokenFile, nil}
}

func (s source) FetchNewCards(ctx context.Context, cfg config.SourceConfig) ([]trello.Card, error) {
	err := s.initializeService(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not initialize google spreadsheet service: %w", err)
	}
	habits, err := s.fetchHabits()
	if err != nil {
		return nil, fmt.Errorf("could not fetch habits: %w", err)
	}

	return toCards(habits, cfg.Label)
}

// fetchHabits retrieves the state of today's habits from the spreadsheet
func (s source) fetchHabits() (map[string]string, error) {
	today := time.Now()
	rangeName := fmt.Sprintf("%s %d!B1:Z%d", today.Month().String()[:3], today.Year(), today.Day()+3)
	rows, err := s.readCells(rangeName)
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}

	states := make(map[string]string)
	for i := 0; i < len(rows[0]); i++ {
		name := fmt.Sprintf("%v", rows[0][i])
		state := fmt.Sprintf("%v", rows[today.Day()+2][i])
		states[name] = state
	}

	return states, nil
}

// readCells reads a range of cell values with the given range
func (s source) readCells(rangeName string) ([][]interface{}, error) {
	resp, err := s.service.Get(s.spreadsheetId, rangeName).Do()
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}
	return resp.Values, nil
}

// toCards returns a slice of trello cards from the given habits which haven't been marked today
func toCards(habits map[string]string, label string) (cards []trello.Card, err error) {
	for habit, state := range habits {
		if state != "" {
			continue
		}

		c, err := trello.NewCard(fmt.Sprintf("%v", habit), label, "", nil)
		if err != nil {
			return nil, fmt.Errorf("could not create habit card: %w", err)
		}

		cards = append(cards, c)
	}
	return cards, nil
}
