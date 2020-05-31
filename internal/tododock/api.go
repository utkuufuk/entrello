package tododock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// fetchTaskResponse represents the HTTP response body returned upon a successful
// GET request to the TodoDock task-fetch API endpoint
type fetchTasksResponse struct {
	Data struct {
		UserId int    `json:"user_id"`
		Tasks  []task `json:"tasks"`
	} `json:"data"`
}

// login logs-in to TodoDock with the configured user's credentials,
// and returns the user ID and JWT obtained from the HTTP response
func (t TodoDockSource) login() (id int, jwt string, err error) {
	req, err := json.Marshal(map[string]string{
		"email":    t.cfg.Email,
		"password": t.cfg.Password,
	})
	if err != nil {
		return -1, "", fmt.Errorf("could not build TodoDock login request body: %w", err)
	}

	url := fmt.Sprintf("%s/login", BASE_URL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(req))
	if err != nil {
		return -1, "", fmt.Errorf("could not login to TodoDock: %w", err)
	}
	defer resp.Body.Close()

	var data map[string]map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&data)
	return int(data["data"]["id"].(float64)), fmt.Sprintf("%s", data["data"]["token"]), nil
}

// fetchTasks retrieves all TodoDock tasks owned by the logged-in user with the given ID
func (t TodoDockSource) fetchTasks(id int, token string) (tasks []task, err error) {
	// build GET request with auth header
	url := fmt.Sprintf("%s/tasks/%d", BASE_URL, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return tasks, fmt.Errorf("could not create GET request to fetch TodoDock tasks: %w", err)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// make http request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return tasks, fmt.Errorf("could not fetch TodoDock tasks for user '%d': %w", id, err)
	}
	defer resp.Body.Close()

	// decode & return tasks
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tasks, fmt.Errorf("could not read response body: %w", err)
	}
	res := new(fetchTasksResponse)
	err = json.Unmarshal(body, &res)
	return res.Data.Tasks, nil
}
