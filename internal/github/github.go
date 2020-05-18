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
	client  *github.Client
	ctx     context.Context
	labelId string
}

func GetSource(ctx context.Context, cfg config.GithubIssues) GithubIssuesSource {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: cfg.Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return GithubIssuesSource{client, ctx, cfg.LabelId}
}

func (g GithubIssuesSource) GetCards() (cards []trello.Card, err error) {
	issues, _, err := g.client.Issues.List(g.ctx, false, nil)
	if err != nil {
		return cards, fmt.Errorf("could not retrieve issues: %w", err)
	}

	cards = make([]trello.Card, 0, len(issues))
	for _, issue := range issues {
		if issue.IsPullRequest() {
			continue
		}

		// convert API url to web URL
		url := strings.Replace(*issue.URL, "api.", "", 1)
		url = strings.Replace(url, "/repos", "", 1)

		cards = append(cards, trello.Card{
			Name:        fmt.Sprintf("[%s] %s", *issue.Repository.Name, *issue.Title),
			Description: url,
			LabelId:     g.labelId,
			DueDate:     nil,
		})
	}
	return cards, nil
}
