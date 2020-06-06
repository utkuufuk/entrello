package habits

import (
	"context"
	"fmt"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

type source struct {
	spreadsheetId   string
	credentialsFile string
	tokenFile       string
	service         service
}

func GetSource(cfg config.Habits) source {
	return source{cfg.SpreadsheetId, cfg.CredentialsFile, cfg.TokenFile, nil}
}

func (s source) FetchNewCards(ctx context.Context, cfg config.SourceConfig) ([]trello.Card, error) {
	err := s.createService(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not create spreadsheet service: %w", err)
	}
	data, _ := s.readCells(s.spreadsheetId, "Jun 2020!B1:M1")
	fmt.Println(data)
	return []trello.Card{}, nil
}
