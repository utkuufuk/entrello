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

type service *sheets.SpreadsheetsValuesService

// createService creates a spreadsheet service sets it in the source
func (s *source) createService(ctx context.Context) error {
	creds, err := ioutil.ReadFile(s.credentialsFile)
	if err != nil {
		return fmt.Errorf("could not read client credentials file: %w", err)
	}

	cfg, err := google.ConfigFromJSON(creds, SCOPE)
	if err != nil {
		return fmt.Errorf("could not parse client secret file: %w", err)
	}

	token, err := readToken(s.tokenFile)
	if err != nil {
		return fmt.Errorf("could not find auth token: %w", err)
	}

	client := cfg.Client(ctx, token)
	service, err := sheets.New(client)
	if err != nil {
		return fmt.Errorf("could not create google spreadsheets service: %w", err)
	}

	s.service = service.Spreadsheets.Values
	return nil
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

// read reads a range of cell values from the spreadsheet
func (s source) readCells(spreadsheetId string, rangeName string) ([][]interface{}, error) {
	resp, err := (*sheets.SpreadsheetsValuesService)(s.service).Get(spreadsheetId, rangeName).Do()
	if err != nil {
		return nil, fmt.Errorf("could not read cells: %w", err)
	}
	return resp.Values, nil
}
