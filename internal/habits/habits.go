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

type habit struct {
	cellName string
	state    string
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
func (s source) fetchHabits() (map[string]habit, error) {
	today := time.Now()
	month := today.Month().String()[:3]
	year := today.Year()
	rowNo := today.Day() + 3

	rangeName := fmt.Sprintf("%s %d!B1:Z%d", month, year, rowNo)
	rows, err := s.readCells(rangeName)
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}

	states := make(map[string]habit)
	for i := 0; i < len(rows[0]); i++ {
		name := fmt.Sprintf("%v", rows[0][i])
		state := fmt.Sprintf("%v", rows[today.Day()+2][i])
		cellName := fmt.Sprintf("%s %d!%s%d", month, year, string('A'+1+i), rowNo)
		states[name] = habit{cellName, state}
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
func toCards(habits map[string]habit, label string) (cards []trello.Card, err error) {
	for name, habit := range habits {
		if habit.state != "" {
			continue
		}

		c, err := trello.NewCard(fmt.Sprintf("%v", name), label, habit.cellName, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create habit card: %w", err)
		}

		cards = append(cards, c)
	}
	return cards, nil
}
