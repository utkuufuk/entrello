package main

import (
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
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.yml", "config file path")
	flag.Parse()

	cfg, err := config.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read configuration: %v", err)
	}

	logger = syslog.NewLogger(cfg.Telegram)

	loc, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		logger.Fatalf("Invalid timezone location: %v", loc)
	}

	sources, labels := getSources(cfg.Sources, time.Now().In(loc))
	if len(sources) == 0 {
		return
	}

	client, err := trello.NewClient(cfg.Trello)
	if err != nil {
		logger.Fatalf("Could not create trello client: %v", err)
	}

	if err := client.LoadBoard(labels); err != nil {
		logger.Fatalf("Could not load existing cards from the board: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(sources))
	for _, src := range sources {
		go process(src, client, &wg)
	}
	wg.Wait()
}
