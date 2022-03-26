package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/pkg/trello"
)

func Poll(cfg config.Config) error {
	loc, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		return fmt.Errorf("invalid timezone location: %v", loc)

	}

	sources, labels := getSources(cfg.Sources, time.Now().In(loc))
	if len(sources) == 0 {
		return nil
	}

	client, err := trello.NewClient(cfg.Trello)
	if err != nil {
		return fmt.Errorf("could not create trello client: %v", err)
	}

	if err := client.LoadBoard(labels); err != nil {
		return fmt.Errorf("Could not load existing cards from the board: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(sources))
	for _, src := range sources {
		go process(src, client, &wg)
	}
	wg.Wait()

	return nil
}
