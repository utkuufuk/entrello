package services

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

// getServicesToPoll returns a slice of services to poll & another slice of relevant service labels
func getServicesToPoll(
	serviceArr []config.Service,
	now time.Time,
) (
	services []config.Service,
	labels []string,
	err error,
) {
	for _, service := range serviceArr {
		if ok, err := shouldPoll(service, now); !ok {
			if err != nil {
				return services, labels, fmt.Errorf(
					"could not check if '%s' should be queried or not: %w",
					service.Name,
					err,
				)
			}
			continue
		}
		services = append(services, service)
		labels = append(labels, service.Label)
	}
	return services, labels, nil
}

// shouldPoll checks if a the service should be polled at the given time instant
func shouldPoll(service config.Service, date time.Time) (bool, error) {
	interval := service.Period.Interval
	if interval < 0 {
		return false, fmt.Errorf("period interval must be a positive integer, got: '%d'", interval)
	}

	switch service.Period.Type {
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

	return false, fmt.Errorf("unrecognized service period type: '%s'", service.Period.Type)
}

// poll polls the given service and creates Trello cards for each item unless
// a corresponding card already exists, also deletes the stale cards if strict mode is enabled
func poll(service config.Service, client trello.Client, wg *sync.WaitGroup) {
	defer wg.Done()

	req, err := http.NewRequest("GET", service.Endpoint, nil)
	if err != nil {
		logger.Error("could not create GET request to service '%s' endpoint: %v", service.Name, err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", service.Secret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error("could not make GET request to service '%s' endpoint: %v", service.Name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		msg := string(body)
		if err != nil {
			msg = err.Error()
		}
		logger.Error("could not retrieve cards from service '%s': %v", service.Name, msg)
		return
	}

	var cards []trello.Card
	if err = json.NewDecoder(resp.Body).Decode(&cards); err != nil {
		logger.Error("could not decode cards received from service '%s': %v", service.Name, err)
		return
	}

	new, stale := client.FilterNewAndStale(cards, service.Label)
	for _, c := range new {
		if err := client.CreateCard(c, service.Label, service.List); err != nil {
			logger.Error("could not create Trello card: %v", err)
			continue
		}
		logger.Info("created new card: %s", c.Name)
	}

	if !service.Strict {
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
