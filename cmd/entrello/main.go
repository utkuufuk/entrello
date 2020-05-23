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
	// GetName returns a human-readable name of the source
	GetName() string

	// GetLabel returns the corresponding card label ID for the source
	GetLabel() string

	// GetNewCards returns a list of Trello cards to be inserted into the board from the source
	GetNewCards() ([]trello.Card, error)
}

func main() {
	cfg, err := config.ReadConfig("config.yml")
	if err != nil {
		log.Fatalf("[-] could not read config variables: %v", err)
	}

	sources, labels := getEnabledSourcesAndLabels(cfg.Sources)
	if len(sources) == 0 {
		log.Println("[+] no sources enabled, aborting...")
		return
	}

	client := trello.NewClient(cfg)
	if err := client.LoadExistingCardNames(labels); err != nil {
		// @todo: send telegram notification instead if enabled
		log.Fatalf("[-] could not load existing cards from the board: %v", err)
	}

	for _, src := range sources {
		cards, err := src.GetNewCards()
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

// getEnabledSourcesAndLabels returns a list of enabled sources & all relevant label IDs
func getEnabledSourcesAndLabels(cfg config.Sources) (sources []Source, labels []string) {
	if cfg.TodoDock.Enabled {
		src := tododock.GetSource(cfg.TodoDock)
		sources = append(sources, src)
		labels = append(labels, src.GetLabel())
	}

	if cfg.GithubIssues.Enabled {
		src := github.GetSource(context.Background(), cfg.GithubIssues)
		sources = append(sources, src)
		labels = append(labels, src.GetLabel())
	}

	return sources, labels
}
