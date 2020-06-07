package habits

import (
	"testing"
)

func TestReadToken(t *testing.T) {
	tt := []struct {
		name     string
		filename string
		err      bool
	}{
		{
			name:     "invalid json file",
			filename: "../../config.example.yml",
			err:      true,
		},
		{
			name:     "non-existing token file",
			filename: "./no-way-this-file-exists.json",
			err:      true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := readToken(tc.filename)
			if (err != nil && !tc.err) || err == nil && tc.err {
				t.Fatalf("did not expect to get '%v' error", err)
			}
		})
	}
}
