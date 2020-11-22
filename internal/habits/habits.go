package habits

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
	"google.golang.org/api/sheets/v4"
)

const (
	nameRowIndex     = 0
	durationRowIndex = 1
	dataRowOffset    = 4 // number of rows before the first data row starts in the spreadsheet
	dataColumnOffset = 1 // number of columns before the first data column starts in the spreadsheet
	dueHour          = 23
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
	Duration string
}

type cell struct {
	col string
	row int
}

func GetSource(cfg config.Habits) source {
	return source{cfg.SpreadsheetId, cfg.CredentialsFile, cfg.TokenFile, nil}
}

func (s source) FetchNewCards(ctx context.Context, cfg config.SourceConfig, now time.Time) ([]trello.Card, error) {
	err := s.initializeService(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not initialize google spreadsheet service: %w", err)
	}

	habits, err := s.fetchHabits(now)
	if err != nil {
		return nil, fmt.Errorf("could not fetch habits: %w", err)
	}

	return toCards(habits, cfg.Label, now)
}

// fetchHabits retrieves the state of today's habits from the spreadsheet
func (s source) fetchHabits(now time.Time) (map[string]habit, error) {
	rangeName, err := getRangeName(now, cell{"A", 1}, cell{"Z", now.Day() + dataRowOffset})
	if err != nil {
		return nil, fmt.Errorf("could not get range name: %w", err)
	}

	rows, err := s.readCells(rangeName)
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}

	return mapHabits(rows, now)
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
func toCards(habits map[string]habit, label string, now time.Time) (cards []trello.Card, err error) {
	for name, habit := range habits {
		if habit.State != "" {
			continue
		}

		// include the day of month in card title to force overwrite in the beginning of the next day
		title := fmt.Sprintf("%v (%d)", name, now.Day())

		// include optional duration info in card description
		description := habit.CellName
		if habit.Duration != "–" {
			description = fmt.Sprintf("%s\nDuration:%s", description, habit.Duration)
		}

		due := time.Date(now.Year(), now.Month(), now.Day(), dueHour, 0, 0, 0, now.Location())
		c, err := trello.NewCard(title, label, description, &due)
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
	for col := dataColumnOffset; col < len(rows[0]); col++ {
		c := cell{string(rune('A' + col)), date.Day() + dataRowOffset}
		cellName, err := getRangeName(date, c, c)
		if err != nil {
			return nil, err
		}

		// handle cases where the last N columns are blank which reduces the slice length by N
		state := ""
		if col < len(rows[date.Day()+dataRowOffset-1]) {
			state = fmt.Sprintf("%v", rows[date.Day()+dataRowOffset-1][col])
		}

		name := fmt.Sprintf("%v", rows[nameRowIndex][col])
		if name == "" {
			return nil, fmt.Errorf("habit name cannot be blank")
		}

		duration := fmt.Sprintf("%v", rows[durationRowIndex][col])
		states[name] = habit{cellName, state, duration}
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
