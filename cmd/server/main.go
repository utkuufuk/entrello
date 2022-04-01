package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/exp/slices"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/logger"
	"github.com/utkuufuk/entrello/internal/service"
	"github.com/utkuufuk/entrello/pkg/trello"
)

var client trello.Client

func main() {
	client = trello.NewClient(config.Trello{
		ApiKey:   config.ServerCfg.TrelloApiKey,
		ApiToken: config.ServerCfg.TrelloApiToken,
		BoardId:  config.ServerCfg.TrelloBoardId,
	})

	http.HandleFunc("/", handlePollRequest)
	http.HandleFunc("/trello-webhook", handleTrelloWebhookRequest)
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerCfg.Port), nil)
}

func handlePollRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		logger.Warn("Method %s not allowed for %s", req.Method, req.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user, pwd, ok := req.BasicAuth()
	if !ok {
		logger.Warn("Could not parse basic auth.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if user != config.ServerCfg.Username {
		logger.Warn("Invalid user name: %s", user)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if pwd != config.ServerCfg.Password {
		logger.Warn("Invalid password: %s", pwd)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("Could not read request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var cfg config.RunnerConfig
	if err = json.Unmarshal(body, &cfg); err != nil {
		logger.Warn("Invalid request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = service.Poll(cfg); err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func handleTrelloWebhookRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method != http.MethodPost {
		logger.Warn("Method %s not allowed for %s", req.Method, req.URL.Path)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	hash := req.Header.Get("x-trello-webhook")
	if hash == "" {
		logger.Warn("Missing 'x-trello-webhook' header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("Could not read Trello webhook request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !trello.VerifyWebhookSignature(
		config.ServerCfg.TrelloWebhookCallbackUrl,
		config.ServerCfg.TrelloSecret,
		hash,
		body,
	) {
		logger.Warn("Invalid Trello webhook signature")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var wrb trello.WebhookRequestBody
	if err = json.Unmarshal(body, &wrb); err != nil {
		logger.Warn("Invalid Trello webhook request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	archivedCardId := trello.ParseArchivedCardId(wrb)
	if archivedCardId == "" {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	archivedCard, err := client.GetCard(archivedCardId)
	if err != nil {
		logger.Error("Could not fetch archived Trello card: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Info("Archived card name: %v", archivedCard.Name)
	logger.Info("Archived card description: %v", archivedCard.Desc)
	logger.Info("Archived card labels: %v", archivedCard.Labels)

	labelIds := make([]string, 0)
	for _, label := range archivedCard.Labels {
		labelIds = append(labelIds, label.ID)
	}

	for _, service := range config.ServerCfg.Services {
		if slices.Contains(labelIds, service.Label) {
			logger.Info("Archived card matches service %v!!", service)
		}
	}

	w.WriteHeader(http.StatusOK)
}
