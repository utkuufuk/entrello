package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/pkg/trello"
)

func Poll(cfg config.RunnerConfig) error {
	loc, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		return fmt.Errorf("invalid timezone location: %v", loc)
	}

	services, labels := getServices(cfg.Services, time.Now().In(loc))
	if len(services) == 0 {
		return nil
	}

	client := trello.NewClient(cfg.Trello)

	if err := client.LoadBoard(labels); err != nil {
		return fmt.Errorf("Could not load existing cards from the board: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(services))
	for _, src := range services {
		go process(src, client, &wg)
	}
	wg.Wait()

	return nil
}
