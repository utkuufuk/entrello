package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/logger"
	"github.com/utkuufuk/entrello/internal/service"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil)
}

func handler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	user, pwd, ok := req.BasicAuth()
	if !ok {
		fmt.Println("Error parsing basic auth")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if user != os.Getenv("USERNAME") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if pwd != os.Getenv("PASSWORD") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var cfg config.Config
	if err = json.Unmarshal(body, &cfg); err != nil {
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
