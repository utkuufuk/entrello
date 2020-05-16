package main

import (
	"log"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/tododock"
	"github.com/utkuufuk/entrello/internal/trello"
)

// Source represents a card source which exports a name and a getter for the cards to be created
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
		log.Fatalf("[-] could not collect card sources: %v", err)
	}

	// fetch all existing cards in the board with the "TodoDock" label
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
			// do not create cards with duplicate names if they both have the "TodoDock" label
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

// collectSources populates & returns an array of card sources to be iterated over
func collectSources(cfg config.Sources) ([]Source, error) {
	s := make([]Source, 0)
	s = append(s, tododock.GetSource(cfg.TodoDock))
	return s, nil
}
