package github

import (
	"errors"
	"testing"

	"github.com/google/go-github/github"
)

func TestToCards(t *testing.T) {
	str := "test"

	tt := []struct {
		name     string
		label    string
		issues   []*github.Issue
		numCards int
		err      error
	}{
		{
			name:     "pull request",
			label:    str,
			issues:   []*github.Issue{{PullRequestLinks: &github.PullRequestLinks{}}},
			numCards: 0,
			err:      nil,
		},
		{
			name:     "valid issue",
			label:    str,
			issues:   []*github.Issue{newIssue(str, str, str)},
			numCards: 1,
			err:      nil,
		},
		{
			name:     "two valid issues",
			label:    str,
			issues:   []*github.Issue{newIssue(str, str, str), newIssue(str, str, str)},
			numCards: 2,
			err:      nil,
		},
		{
			name:     "empty label",
			label:    "",
			issues:   []*github.Issue{newIssue(str, str, str)},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:     "empty URL",
			label:    str,
			issues:   []*github.Issue{newIssue(str, str, "")},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:     "empty repo name",
			label:    str,
			issues:   []*github.Issue{newIssue(str, "", str)},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:     "empty title",
			label:    str,
			issues:   []*github.Issue{newIssue("", str, str)},
			numCards: 0,
			err:      errors.New(""),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cards, err := toCards(tc.issues, tc.label)
			if same := (err == nil && tc.err == nil) || tc.err != nil && err != nil; !same {
				t.Fatalf("want '%v', got '%v'", tc.err, err)
			}

			if len(cards) != tc.numCards {
				t.Errorf("expected %d cards, got %d", tc.numCards, len(cards))
			}
		})
	}
}

func newIssue(title, repoName, url string) *github.Issue {
	return &github.Issue{
		Title:      &title,
		Repository: &github.Repository{Name: &repoName},
		URL:        &url,
	}
}
