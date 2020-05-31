package main

import (
	"context"
	"fmt"
	"log"

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
	cards, err := src.FetchNewCards()
	if err != nil {
		q.err <- fmt.Errorf("could not fetch cards for source '%s': %v", src.GetName(), err)
		return
	}

	new, stale := client.CompareWithExisting(cards, src.GetLabel())

	for _, c := range new {
		q.add <- c
	}

	if !src.IsStrict() {
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
			// @todo: send telegram notification instead if enabled
			if err := client.CreateCard(c); err != nil {
				log.Printf("[-] error occurred while creating card: %v", err)
				break
			}
			log.Printf("[+] created new card: %s", c.Name)
		case c := <-q.del:
			// @todo: send telegram notification instead if enabled
			if err := client.ArchiveCard(c); err != nil {
				log.Printf("[-] error occurred while archiving card: %v", err)
				break
			}
			log.Printf("[+] archived stale card: %s", c.Name)
		case err := <-q.err:
			// @todo: send telegram notification instead if enabled
			log.Printf("[-] %v", err)
		case <-ctx.Done():
			return
		}
	}
}
