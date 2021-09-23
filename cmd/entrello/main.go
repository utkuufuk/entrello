package main

import (
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/syslog"
	"github.com/utkuufuk/entrello/pkg/trello"
)

var (
	logger syslog.Logger
	now    time.Time
)

func main() {
	// read configuration file path
	var configFile string
	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.Parse()

	// read config params
	cfg, err := config.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read configuration: %v", err)
	}

	// get a system logger instance
	logger = syslog.NewLogger(cfg.Telegram)

	// get current time for the configured location
	loc, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		logger.Fatalf("Invalid timezone location: %v", loc)
	}
	now = time.Now().In(loc)

	// set global timeout
	timeout := time.Second * time.Duration(cfg.TimeoutSeconds)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// get a list of sources and the corresponding labels for each source
	sources, labels := getSources(cfg.Sources)
	if len(sources) == 0 {
		return
	}

	// initialize the Trello client
	client, err := trello.NewClient(cfg.Trello)
	if err != nil {
		logger.Fatalf("Could not create trello client: %v", err)
	}

	// load Trello cards from the board with relevant labels
	if err := client.LoadBoard(labels); err != nil {
		logger.Fatalf("Could not load existing cards from the board: %v", err)
	}

	// fetch new cards from each source and handle the new & stale ones
	var wg sync.WaitGroup
	wg.Add(len(sources))
	for _, src := range sources {
		go process(src, ctx, client, &wg)
	}
	wg.Wait()
}
