package main

import (
	"context"
	"fmt"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/github"
	"github.com/utkuufuk/entrello/internal/tododock"
	"github.com/utkuufuk/entrello/internal/trello"
)

type source struct {
	cfg config.SourceConfig
	api interface {
		FetchNewCards() ([]trello.Card, error)
	}
}

// getEnabledSourcesAndLabels returns a list of enabled sources & all relevant label IDs
func getEnabledSourcesAndLabels(ctx context.Context, cfg config.Sources) (s []source, l []string) {
	sources := []source{
		{cfg.GithubIssues.SourceConfig, github.GetSource(ctx, cfg.GithubIssues)},
		{cfg.TodoDock.SourceConfig, tododock.GetSource(cfg.TodoDock)},
	}
	now := time.Now()

	for _, src := range sources {
		if ok, err := shouldQuery(src.cfg, now); !ok {
			if err != nil {
				logger.Errorf("could not check if '%s' should be queried or not, skipping", src.cfg.Name)
			}
			continue
		}
		s = append(s, src)
		l = append(l, src.cfg.Label)
	}
	return s, l
}

// shouldQuery checks if a query should be executed at the given time given the source configuration
func shouldQuery(cfg config.SourceConfig, now time.Time) (bool, error) {
	if !cfg.Enabled {
		return false, nil
	}

	interval := cfg.Period.Interval
	if interval < 0 {
		return false, fmt.Errorf("period interval must be a positive integer, got: '%d'", interval)
	}

	switch cfg.Period.Type {
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

	return false, fmt.Errorf("unrecognized source period type: '%s'", cfg.Period.Type)
}
