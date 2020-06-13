package main

import (
	"context"
	"log"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/syslog"
	"github.com/utkuufuk/entrello/internal/trello"
)

var (
	logger syslog.Logger
)

type CardQueue struct {
	new   chan trello.Card
	stale chan trello.Card
}

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

	// concurrently fetch new cards from each source and start processing new & stale cards
	q := CardQueue{make(chan trello.Card), make(chan trello.Card)}
	for _, src := range sources {
		go src.queryAndQueue(ctx, client, q)
	}
	processActionables(ctx, client, q)
}

// processActionables listens to the card queue in an infinite loop and creates/deletes Trello cards
// depending on which channel the cards come from. Terminates whenever the global timeout is reached.
func processActionables(ctx context.Context, client trello.Client, q CardQueue) {
	for {
		select {
		case c := <-q.new:
			if err := client.CreateCard(c); err != nil {
				logger.Errorf("could not create Trello card: %v", err)
				break
			}
			logger.Printf("created new card: %s", c.Name)
		case c := <-q.stale:
			if err := client.DeleteCard(c); err != nil {
				logger.Errorf("could not delete Trello card: %v", err)
				break
			}
			logger.Printf("deleted stale card: %s", c.Name)
		case <-ctx.Done():
			return
		}
	}
}
