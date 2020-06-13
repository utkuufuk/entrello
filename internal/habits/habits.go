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
	CellName string
	State    string
}

type cell struct {
	col string
	row int
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
	rangeName, err := getRangeName(today, cell{"B", 1}, cell{"Z", today.Day() + 3})
	if err != nil {
		return nil, fmt.Errorf("could not get range name: %w", err)
	}

	rows, err := s.readCells(rangeName)
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}

	return mapHabits(rows, today)
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
		if habit.State != "" {
			continue
		}

		c, err := trello.NewCard(fmt.Sprintf("%v", name), label, habit.CellName, nil)
		if err != nil {
			return nil, fmt.Errorf("could not create habit card: %w", err)
		}

		cards = append(cards, c)
	}
	return cards, nil
}

// mapHabits creates a state map for given a date and a spreadsheet row data
func mapHabits(rows [][]interface{}, date time.Time) (map[string]habit, error) {
	states := make(map[string]habit)
	for i := 0; i < len(rows[0]); i++ {
		c := cell{string('A' + 1 + i), date.Day() + 3}
		cellName, err := getRangeName(date, c, c)
		if err != nil {
			return states, err
		}

		// handle cases where the last N columns are blank which reduces the slice length by N
		state := ""
		if i < len(rows[date.Day()+2]) {
			state = fmt.Sprintf("%v", rows[date.Day()+2][i])
		}

		name := fmt.Sprintf("%v", rows[0][i])
		states[name] = habit{cellName, state}
	}
	return states, nil
}

// getRangeName gets the range name given a date and start & end cells
func getRangeName(date time.Time, start, end cell) (string, error) {
	if start.col < "A" || start.col > "Z" || start.row <= 0 {
		return "", fmt.Errorf("invalid start cell: %s%d", start.col, start.row)
	}

	month := date.Month().String()[:3]
	year := date.Year()

	// assume single cell if no end date specified
	if end.col == "" || end.row == 0 || (end.col == start.col && end.row == start.row) {
		return fmt.Sprintf("%s %d!%s%d", month, year, start.col, start.row), nil
	}

	if end.col < "A" || end.col > "Z" || end.row <= 0 {
		return "", fmt.Errorf("invalid end cell: %s%d", end.col, end.row)
	}

	return fmt.Sprintf("%s %d!%s%d:%s%d", month, year, start.col, start.row, end.col, end.row), nil
}
