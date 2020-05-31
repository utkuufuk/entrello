package trello

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/utkuufuk/entrello/internal/config"
)

func TestNewCard(t *testing.T) {
	tt := []struct {
		name   string
		cName  string
		cLabel string
		cDesc  string
		err    error
	}{
		{
			name:   "no errors",
			cName:  "name",
			cLabel: "label",
			cDesc:  "desc",
			err:    nil,
		},
		{
			name:   "missing name",
			cName:  "",
			cLabel: "label",
			cDesc:  "desc",
			err:    fmt.Errorf("card name cannot be blank"),
		},
		{
			name:   "missing description",
			cName:  "name",
			cLabel: "label",
			cDesc:  "",
			err:    fmt.Errorf("description cannot be blank"),
		},
		{
			name:   "missing label ID",
			cName:  "name",
			cLabel: "",
			cDesc:  "desc",
			err:    fmt.Errorf("label ID cannot be blank"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var opts cmp.Options
			opts = append(opts, cmp.Comparer(func(x, y error) bool {
				return (x == nil && y == nil) || (x.Error() == y.Error())
			}))

			_, err := NewCard(tc.cName, tc.cLabel, tc.cDesc, nil)
			if diff := cmp.Diff(err, tc.err, opts...); diff != "" {
				t.Errorf("errors diff: %s", diff)
			}
		})
	}
}

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
				Trello: config.Trello{
					ApiKey:   tc.apiKey,
					ApiToken: tc.apiToken,
					BoardId:  tc.boardId,
					ListId:   tc.listId,
				},
				Sources: config.Sources{},
			}
			_, err := NewClient(cfg.Trello)
			if (err != nil && !tc.err) || err == nil && tc.err {
				t.Fatalf("did not expect the error outcome to be: '%t'", tc.err)
			}
		})
	}
}
