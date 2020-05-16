package main

import (
	"log"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/tododock"
	"github.com/utkuufuk/entrello/internal/trello"
)

type Source interface {
	GetCards() ([]trello.Card, error)
	GetName() string
}

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("[-] could not read config variables: %v", err)
	}

	sources, err := collectSources(cfg.Sources)
	if err != nil {
		log.Fatalf("[-] could not collect trello card sources: %v", err)
	}

	client := trello.NewClient(cfg)
	cardMap, err := client.FetchBoardCards()

	if err != nil {
		log.Fatalf("[-] could not fetch cards in Tasks board: %v", err)
	}

	for _, source := range sources {
		cards, err := source.GetCards()
		if err != nil {
			log.Printf("[-] could not get cards from source '%s': %v", source.GetName(), err)
		}

		for _, card := range cards {
			if _, ok := cardMap[card.Name]; ok {
				log.Printf("[+] skipping '%s' as it already exists...\n", card.Name)
				continue
			}

			err = client.AddCard(card)
			if err != nil {
				log.Printf("[-] could not create card '%s': %v", card.Name, err)
			}
			log.Printf("[+] created new card: '%s'\n", card.Name)
		}
	}
}

func collectSources(cfg config.Sources) ([]Source, error) {
	s := make([]Source, 0)
	s = append(s, tododock.GetSource(cfg.TodoDock))
	return s, nil
}
