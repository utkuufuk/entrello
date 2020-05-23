package trello

import (
	"testing"

	"github.com/utkuufuk/entrello/internal/config"
)

func TestNewClient(t *testing.T) {
	str := "placeholder"

	tt := []struct {
		name     string
		boardId  string
		listId   string
		apiKey   string
		apiToken string
		err      bool
	}{
		{
			name:     "no errors",
			boardId:  str,
			listId:   str,
			apiKey:   str,
			apiToken: str,
			err:      false,
		},
		{
			name:     "missing board ID",
			boardId:  "",
			listId:   str,
			apiKey:   str,
			apiToken: str,
			err:      true,
		},
		{
			name:     "missing list ID",
			boardId:  str,
			listId:   "",
			apiKey:   str,
			apiToken: str,
			err:      true,
		},
		{
			name:     "missing api key",
			boardId:  str,
			listId:   str,
			apiKey:   "",
			apiToken: str,
			err:      true,
		},
		{
			name:     "missing api token",
			boardId:  str,
			listId:   str,
			apiKey:   str,
			apiToken: "",
			err:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Config{
				TrelloApiKey:   tc.apiKey,
				TrelloApiToken: tc.apiToken,
				BoardId:        tc.boardId,
				ListId:         tc.listId,
				Sources:        config.Sources{},
			}
			_, err := NewClient(cfg)
			if (err != nil && !tc.err) || err == nil && tc.err {
				t.Fatalf("did not expect the error outcome to be: '%t'", tc.err)
			}
		})
	}
}

func TestContains(t *testing.T) {
	str := "placeholder"

	tt := []struct {
		name    string
		list    []string
		query   string
		outcome bool
	}{
		{
			name:    "all empty",
			list:    []string{},
			query:   "",
			outcome: false,
		},
		{
			name:    "empty query",
			list:    []string{str},
			query:   "",
			outcome: false,
		},
		{
			name:    "empty list",
			list:    []string{},
			query:   str,
			outcome: false,
		},
		{
			name:    "no match",
			list:    []string{"a", "b", "c"},
			query:   "d",
			outcome: false,
		},
		{
			name:    "match",
			list:    []string{"a", "b", "c", str},
			query:   str,
			outcome: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			outcome := contains(tc.list, tc.query)
			if outcome == tc.outcome {
				return
			}

			prefix := "expected "
			if !tc.outcome {
				prefix = "did not expect "
			}
			t.Fatalf("%s %s to be in the list %v, got %t", prefix, tc.query, tc.list, outcome)
		})
	}
}
