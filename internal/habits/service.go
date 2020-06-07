package habits

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const (
	SCOPE = "https://www.googleapis.com/auth/spreadsheets"
)

// initializeService creates and sets a spreadsheet service within the source struct
func (s *source) initializeService(ctx context.Context) error {
	cfg, token, err := readCreds(s.credentialsFile, s.tokenFile)
	if err != nil {
		return fmt.Errorf("failed to get credentials for google spreadsheets: %w", err)
	}

	client := cfg.Client(ctx, token)
	service, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("could not create google spreadsheets client: %w", err)
	}

	s.service = service.Spreadsheets.Values
	return nil
}

// readCreds reads and returns credentials from the configured files
func readCreds(credentialsFile, tokenFile string) (*oauth2.Config, *oauth2.Token, error) {
	c, err := ioutil.ReadFile(credentialsFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not read client credentials file: %w", err)
	}

	cfg, err := google.ConfigFromJSON(c, SCOPE)
	if err != nil {
		return nil, nil, fmt.Errorf("could not parse client secret file: %w", err)
	}

	token, err := readToken(tokenFile)
	if err != nil {
		return nil, nil, fmt.Errorf("could not find auth token: %w", err)
	}
	return cfg, token, nil
}

// readToken reads the client auth token from a JSON file
func readToken(tokenPath string) (*oauth2.Token, error) {
	f, err := os.Open(tokenPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}
