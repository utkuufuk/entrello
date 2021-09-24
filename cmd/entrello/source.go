package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/pkg/trello"
)

// getSources returns a slice of sources & their labels as a separate slice
func getSources(cfg config.Sources, now time.Time) (sources []config.Source, labels []string) {
	arr := []config.Source{
		cfg.GithubIssues,
		cfg.TodoDock,
		cfg.Habits,
	}

	for _, src := range arr {
		if ok, err := shouldQuery(src, now); !ok {
			if err != nil {
				logger.Errorf("could not check if '%s' should be queried or not, skipping", src.Name)
			}
			continue
		}
		sources = append(sources, src)
		labels = append(labels, src.Label)
	}
	return sources, labels
}

// shouldQuery checks if a the source should be queried at the given time
func shouldQuery(src config.Source, date time.Time) (bool, error) {
	interval := src.Period.Interval
	if interval < 0 {
		return false, fmt.Errorf("period interval must be a positive integer, got: '%d'", interval)
	}

	switch src.Period.Type {
	case config.PERIOD_TYPE_DEFAULT:
		return true, nil
	case config.PERIOD_TYPE_DAY:
		if interval > 31 {
			return false, fmt.Errorf("daily interval cannot be more than 14, got: '%d'", interval)
		}
		return date.Day()%interval == 0 && date.Hour() == 0 && date.Minute() == 0, nil
	case config.PERIOD_TYPE_HOUR:
		if interval > 23 {
			return false, fmt.Errorf("hourly interval cannot be more than 23, got: '%d'", interval)
		}
		return date.Hour()%interval == 0 && date.Minute() == 0, nil
	case config.PERIOD_TYPE_MINUTE:
		if interval > 60 {
			return false, fmt.Errorf("minute interval cannot be more than 60, got: '%d'", interval)
		}
		return date.Minute()%interval == 0, nil
	}

	return false, fmt.Errorf("unrecognized source period type: '%s'", src.Period.Type)
}

// process fetches cards from the source and creates the ones that don't already exist,
// also deletes the stale cards if strict mode is enabled
func process(src config.Source, client trello.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(src.Endpoint)
	if err != nil {
		logger.Errorf("could not make GET request to source '%s' endpoint: %v", src.Name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		msg := string(body)
		if err != nil {
			msg = err.Error()
		}
		logger.Errorf("could not retrieve cards from source '%s': %v", src.Name, msg)
		return
	}

	var cards []trello.Card
	if err = json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		logger.Errorf("could not decode cards received from source '%s': %v", src.Name, err)
		return
	}

	new, stale := client.FilterNewAndStale(cards, src.Label)
	for _, c := range new {
		if err := client.CreateCard(c, src.Label, src.List); err != nil {
			logger.Errorf("could not create Trello card: %v", err)
			continue
		}
		logger.Debugf("created new card: %s", c.Name)
	}

	if !src.Strict {
		return
	}

	for _, c := range stale {
		if err := client.DeleteCard(c); err != nil {
			logger.Errorf("could not delete Trello card: %v", err)
			continue
		}
		logger.Debugf("deleted stale card: %s", c.Name)
	}
}
