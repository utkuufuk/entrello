package main

import (
	"log"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		// @todo: send telegram notification instead if enabled
		log.Fatalf("[-] could not read config variables: %v", err)
	}

	sources, labels := getEnabledSourcesAndLabels(cfg.Sources)
	if len(sources) == 0 {
		log.Println("[+] no sources enabled, aborting...")
		return
	}

	client, err := trello.NewClient(cfg)
	if err != nil {
		// @todo: send telegram notification instead if enabled
		log.Fatalf("[-] could not create trello client: %v", err)
	}

	if err := client.LoadExistingCards(labels); err != nil {
		// @todo: send telegram notification instead if enabled
		log.Fatalf("[-] could not load existing cards from the board: %v", err)
	}

	for _, src := range sources {
		cards, err := src.FetchNewCards()
		if err != nil {
			// @todo: send telegram notification instead if enabled
			log.Printf("[-] could not get cards for source '%s': %v", src.GetName(), err)
			continue
		}

		if err := client.UpdateCards(cards); err != nil {
			// @todo: send telegram notification instead if enabled
			log.Printf("[-] error occurred while processing source '%s': %v", src.GetName(), err)
		}
	}
}
