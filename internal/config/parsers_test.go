package config

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseServices(t *testing.T) {
	tt := []struct {
		name     string
		input    string
		isValid  bool
		services []Service
	}{
		{
			name:     "simple service without secret",
			input:    "label@http://example.com",
			isValid:  true,
			services: []Service{{Label: "label", Secret: "", Endpoint: "http://example.com"}},
		},
		{
			name:     "simple service with secret",
			input:    "label:secret@http://example.com",
			isValid:  true,
			services: []Service{{Label: "label", Secret: "secret", Endpoint: "http://example.com"}},
		},
		{
			name:     "service with secret containing numbers and uppercase letters",
			input:    "label:aBcD1230XyZ@http://example.com",
			isValid:  true,
			services: []Service{{Label: "label", Secret: "aBcD1230XyZ", Endpoint: "http://example.com"}},
		},
		{
			name:    "endpoint URL does not start with 'http'",
			input:   "label@example.com",
			isValid: false,
		},
		{
			name:    "no '@' delimiter",
			input:   "label-http://example.com",
			isValid: false,
		},
		{
			name:    "multiple '@' delimiters",
			input:   "label@joe@example.com",
			isValid: false,
		},
		{
			name:    "multiple ':' delimiters",
			input:   "label:super:secret:password@http://example.com",
			isValid: false,
		},
		{
			name:    "non-alphanumeric characters in label",
			input:   "definitely$$not_*?a+Trello.Label@http://example.com",
			isValid: false,
		},
		{
			name:    "non-alphanumeric characters in secret",
			input:   "label:-?_*@http://example.com",
			isValid: false,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			services, err := parseServices(tc.input)
			if tc.isValid != (err == nil) {
				t.Errorf("expected valid output? %v. Got error: %s", tc.isValid, err)
				return
			}
			if diff := cmp.Diff(services, tc.services); diff != "" {
				t.Errorf("services diff: %s", diff)
			}
		})
	}
}
