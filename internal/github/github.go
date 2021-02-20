package github

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/pkg/trello"
	"golang.org/x/oauth2"
)

type source struct {
	client *github.Client
}

func GetSource(cfg config.GithubIssues) source {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)
	return source{client}
}

func (s source) FetchNewCards(ctx context.Context, cfg config.SourceConfig, now time.Time) ([]trello.Card, error) {
	issues, _, err := s.client.Issues.List(ctx, false, nil)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve issues: %w", err)
	}

	return toCards(issues, cfg.Label)
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
			return nil, fmt.Errorf("could not create github issue card: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}

// toCard converts the given issue into a trello card
func toCard(issue *github.Issue, label string) (c trello.Card, err error) {
	if *issue.Repository.Name == "" || *issue.Title == "" || *issue.URL == "" || label == "" {
		e := "could not create card from issue; title, repo name, url and label are mandatory"
		return c, fmt.Errorf(e)
	}
	name := fmt.Sprintf("[%s] %s", *issue.Repository.Name, *issue.Title)
	url := strings.Replace(*issue.URL, "api.", "", 1)
	url = strings.Replace(url, "/repos", "", 1)
	return trello.NewCard(name, label, url, nil)
}
