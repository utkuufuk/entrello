package trello

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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
