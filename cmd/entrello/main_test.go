package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/utkuufuk/entrello/internal/config"
)

func TestGetEnabledSourcesAndLabels(t *testing.T) {
	tt := []struct {
		name         string
		githubIssues config.GithubIssues
		todoDock     config.TodoDock
		numResults   int
		labels       []string
	}{
		{
			name:         "nothing enabled",
			githubIssues: config.GithubIssues{Enabled: false},
			todoDock:     config.TodoDock{Enabled: false},
			numResults:   0,
		},
		{
			name: "only github issues enabled",
			githubIssues: config.GithubIssues{
				Enabled: true,
				Label:   "github-label",
			},
			todoDock:   config.TodoDock{Enabled: false},
			numResults: 1,
			labels:     []string{"github-label"},
		},
		{
			name:         "only tododock enabled",
			githubIssues: config.GithubIssues{Enabled: false},
			todoDock: config.TodoDock{
				Enabled: true,
				Label:   "tododock-label",
			},
			numResults: 1,
			labels:     []string{"tododock-label"},
		},
		{
			name: "all enabled",
			githubIssues: config.GithubIssues{
				Enabled: true,
				Label:   "github-label",
			},
			todoDock: config.TodoDock{
				Enabled: true,
				Label:   "tododock-label",
			},
			numResults: 2,
			labels:     []string{"github-label", "tododock-label"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Sources{
				GithubIssues: tc.githubIssues,
				TodoDock:     tc.todoDock,
			}

			sources, labels := getEnabledSourcesAndLabels(cfg)
			if len(sources) != tc.numResults {
				t.Errorf("expected %d source(s); got %v", tc.numResults, len(sources))
			}

			if len(labels) != tc.numResults {
				t.Errorf("expected %d label(s); got %v", tc.numResults, len(labels))
			}

			if tc.numResults == 0 {
				return
			}

			var opts cmp.Options
			if diff := cmp.Diff(labels, tc.labels, opts...); diff != "" {
				t.Errorf("labels diff: %s", diff)
			}
		})
	}
}
