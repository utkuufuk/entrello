package habits

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
	"google.golang.org/api/sheets/v4"
)

const (
	nameRowIndex     = 0
	durationRowIndex = 1
	scoreRowIndex    = 2
	dataRowOffset    = 3 // number of rows before the first data row starts in the spreadsheet
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
	Score    float64
}

type cell struct {
	col string
	row int
}

func GetSource(cfg config.Habits) source {
	return source{cfg.SpreadsheetId, cfg.CredentialsFile, cfg.TokenFile, nil}
}

func (s source) FetchNewCards(
	ctx context.Context,
	cfg config.SourceConfig,
	now time.Time,
) ([]trello.Card, error) {
	if err := s.initializeService(ctx); err != nil {
		return nil, fmt.Errorf("could not initialize google spreadsheet service: %w", err)
	}

	habits, err := s.fetchHabits(now)
	if err != nil {
		return nil, fmt.Errorf("could not fetch habits: %w", err)
	}

	if err = s.updateScores(habits, now); err != nil {
		return nil, err
	}

	return toCards(habits, cfg.Label, now)
}

func (s source) updateScores(habits map[string]habit, now time.Time) error {
	scores := make([]float64, len(habits))
	var cellNameComponents []string
	for _, habit := range habits {
		cellNameComponents = strings.Split(habit.CellName, "!")
		col := []rune(cellNameComponents[1][0:1])[0]
		idx := int(col) - int('A') - 1
		scores[idx] = habit.Score
	}

	row := scoreRowIndex + 1
	firstCol := string(rune(int('A') + dataColumnOffset))
	lastCol := string(rune(int('A') + len(habits)))
	rangeName, err := getRangeName(now, cell{firstCol, row}, cell{lastCol, row})
	if err != nil {
		return fmt.Errorf("could not update habit scores: %w", err)
	}

	values := make([][]interface{}, 1)
	values[0] = make([]interface{}, len(scores))
	for i, score := range scores {
		values[0][i] = score
	}
	return s.writeCells(values, rangeName)
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

// writeCells writes a 2D array of values into a range of cells
func (s source) writeCells(values [][]interface{}, rangeName string) error {
	_, err := s.service.
		Update(s.spreadsheetId, rangeName, &sheets.ValueRange{Values: values}).
		ValueInputOption("USER_ENTERED").
		Do()

	if err != nil {
		return fmt.Errorf("could not write cells: %w", err)
	}
	return nil
}

// toCards returns a slice of trello cards from the given habits which haven't been marked today
func toCards(
	habits map[string]habit,
	label string,
	now time.Time,
) (cards []trello.Card, err error) {
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

// mapHabits creates a map of habits for given a date and a spreadsheet row data
func mapHabits(rows [][]interface{}, date time.Time) (map[string]habit, error) {
	habits := make(map[string]habit)
	for col := dataColumnOffset; col < len(rows[0]); col++ {
		// evaluate the habit's cell name for today
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

		// read habit name
		name := fmt.Sprintf("%v", rows[nameRowIndex][col])
		if name == "" {
			return nil, fmt.Errorf("habit name cannot be blank")
		}

		// read optional habit duration value
		duration := fmt.Sprintf("%v", rows[durationRowIndex][col])

		// calculate habit score
		nom := 0
		denom := date.Day()
		for row := dataRowOffset; row < date.Day()+dataRowOffset; row++ {
			if len(rows[row]) < col+1 {
				continue
			}

			val := rows[row][col]
			if val == "✔" {
				nom++
			}

			if val == "–" {
				denom--
			}
		}
		score := (float64(nom) / float64(denom))

		habits[name] = habit{cellName, state, duration, score}
	}
	return habits, nil
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
