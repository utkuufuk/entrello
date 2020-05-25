package config

import (
	"testing"
)

func TestGetEnabledSourcesAndLabels(t *testing.T) {
	tt := []struct {
		name     string
		filename string
		err      bool
	}{
		{
			name:     "valid file",
			filename: "../../config.example.yml",
			err:      false,
		},
		{
			name:     "non-existing file",
			filename: "../../config.example1.yml",
			err:      true,
		},
		{
			name:     "illegal YAML syntax",
			filename: "../../README.md",
			err:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ReadConfig(tc.filename)

			if (err != nil && !tc.err) || err == nil && tc.err {
				t.Fatalf("did not expect to get '%v' error", err)
			}
		})
	}
}
