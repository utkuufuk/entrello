package tododock

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/utkuufuk/entrello/internal/trello"
)

type task struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	State         string      `json:"state"`
	Color         string      `json:"color"`
	Notes         interface{} `json:"notes"`
	Period        int         `json:"period"`
	NextResetDate string      `json:"next_reset_date"`
	MuteEmails    int         `json:"mute_reminder_emails"`
}

type fetchTasksResponse struct {
	Data struct {
		UserID int    `json:"user_id"`
		Tasks  []task `json:"tasks"`
	} `json:"data"`
}

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

func toCards(tasks []task) (c []trello.Card, err error) {
	c = make([]trello.Card, 0, len(tasks))
	for _, t := range tasks {
		d, err := time.Parse("2006-01-06 15:04:05", t.NextResetDate)
		if err != nil {
			return c, fmt.Errorf("could not parse next reset date '%s': %w", t.NextResetDate, err)
		}

		c = append(c, trello.Card{
			Name:    t.Name,
			Label:   "TodoDock",
			List:    "To-Do",
			DueDate: d,
		})
	}
	return c, nil
}
