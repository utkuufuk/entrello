package habits

import (
	"errors"
	"testing"
)

func TestToCards(t *testing.T) {
	str := "test"

	tt := []struct {
		name     string
		label    string
		habits   map[string]habit
		numCards int
		err      error
	}{
		{
			name:     "blank habit name",
			label:    str,
			habits:   map[string]habit{"": {}},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:     "missing label",
			label:    "",
			habits:   map[string]habit{str: {}},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:  "marked habits",
			label: str,
			habits: map[string]habit{
				"a": {str, "✔"},
				"b": {str, "x"},
				"c": {str, "✘"},
				"d": {str, "–"},
				"e": {str, "-"},
			},
			numCards: 0,
			err:      nil,
		},
		{
			name:  "some marked some unhabits",
			label: str,
			habits: map[string]habit{
				"a": {str, "✔"},
				"b": {str, "x"},
				"c": {str, "✘"},
				"d": {str, "–"},
				"e": {str, "-"},
				"f": {str, ""},
				"g": {str, ""},
			},
			numCards: 2,
			err:      nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := toCards(tc.habits, tc.label)
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if len(cards) != tc.numCards {
				t.Errorf("expected %d cards, got %d", tc.numCards, len(cards))
			}
		})
	}
}
