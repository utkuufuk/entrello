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
			input:    "label@endpoint",
			isValid:  true,
			services: []Service{{Label: "label", Secret: "", Endpoint: "endpoint"}},
		},
		{
			name:     "simple service with secret",
			input:    "label:secret@endpoint",
			isValid:  true,
			services: []Service{{Label: "label", Secret: "secret", Endpoint: "endpoint"}},
		},
		{
			name:     "service with secret containing numbers and uppercase letters",
			input:    "label:aBcD1230XyZ@endpoint",
			isValid:  true,
			services: []Service{{Label: "label", Secret: "aBcD1230XyZ", Endpoint: "endpoint"}},
		},
		{
			name:    "no '@' delimiter",
			input:   "label-endpoint",
			isValid: false,
		},
		{
			name:    "multiple '@' delimiters",
			input:   "label@joe@example.com",
			isValid: false,
		},
		{
			name:    "multiple ':' delimiters",
			input:   "label:super:secret:password@endpoint",
			isValid: false,
		},
		{
			name:    "non-alphanumeric characters in label",
			input:   "definitely$$not_*?a+Trello.Label@endpoint",
			isValid: false,
		},
		{
			name:    "non-alphanumeric characters in secret",
			input:   "label:-?_*@endpoint",
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
