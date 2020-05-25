package tododock

import (
	"testing"
)

func TestToCards(t *testing.T) {
	str := "placeholder"
	past := "2000-01-01 15:00:00"
	future := "2099-01-01 05:30:30"

	tt := []struct {
		name     string
		tasks    []task
		label    string
		numCards int
		err      bool
	}{
		{
			name:     "no tasks",
			tasks:    []task{},
			label:    str,
			numCards: 0,
			err:      false,
		},
		{
			name:     "only archived tasks",
			tasks:    []task{newTask(str, "archived", str, past)},
			label:    str,
			numCards: 0,
			err:      false,
		},
		{
			name:     "only future tasks",
			tasks:    []task{newTask(str, "active", str, future)},
			label:    str,
			numCards: 0,
			err:      false,
		},
		{
			name:     "missing label",
			tasks:    []task{newTask(str, "active", str, past)},
			label:    "",
			numCards: 0,
			err:      true,
		},
		{
			name:     "invalid date",
			tasks:    []task{newTask(str, "active", str, "xxx")},
			label:    str,
			numCards: 0,
			err:      true,
		},
		{
			name:     "all eligible",
			tasks:    []task{newTask(str, "active", str, past), newTask(str, "active", str, past)},
			label:    str,
			numCards: 2,
			err:      false,
		},
		{
			name: "some ineligible",
			tasks: []task{
				newTask(str, "active", str, past),
				newTask(str, "archived", str, past),
				newTask(str, "archived", str, past),
			},
			label:    str,
			numCards: 1,
			err:      false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := toCards(tc.tasks, tc.label)

			if (err != nil && !tc.err) || err == nil && tc.err {
				t.Fatalf("did not expect to get '%v' error", err)
			}

			if len(cards) != tc.numCards {
				t.Errorf("wanted %d cards to be created, got %d", tc.numCards, len(cards))
			}
		})
	}
}

func newTask(name, state, notes, due string) task {
	return task{
		Id:            0,
		Name:          name,
		State:         state,
		Color:         "",
		Notes:         notes,
		Period:        0,
		NextResetDate: due,
		MuteEmails:    0,
	}
}
