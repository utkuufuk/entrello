package main

import (
	"context"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/github"
	"github.com/utkuufuk/entrello/internal/tododock"
	"github.com/utkuufuk/entrello/internal/trello"
)

// Source defines an interface for a Trello card source
type Source interface {
	// IsEnabled returns true if the source is enabled by configuration.
	IsEnabled() bool

	// GetName returns a human-readable name of the source
	GetName() string

	// GetLabel returns the corresponding card label ID for the source
	GetLabel() string

	// GetPeriod returns the period in minutes that the source should be checked
	GetPeriod() config.Period

	// FetchNewCards returns a list of Trello cards to be inserted into the board from the source
	FetchNewCards() ([]trello.Card, error)
}

// getEnabledSourcesAndLabels returns a list of enabled sources & all relevant label IDs
func getEnabledSourcesAndLabels(cfg config.Sources) (sources []Source, labels []string) {
	arr := []Source{
		github.GetSource(context.Background(), cfg.GithubIssues),
		tododock.GetSource(cfg.TodoDock),
	}
	now := time.Now()

	for _, src := range arr {
		if !src.IsEnabled() || !shouldCheck(now, src.GetPeriod()) {
			continue
		}
		sources = append(sources, src)
		labels = append(labels, src.GetLabel())
	}
	return sources, labels
}

// @todo: implement
// shouldCheck returns true if the given time instance is a valid point in time for checking the source
func shouldCheck(now time.Time, period config.Period) bool {
	// m := now.Minute()
	return true
}
