package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/github"
	"github.com/utkuufuk/entrello/internal/tododock"
	"github.com/utkuufuk/entrello/internal/trello"
)

// source defines an interface for a Trello card source
type source interface {
	// IsEnabled returns true if the source is enabled.
	IsEnabled() bool

	// IsStrict returns true if "strict" mode is enabled for the source
	IsStrict() bool

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
func getEnabledSourcesAndLabels(ctx context.Context, cfg config.Sources) (sources []source, labels []string) {
	arr := []source{
		github.GetSource(ctx, cfg.GithubIssues),
		tododock.GetSource(cfg.TodoDock),
	}
	now := time.Now()

	for _, src := range arr {
		if ok, err := shouldQuery(src, now); !ok {
			if err != nil {
				// @todo: send telegram notification instead if enabled
				log.Printf("[-] could not check if '%s' should be queried or not, skipping", src.GetName())
			}
			continue
		}
		sources = append(sources, src)
		labels = append(labels, src.GetLabel())
	}
	return sources, labels
}

// shouldQuery checks if the given source should be queried at the given time
func shouldQuery(src source, now time.Time) (bool, error) {
	if !src.IsEnabled() {
		return false, nil
	}

	interval := src.GetPeriod().Interval
	if interval < 0 {
		return false, fmt.Errorf("period interval must be a positive integer, got: '%d'", interval)
	}

	switch src.GetPeriod().Type {
	case config.PERIOD_TYPE_DEFAULT:
		return true, nil
	case config.PERIOD_TYPE_DAY:
		if interval > 31 {
			return false, fmt.Errorf("daily interval cannot be more than 14, got: '%d'", interval)
		}
		return now.Day()%interval == 0 && now.Hour() == 0 && now.Minute() == 0, nil
	case config.PERIOD_TYPE_HOUR:
		if interval > 23 {
			return false, fmt.Errorf("hourly interval cannot be more than 23, got: '%d'", interval)
		}
		return now.Hour()%interval == 0 && now.Minute() == 0, nil
	case config.PERIOD_TYPE_MINUTE:
		if interval > 60 {
			return false, fmt.Errorf("minute interval cannot be more than 60, got: '%d'", interval)
		}
		return now.Minute()%interval == 0, nil
	}

	return false, fmt.Errorf("unrecognized source period type: '%s'", src.GetPeriod().Type)
}

// queueActionables fetches new cards from the source, then pushes those to be created and
// to be deleted into the corresponding channels, as well as any errors encountered.
func queueActionables(src source, client trello.Client, q CardQueue) {
	cards, err := src.FetchNewCards()
	if err != nil {
		q.err <- fmt.Errorf("could not fetch cards for source '%s': %v", src.GetName(), err)
		return
	}

	new, stale := client.CompareWithExisting(cards, src.GetLabel())

	for _, c := range new {
		q.add <- c
	}

	if !src.IsStrict() {
		return
	}

	for _, c := range stale {
		q.del <- c
	}
}