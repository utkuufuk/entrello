package trello

import (
	"fmt"
	"testing"

	"github.com/adlio/trello"
	"github.com/google/go-cmp/cmp"
)

func TestNewCard(t *testing.T) {
	tt := []struct {
		name  string
		cName string
		cDesc string
		err   error
	}{
		{
			name:  "no errors",
			cName: "name",
			cDesc: "desc",
			err:   nil,
		},
		{
			name:  "missing name",
			cName: "",
			cDesc: "desc",
			err:   fmt.Errorf("card name cannot be blank"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var opts cmp.Options
			opts = append(opts, cmp.Comparer(func(x, y error) bool {
				return (x == nil && y == nil) || (x.Error() == y.Error())
			}))

			_, err := NewCard(tc.cName, tc.cDesc, nil)
			if diff := cmp.Diff(err, tc.err, opts...); diff != "" {
				t.Errorf("errors diff: %s", diff)
			}
		})
	}
}

func TestFilterNewAndStale(t *testing.T) {
	label := "label"
	tt := []struct {
		name     string
		client   Client
		cards    []Card
		label    string
		numNew   int
		numStale int
	}{
		{
			name: "only existing cards",
			client: Client{existingCards: map[string][]Card{label: {
				newTestCardByName("a"),
				newTestCardByName("b"),
			}}},
			cards:    []Card{newTestCardByName("a"), newTestCardByName("b")},
			numNew:   0,
			numStale: 0,
		},
		{
			name:     "only new cards",
			client:   Client{existingCards: map[string][]Card{label: {}}},
			cards:    []Card{newTestCardByName("a"), newTestCardByName("b")},
			numNew:   2,
			numStale: 0,
		},
		{
			name: "only stale cards",
			client: Client{existingCards: map[string][]Card{label: {
				newTestCardByName("a"),
				newTestCardByName("b"),
			}}},
			cards:    []Card{},
			numNew:   0,
			numStale: 2,
		},
		{
			name: "new, stale and existing cards",
			client: Client{existingCards: map[string][]Card{label: {
				newTestCardByName("a"),
				newTestCardByName("b"), // stale
				newTestCardByName("c"), // stale
			}}},
			cards: []Card{
				newTestCardByName("a"), // existing
				newTestCardByName("d"), // new
			},
			numNew:   1,
			numStale: 2,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			new, stale := tc.client.FilterNewAndStale(tc.cards, label)

			if len(new) != tc.numNew {
				t.Errorf("wanted %d new cards, got %d", tc.numNew, len(new))
			}

			if len(stale) != tc.numStale {
				t.Errorf("wanted %d stale cards, got %d", tc.numStale, len(stale))
			}
		})
	}
}

func newTestCardByName(name string) *trello.Card {
	return &trello.Card{
		Name: name,
	}
}
