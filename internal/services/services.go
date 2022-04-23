package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/pkg/trello"
	"golang.org/x/exp/slices"
)

// Poll polls any number of configured services that should be polled at the given time instant.
func Poll(cfg config.RunnerConfig) error {
	loc, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		return fmt.Errorf("invalid timezone location: %v", loc)
	}

	services, labels, err := getServicesToPoll(cfg.Services, time.Now().In(loc))
	if err != nil {
		return fmt.Errorf("failed to get services to poll: %w", err)
	}
	if len(services) == 0 {
		return nil
	}

	client := trello.NewClient(cfg.Trello)

	if err := client.LoadBoard(labels); err != nil {
		return fmt.Errorf("Could not load existing cards from the board: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, src := range services {
		go poll(src, client, &wg)
	}
	wg.Wait()

	return nil
}

// Notify notifies any number of configured services with the latest state of the
// given archived Trello card.
func Notify(card trello.Card, services []config.Service) error {
	labelIds := make([]string, 0)
	for _, label := range card.Labels {
		labelIds = append(labelIds, label.ID)
	}

	for _, service := range services {
		if slices.Contains(labelIds, service.Label) {
			postBody, err := json.Marshal(card)
			if err != nil {
				return fmt.Errorf("could not marshal archived card: %w", err)
			}

			req, err := http.NewRequest("POST", service.Endpoint, bytes.NewBuffer(postBody))
			if err != nil {
				return fmt.Errorf("could not create POST request to %s: %w", service.Endpoint, err)
			}
			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("X-Api-Key", service.Secret)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return fmt.Errorf("could not post archived card data to %s: %w", service.Endpoint, err)
			}
			defer resp.Body.Close()
		}
	}
	return nil
}
