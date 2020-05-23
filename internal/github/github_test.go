package github

import (
	"errors"
	"testing"

	"github.com/google/go-github/github"
)

func TestToCards(t *testing.T) {
	testStr := "test"

	tt := []struct {
		name     string
		label    string
		issues   []*github.Issue
		numCards int
		err      error
	}{
		{
			name:     "pull request",
			label:    testStr,
			issues:   []*github.Issue{{PullRequestLinks: &github.PullRequestLinks{}}},
			numCards: 0,
			err:      nil,
		},
		{
			name:     "valid issue",
			label:    testStr,
			issues:   []*github.Issue{createIssue(testStr, testStr, testStr)},
			numCards: 1,
			err:      nil,
		},
		{
			name:  "two valid issues",
			label: testStr,
			issues: []*github.Issue{
				createIssue(testStr, testStr, testStr),
				createIssue(testStr, testStr, testStr),
			},
			numCards: 2,
			err:      nil,
		},
		{
			name:     "empty label",
			label:    "",
			issues:   []*github.Issue{createIssue(testStr, testStr, testStr)},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:     "empty URL",
			label:    testStr,
			issues:   []*github.Issue{createIssue(testStr, testStr, "")},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:     "empty repo name",
			label:    testStr,
			issues:   []*github.Issue{createIssue(testStr, "", testStr)},
			numCards: 0,
			err:      errors.New(""),
		},
		{
			name:     "empty title",
			label:    testStr,
			issues:   []*github.Issue{createIssue("", testStr, testStr)},
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

func createIssue(title, repoName, url string) *github.Issue {
	return &github.Issue{
		Title:      &title,
		Repository: &github.Repository{Name: &repoName},
		URL:        &url,
	}
}
