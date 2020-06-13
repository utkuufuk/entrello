package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/utkuufuk/entrello/internal/config"
)

func TestGetEnabledSources(t *testing.T) {
	period := config.Period{
		Type:     config.PERIOD_TYPE_DEFAULT,
		Interval: 0,
	}

	tt := []struct {
		name            string
		githubIssuesCfg config.SourceConfig
		todoDockCfg     config.SourceConfig
		numResults      int
		labels          []string
	}{
		{
			name:            "nothing enabled",
			githubIssuesCfg: config.SourceConfig{Enabled: false, Period: period},
			todoDockCfg:     config.SourceConfig{Enabled: false, Period: period},
			numResults:      0,
		},
		{
			name: "only github issues enabled",
			githubIssuesCfg: config.SourceConfig{
				Enabled: true,
				Period:  period,
				Label:   "github-label",
			},
			todoDockCfg: config.SourceConfig{Enabled: false, Period: period},
			numResults:  1,
			labels:      []string{"github-label"},
		},
		{
			name:            "only tododock enabled",
			githubIssuesCfg: config.SourceConfig{Enabled: false, Period: period},
			todoDockCfg: config.SourceConfig{
				Enabled: true,
				Period:  period,
				Label:   "tododock-label",
			},
			numResults: 1,
			labels:     []string{"tododock-label"},
		},
		{
			name: "all enabled",
			githubIssuesCfg: config.SourceConfig{
				Enabled: true,
				Period:  period,
				Label:   "github-label",
			},
			todoDockCfg: config.SourceConfig{
				Enabled: true,
				Period:  period,
				Label:   "tododock-label",
			},
			numResults: 2,
			labels:     []string{"github-label", "tododock-label"},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.Sources{
				GithubIssues: config.GithubIssues{SourceConfig: tc.githubIssuesCfg},
				TodoDock:     config.TodoDock{SourceConfig: tc.todoDockCfg},
			}

			sources, labels := getEnabledSources(cfg)
			if len(sources) != tc.numResults {
				t.Errorf("expected %d source(s); got %v", tc.numResults, len(sources))
			}

			if len(labels) != tc.numResults {
				t.Errorf("expected %d label(s); got %v", tc.numResults, len(labels))
			}

			if tc.numResults == 0 {
				return
			}

			if diff := cmp.Diff(labels, tc.labels); diff != "" {
				t.Errorf("labels diff: %s", diff)
			}
		})
	}
}

func TestShouldQuery(t *testing.T) {
	tt := []struct {
		name      string
		pType     string
		pInterval int
		date      time.Time
		ok        bool
		err       error
	}{
		{
			name:      "default period",
			pType:     "default",
			pInterval: 0,
			date:      time.Now(),
			ok:        true,
			err:       nil,
		},
		{
			name:      "invalid period type",
			pType:     "foo",
			pInterval: 0,
			date:      time.Now(),
			ok:        false,
			err:       fmt.Errorf("unrecognized source period type: 'foo'"),
		},
		{
			name:      "negative period interval",
			pType:     config.PERIOD_TYPE_DAY,
			pInterval: -1,
			date:      time.Now(),
			ok:        false,
			err:       fmt.Errorf("period interval must be a positive integer, got: '-1'"),
		},
		{
			name:      "every 3 days, on 6th at midnight, should query",
			pType:     config.PERIOD_TYPE_DAY,
			pInterval: 3,
			date:      time.Date(1990, time.Month(2), 6, 0, 0, 0, 0, time.UTC),
			ok:        true,
			err:       nil,
		},
		{
			name:      "every 3 days, on 6th at 01:00, should not query",
			pType:     config.PERIOD_TYPE_DAY,
			pInterval: 3,
			date:      time.Date(1990, time.Month(2), 6, 1, 0, 0, 0, time.UTC),
			ok:        false,
			err:       nil,
		},
		{
			name:      "every 3 days, on 6th at 00:15, should not query",
			pType:     config.PERIOD_TYPE_DAY,
			pInterval: 3,
			date:      time.Date(1990, time.Month(2), 6, 0, 15, 0, 0, time.UTC),
			ok:        false,
			err:       nil,
		},
		{
			name:      "every 3 days, on 4th, should not query",
			pType:     config.PERIOD_TYPE_DAY,
			pInterval: 3,
			date:      time.Date(1990, time.Month(2), 4, 0, 0, 0, 0, time.UTC),
			ok:        false,
			err:       nil,
		},
		{
			name:      "invalid daily period interval",
			pType:     config.PERIOD_TYPE_DAY,
			pInterval: 40,
			date:      time.Date(1990, time.Month(2), 4, 0, 0, 0, 0, time.UTC),
			ok:        false,
			err:       fmt.Errorf("daily interval cannot be more than 14, got: '40'"),
		},
		{
			name:      "every 5 hours, at 15:00, should query",
			pType:     config.PERIOD_TYPE_HOUR,
			pInterval: 5,
			date:      time.Date(1990, time.Month(2), 1, 15, 0, 0, 0, time.UTC),
			ok:        true,
			err:       nil,
		},
		{
			name:      "every 4 hours, at 16:33, should not query",
			pType:     config.PERIOD_TYPE_HOUR,
			pInterval: 4,
			date:      time.Date(1990, time.Month(2), 4, 16, 33, 0, 0, time.UTC),
			ok:        false,
			err:       nil,
		},
		{
			name:      "every 2 hours, at 21:00, should not query",
			pType:     config.PERIOD_TYPE_HOUR,
			pInterval: 2,
			date:      time.Date(1990, time.Month(2), 4, 21, 0, 0, 0, time.UTC),
			ok:        false,
			err:       nil,
		},
		{
			name:      "invalid hourly period interval",
			pType:     config.PERIOD_TYPE_HOUR,
			pInterval: 25,
			date:      time.Date(1990, time.Month(2), 4, 1, 0, 0, 0, time.UTC),
			ok:        false,
			err:       fmt.Errorf("hourly interval cannot be more than 23, got: '25'"),
		},
		{
			name:      "every 7 minutes, at 14:56, should query",
			pType:     config.PERIOD_TYPE_MINUTE,
			pInterval: 7,
			date:      time.Date(1990, time.Month(2), 1, 14, 56, 0, 0, time.UTC),
			ok:        true,
			err:       nil,
		},
		{
			name:      "every 6 minutes, at 02:13, but should not query",
			pType:     config.PERIOD_TYPE_MINUTE,
			pInterval: 6,
			date:      time.Date(1990, time.Month(2), 4, 2, 13, 0, 0, time.UTC),
			ok:        false,
			err:       nil,
		},
		{
			name:      "invalid minute period interval",
			pType:     config.PERIOD_TYPE_MINUTE,
			pInterval: 61,
			date:      time.Date(1990, time.Month(2), 4, 1, 0, 0, 0, time.UTC),
			ok:        false,
			err:       fmt.Errorf("minute interval cannot be more than 60, got: '61'"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cfg := config.SourceConfig{
				Enabled: true,
				Period: config.Period{
					Type:     tc.pType,
					Interval: tc.pInterval,
				},
			}
			src := source{cfg: cfg}
			ok, err := src.shouldQuery(tc.date)

			if err != nil || tc.err != nil {
				if err == nil || tc.err == nil || err.Error() != tc.err.Error() {
					t.Errorf("expected error to be %v, got '%v'", tc.err, err)
				}
			}

			if ok != tc.ok {
				t.Errorf("expected %t, got %t", tc.ok, ok)
			}
		})
	}
}
