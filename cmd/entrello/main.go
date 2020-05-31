package main

import (
	"context"
	"log"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

type CardQueue struct {
	add chan trello.Card
	del chan trello.Card
	err chan error
}

func main() {
	// read config params
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		// @todo: send telegram notification instead if enabled
		log.Fatalf("[-] could not read config variables: %v", err)
	}

	// set global timeout
	timeout := time.Second * time.Duration(cfg.TimeoutSeconds)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// get a list of enabled sources and the corresponding labels for each source
	sources, labels := getEnabledSourcesAndLabels(ctx, cfg.Sources)
	if len(sources) == 0 {
		return
	}

	// initialize the Trello client
	client, err := trello.NewClient(cfg.Trello)
	if err != nil {
		// @todo: send telegram notification instead if enabled
		log.Fatalf("[-] could not create trello client: %v", err)
	}

	// within the Trello client, load the existing cards (only with relevant labels)
	if err := client.LoadCards(labels); err != nil {
		// @todo: send telegram notification instead if enabled
		log.Fatalf("[-] could not load existing cards from the board: %v", err)
	}

	// initialize channels, then start listening each source for cards to create/delete and errors
	q := CardQueue{
		add: make(chan trello.Card),
		del: make(chan trello.Card),
		err: make(chan error),
	}

	// concurrently fetch new cards from each source and start queuing cards to be created & deleted
	for _, src := range sources {
		go queueActionables(src, client, q)
	}

	//
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
