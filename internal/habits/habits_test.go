package habits

import (
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
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

func TestMapHabits(t *testing.T) {
	any := "."

	tt := []struct {
		name string
		rows [][]string
		out  map[string]habit
		err  error
	}{
		{
			name: "all marked",
			rows: [][]string{
				{"", "a", "b", "c"},
				{any, any, any, any},
				{any, any, any, any},
				{any, "✔", "✘", "–"},
			},
			out: map[string]habit{
				"a": {"Jan 2020!B4", "✔"},
				"b": {"Jan 2020!C4", "✘"},
				"c": {"Jan 2020!D4", "–"},
			},
		},
		{
			name: "blank mid rows",
			rows: [][]string{
				{"", "a", "b", "c"},
				{},
				{},
				{any, "✔", "✘", "–"},
			},
			out: map[string]habit{
				"a": {"Jan 2020!B4", "✔"},
				"b": {"Jan 2020!C4", "✘"},
				"c": {"Jan 2020!D4", "–"},
			},
		},
		{
			name: "blank cell in the middle",
			rows: [][]string{
				{"", "a", "b", "c", "d"},
				{any, any, any, any, any},
				{any, any, any, any, any},
				{any, "✔", "✘", "", "–"},
			},
			out: map[string]habit{
				"a": {"Jan 2020!B4", "✔"},
				"b": {"Jan 2020!C4", "✘"},
				"c": {"Jan 2020!D4", ""},
				"d": {"Jan 2020!E4", "–"},
			},
		},
		{
			name: "blank cells in the end",
			rows: [][]string{
				{any, "a", "b", "c", "d", "e"},
				{any, any, any, any, any, any},
				{any, any, any, any, any, any},
				{any, "✔", "✘", "–"},
			},
			out: map[string]habit{
				"a": {"Jan 2020!B4", "✔"},
				"b": {"Jan 2020!C4", "✘"},
				"c": {"Jan 2020!D4", "–"},
				"d": {"Jan 2020!E4", ""},
				"e": {"Jan 2020!F4", ""},
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			date := time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)

			data := make([][]interface{}, 0, len(tc.rows))
			for r, row := range tc.rows {
				data = append(data, make([]interface{}, 0, len(row)))
				for _, col := range row {
					data[r] = append(data[r], col)
				}
			}

			habits, err := mapHabits(data, date)
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if diff := cmp.Diff(habits, tc.out); diff != "" {
				t.Errorf("output diff: %s", diff)
			}
		})
	}
}

func TestGetRangeName(t *testing.T) {
	tt := []struct {
		name  string
		year  int
		month int
		start cell
		end   cell
		out   string
		err   error
	}{
		{
			name:  "invalid start col",
			year:  2020,
			month: 1,
			start: cell{"", 1},
			end:   cell{},
			err:   errors.New(""),
		},
		{
			name:  "invalid start row",
			year:  2020,
			month: 1,
			start: cell{"A", 0},
			end:   cell{},
			err:   errors.New(""),
		},
		{
			name:  "invalid end col",
			year:  2020,
			month: 1,
			start: cell{"A", 1},
			end:   cell{"0", 1},
			err:   errors.New(""),
		},
		{
			name:  "invalid end row",
			year:  2020,
			month: 1,
			start: cell{"A", 1},
			end:   cell{"A", -1},
			err:   errors.New(""),
		},
		{
			name:  "implicit single cell",
			year:  2020,
			month: 1,
			start: cell{"A", 1},
			end:   cell{},
			out:   "Jan 2020!A1",
			err:   nil,
		},
		{
			name:  "explicit single cell",
			year:  2020,
			month: 1,
			start: cell{"A", 1},
			end:   cell{"A", 1},
			out:   "Jan 2020!A1",
			err:   nil,
		},
		{
			name:  "valid range",
			year:  2020,
			month: 1,
			start: cell{"B", 3},
			end:   cell{"D", 5},
			out:   "Jan 2020!B3:D5",
			err:   nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			date := time.Date(tc.year, time.Month(tc.month), 1, 0, 0, 0, 0, time.UTC)
			out, err := getRangeName(date, tc.start, tc.end)
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if err == nil {
				return
			}

			if out != tc.out {
				t.Fatalf("range name mismatch; want '%s', got '%s'", tc.out, out)
			}
		})
	}
}
