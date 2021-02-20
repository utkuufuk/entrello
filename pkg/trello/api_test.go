package trello

import (
	"testing"

	"github.com/adlio/trello"
)

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

func TestSetExistingCards(t *testing.T) {
	tt := []struct {
		name      string
		client    Client
		cards     []*trello.Card
		labels    []string
		numLabels int
		numCards  map[string]int
	}{
		{
			name:   "no labels",
			client: Client{existingCards: make(map[string][]Card)},
			cards: []*trello.Card{
				newTestCardByLabel([]string{"a"}),
				newTestCardByLabel([]string{"b"}),
			},
			labels:    []string{},
			numLabels: 0,
		},
		{
			name:   "no matching labels",
			client: Client{existingCards: make(map[string][]Card)},
			cards: []*trello.Card{
				newTestCardByLabel([]string{"a"}),
				newTestCardByLabel([]string{"b"}),
			},
			labels:    []string{"c"},
			numLabels: 1,
			numCards:  map[string]int{"c": 0},
		},
		{
			name:   "all matching labels",
			client: Client{existingCards: make(map[string][]Card)},
			cards: []*trello.Card{
				newTestCardByLabel([]string{"a"}),
				newTestCardByLabel([]string{"b"}),
			},
			labels:    []string{"b", "a"},
			numLabels: 2,
			numCards:  map[string]int{"a": 1, "b": 1},
		},
		{
			name:   "all matching overlapping labels",
			client: Client{existingCards: make(map[string][]Card)},
			cards: []*trello.Card{
				newTestCardByLabel([]string{"a", "b"}),
				newTestCardByLabel([]string{"a", "b"}),
			},
			labels:    []string{"b", "a"},
			numLabels: 2,
			numCards:  map[string]int{"a": 2, "b": 2},
		},
		{
			name:   "some matching labels",
			client: Client{existingCards: make(map[string][]Card)},
			cards: []*trello.Card{
				newTestCardByLabel([]string{"a"}),
				newTestCardByLabel([]string{"b"}),
				newTestCardByLabel([]string{"c"}),
			},
			labels:    []string{"b", "a"},
			numLabels: 2,
			numCards:  map[string]int{"a": 1, "b": 1},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.client.setExistingCards(tc.cards, tc.labels)

			if len(tc.client.existingCards) != tc.numLabels {
				t.Errorf("wanted %d keys in the map, got %d", tc.numLabels, len(tc.client.existingCards))
			}

			if tc.numLabels == 0 {
				return
			}

			for k, v := range tc.client.existingCards {
				if tc.numCards[k] != len(v) {
					t.Errorf("wanted %d cards for key %s, got %d", tc.numCards[k], k, len(v))
				}
			}
		})
	}
}

func newTestCardByLabel(labels []string) *trello.Card {
	return &trello.Card{
		IDLabels: labels,
	}
}
