package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/logger"
	"github.com/utkuufuk/entrello/pkg/trello"
)

// getSources returns a slice of sources & all source labels as a separate slice
func getSources(srcArr []config.Source, now time.Time) (sources []config.Source, labels []string) {
	for _, src := range srcArr {
		if ok, err := shouldQuery(src, now); !ok {
			if err != nil {
				logger.Error("could not check if '%s' should be queried or not, skipping", src.Name)
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
	case config.PeriodTypeDefault:
		return true, nil
	case config.PeriodTypeDay:
		if interval > 31 {
			return false, fmt.Errorf("daily interval cannot be more than 14, got: '%d'", interval)
		}
		return date.Day()%interval == 0 && date.Hour() == 0 && date.Minute() == 0, nil
	case config.PeriodTypeHour:
		if interval > 23 {
			return false, fmt.Errorf("hourly interval cannot be more than 23, got: '%d'", interval)
		}
		return date.Hour()%interval == 0 && date.Minute() == 0, nil
	case config.PeriodTypeMinute:
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
		logger.Error("could not make GET request to source '%s' endpoint: %v", src.Name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		msg := string(body)
		if err != nil {
			msg = err.Error()
		}
		logger.Error("could not retrieve cards from source '%s': %v", src.Name, msg)
		return
	}

	var cards []trello.Card
	if err = json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		logger.Error("could not decode cards received from source '%s': %v", src.Name, err)
		return
	}

	new, stale := client.FilterNewAndStale(cards, src.Label)
	for _, c := range new {
		if err := client.CreateCard(c, src.Label, src.List); err != nil {
			logger.Error("could not create Trello card: %v", err)
			continue
		}
		logger.Info("created new card: %s", c.Name)
	}

	if !src.Strict {
		return
	}

	for _, c := range stale {
		if err := client.DeleteCard(c); err != nil {
			logger.Error("could not delete Trello card: %v", err)
			continue
		}
		logger.Info("deleted stale card: %s", c.Name)
	}
}
