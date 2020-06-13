package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/syslog"
	"github.com/utkuufuk/entrello/internal/trello"
)

var (
	logger syslog.Logger
)

func main() {
	// read config params
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("Could not read config variables: %v", err)
	}

	// get a system logger instance
	logger = syslog.NewLogger(cfg.Telegram)

	// set global timeout
	timeout := time.Second * time.Duration(cfg.TimeoutSeconds)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// get a list of enabled sources and the corresponding labels for each source
	sources, labels := getEnabledSources(cfg.Sources)
	if len(sources) == 0 {
		return
	}

	// initialize the Trello client
	client, err := trello.NewClient(cfg.Trello)
	if err != nil {
		logger.Fatalf("could not create trello client: %v", err)
	}

	// load Trello cards from the board with relevant labels
	if err := client.LoadBoard(labels); err != nil {
		logger.Fatalf("could not load existing cards from the board: %v", err)
	}

	// fetch new cards from each source and handle the new & stale ones
	var wg sync.WaitGroup
	wg.Add(len(sources))
	for _, src := range sources {
		go src.process(ctx, client, &wg)
	}
	wg.Wait()
}
