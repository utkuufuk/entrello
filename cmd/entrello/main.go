package main

import (
	"context"
	"log"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

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

	// concurrently fetch new cards from sources and start processing cards to be created & deleted
	q := CardQueue{make(chan trello.Card), make(chan trello.Card), make(chan error)}
	for _, src := range sources {
		go queueActionables(src, client, q)
	}
	processActionables(ctx, client, q)
}
