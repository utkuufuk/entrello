package main

import (
	"flag"
	"log"
	"os"
	"sync"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/logger"
	"github.com/utkuufuk/entrello/pkg/trello"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "c", "config.json", "config file path")
	flag.Parse()

	cfg, err := config.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Could not read configuration: %v", err)
	}

	loc, err := time.LoadLocation(cfg.TimezoneLocation)
	if err != nil {
		logger.Error("Invalid timezone location: %v", loc)
		os.Exit(1)
	}

	sources, labels := getSources(cfg.Sources, time.Now().In(loc))
	if len(sources) == 0 {
		return
	}

	client, err := trello.NewClient(cfg.Trello)
	if err != nil {
		logger.Error("Could not create trello client: %v", err)
		os.Exit(1)
	}

	if err := client.LoadBoard(labels); err != nil {
		logger.Error("Could not load existing cards from the board: %v", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	wg.Add(len(sources))
	for _, src := range sources {
		go process(src, client, &wg)
	}
	wg.Wait()
}
