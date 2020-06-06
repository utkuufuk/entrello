package main

import (
	"context"
	"fmt"

	"github.com/utkuufuk/entrello/internal/trello"
)

type CardQueue struct {
	add chan trello.Card
	del chan trello.Card
	err chan error
}

// queueActionables fetches new cards from the source, then pushes those to be created and
// to be deleted into the corresponding channels, as well as any errors encountered.
func queueActionables(src source, client trello.Client, q CardQueue) {
	cards, err := src.api.FetchNewCards(src.cfg)
	if err != nil {
		q.err <- fmt.Errorf("could not fetch cards for source '%s': %v", src.cfg.Name, err)
		return
	}

	new, stale := client.CompareWithExisting(cards, src.cfg.Label)

	for _, c := range new {
		q.add <- c
	}

	if !src.cfg.Strict {
		return
	}

	for _, c := range stale {
		q.del <- c
	}
}

// processActionables listens to the card queue in an infinite loop and creates/deletes Trello cards
// depending on which channel the cards come from. Terminates whenever the global timeout is reached.
func processActionables(ctx context.Context, client trello.Client, q CardQueue) {
	for {
		select {
		case c := <-q.add:
			if err := client.CreateCard(c); err != nil {
				logger.Errorf("could not create Trello card: %v", err)
				break
			}
			logger.Printf("created new card: %s", c.Name)
		case c := <-q.del:
			if err := client.ArchiveCard(c); err != nil {
				logger.Errorf("could not archive card card: %v", err)
				break
			}
			logger.Printf("archived stale card: %s", c.Name)
		case err := <-q.err:
			logger.Errorf("%v", err)
		case <-ctx.Done():
			return
		}
	}
}
