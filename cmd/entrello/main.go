package main

import (
	"context"
	"log"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/github"
	"github.com/utkuufuk/entrello/internal/tododock"
	"github.com/utkuufuk/entrello/internal/trello"
)

// Source represents a card source which exports a name and a getter for the cards to be created
type Source interface {
	// GetCards returns a list of Trello cards to be inserted into the board from the source
	GetCards() ([]trello.Card, error)
}

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("[-] could not read config variables: %v", err)
	}

	sources := collectSources(cfg.Sources)
	if len(sources) == 0 {
		log.Println("[+] no sources enabled, aborting...")
		return
	}

	// fetch all existing cards in the board with their corresponding labels
	client := trello.NewClient(cfg)
	cardMap, err := client.FetchBoardCards()
	if err != nil {
		log.Fatalf("[-] could not fetch cards in Tasks board: %v", err)
	}

	for name, source := range sources {
		cards, err := source.GetCards()
		if err != nil {
			log.Printf("[-] could not get cards from source '%s': %v", name, err)
		}

		for _, card := range cards {
			// if card name already exists with the same label, do not create a duplicate one
			if labels, ok := cardMap[card.Name]; ok {
				if contains(labels, card.LabelId) {
					continue
				}
			}

			err = client.AddCard(card)
			if err != nil {
				log.Printf("[-] could not create card '%s': %v", card.Name, err)
			}
			log.Printf("[+] created new card: '%s'\n", card.Name)
		}
	}
}

// collectSources populates & returns a map of card sources to be iterated over
func collectSources(cfg config.Sources) (s map[string]Source) {
	s = make(map[string]Source)

	if cfg.TodoDock.Enabled {
		s["TodoDock"] = tododock.GetSource(cfg.TodoDock)
	}

	if cfg.GithubIssues.Enabled {
		s["Github Issues"] = github.GetSource(context.Background(), cfg.GithubIssues)
	}

	return s
}

// contains returns true if the list of label IDs contain the given label ID
func contains(labels []string, label string) bool {
	for _, l := range labels {
		if l == label {
			return true
		}
	}
	return false
}
