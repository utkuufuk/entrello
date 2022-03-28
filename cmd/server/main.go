package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/logger"
	"github.com/utkuufuk/entrello/internal/service"
	"github.com/utkuufuk/entrello/pkg/trello"
)

func main() {
	http.HandleFunc("/", controller)
	http.ListenAndServe(fmt.Sprintf(":%s", config.ServerCfg.Port), nil)
}

func controller(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodHead {
		w.WriteHeader(http.StatusOK)
		return
	}

	if req.Method == http.MethodGet {
		handlePollRequest(w, req)
		return
	}

	if req.Method == http.MethodPost {
		handleTrelloWebhookRequest(w, req)
		return
	}

	logger.Warn("Method %s not allowed", req.Method)
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func handlePollRequest(w http.ResponseWriter, req *http.Request) {
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
	hash := req.Header.Get("x-trello-webhook")
	if hash == "" {
		logger.Warn("Missing 'x-trello-webhook' header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		logger.Error("Could not read request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	logger.Info("Callback URL:", config.ServerCfg.TrelloWebhookCallbackUrl)
	logger.Info("Hash:", hash)
	logger.Info("Body:", body)
	if !trello.VerifyTrelloSignature(
		config.ServerCfg.TrelloWebhookCallbackUrl,
		config.ServerCfg.TrelloSecret,
		hash,
		body,
	) {
		logger.Warn("Invalid Trello webhook signature")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
