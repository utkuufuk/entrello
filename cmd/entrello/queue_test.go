package main

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/utkuufuk/entrello/internal/trello"
)

func TestProcessActionables(t *testing.T) {
	tt := []struct {
		name    string
		timeout int
	}{
		{
			name:    "no timeout",
			timeout: 0,
		},
		{
			name:    "1 second timeout",
			timeout: 1,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(tc.timeout))
			defer cancel()

			q := CardQueue{make(chan trello.Card), make(chan trello.Card), make(chan error)}
			processActionables(ctx, trello.Client{}, q)

			if secs := int(math.Round(time.Since(start).Seconds())); secs != tc.timeout {
				t.Errorf("wanted process to take %d seconds, it took %d", tc.timeout, secs)
			}
		})
	}
}
