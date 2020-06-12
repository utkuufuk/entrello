package main

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

type source struct {
	cfg config.SourceConfig
	api interface {
		FetchNewCards(ctx context.Context, cfg config.SourceConfig) ([]trello.Card, error)
	}
}

// shouldQuery checks if a the source should be queried at the given time
func (s source) shouldQuery(now time.Time) (bool, error) {
	if !s.cfg.Enabled {
		return false, nil
	}

	interval := s.cfg.Period.Interval
	if interval < 0 {
		return false, fmt.Errorf("period interval must be a positive integer, got: '%d'", interval)
	}

	switch s.cfg.Period.Type {
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

	return false, fmt.Errorf("unrecognized source period type: '%s'", s.cfg.Period.Type)
}

// queueActionables fetches new cards from the source, then pushes those to be created and
// to be deleted into the corresponding channels, as well as any errors encountered.
func (s source) queueActionables(ctx context.Context, client trello.Client, q CardQueue) {
	cards, err := s.api.FetchNewCards(ctx, s.cfg)
	if err != nil {
		q.err <- fmt.Errorf("could not fetch cards for source '%s': %v", s.cfg.Name, err)
		return
	}

	new, stale := client.FilterNewAndStale(cards, s.cfg.Label)

	for _, c := range new {
		q.add <- c
	}

	if !s.cfg.Strict {
		return
	}

	for _, c := range stale {
		q.del <- c
	}
}
