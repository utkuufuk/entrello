package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
	"golang.org/x/oauth2"
)

type GithubIssuesSource struct {
	client *github.Client
	ctx    context.Context
	cfg    config.GithubIssues
}

func GetSource(ctx context.Context, cfg config.GithubIssues) GithubIssuesSource {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return GithubIssuesSource{client, ctx, cfg}
}

func (g GithubIssuesSource) FetchNewCards() ([]trello.Card, error) {
	issues, _, err := g.client.Issues.List(g.ctx, false, nil)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve issues: %w", err)
	}

	return toCards(issues, g.cfg.SourceConfig.Label)
}

// toCards converts a list of issues into a list of trello card
func toCards(issues []*github.Issue, label string) ([]trello.Card, error) {
	cards := make([]trello.Card, 0, len(issues))
	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}

		c, err := toCard(issue, label)
		if err != nil {
			return nil, fmt.Errorf("could not create card: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}

// toCard converts the given issue into a trello card
func toCard(issue *github.Issue, label string) (c trello.Card, err error) {
	if *issue.Repository.Name == "" || *issue.Title == "" || *issue.URL == "" || label == "" {
		return c, fmt.Errorf("could not convert issue to card, title, repo name, url and label cannot be blank")
	}
	name := fmt.Sprintf("[%s] %s", *issue.Repository.Name, *issue.Title)
	url := strings.Replace(*issue.URL, "api.", "", 1)
	url = strings.Replace(url, "/repos", "", 1)
	return trello.NewCard(name, label, url, nil)
}
