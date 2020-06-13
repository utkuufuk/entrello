package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/github"
	"github.com/utkuufuk/entrello/internal/habits"
	"github.com/utkuufuk/entrello/internal/tododock"
	"github.com/utkuufuk/entrello/internal/trello"
)

type source struct {
	cfg config.SourceConfig
	api interface {
		FetchNewCards(ctx context.Context, cfg config.SourceConfig) ([]trello.Card, error)
	}
}

// getEnabledSources returns a slice of enabled sources & their labels as a separate slice
func getEnabledSources(cfg config.Sources) (sources []source, labels []string) {
	arr := []source{
		{cfg: cfg.GithubIssues.SourceConfig, api: github.GetSource(cfg.GithubIssues)},
		{cfg: cfg.TodoDock.SourceConfig, api: tododock.GetSource(cfg.TodoDock)},
		{cfg: cfg.Habits.SourceConfig, api: habits.GetSource(cfg.Habits)},
	}

	now := time.Now()

	for _, src := range arr {
		if ok, err := src.shouldQuery(now); !ok {
			if err != nil {
				logger.Errorf("could not check if '%s' should be queried or not, skipping", src.cfg.Name)
			}
			continue
		}
		sources = append(sources, src)
		labels = append(labels, src.cfg.Label)
	}
	return sources, labels
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

// process fetches cards from the source and creates the ones that don't already exist,
// also deletes the stale cards if strict mode is enabled
func (s source) process(ctx context.Context, client trello.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	cards, err := s.api.FetchNewCards(ctx, s.cfg)
	if err != nil {
		logger.Errorf("could not fetch cards for source '%s': %v", s.cfg.Name, err)
		return
	}

	new, stale := client.FilterNewAndStale(cards, s.cfg.Label)

	for _, c := range new {
		if err := client.CreateCard(c); err != nil {
			logger.Errorf("could not create Trello card: %v", err)
			continue
		}
		logger.Printf("created new card: %s", c.Name)
	}

	if !s.cfg.Strict {
		return
	}

	for _, c := range stale {
		if err := client.DeleteCard(c); err != nil {
			logger.Errorf("could not delete Trello card: %v", err)
			continue
		}
		logger.Printf("deleted stale card: %s", c.Name)
	}
}
