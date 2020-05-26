package tododock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/utkuufuk/entrello/internal/config"
	"github.com/utkuufuk/entrello/internal/trello"
)

const (
	BASE_URL = "https://tododock.com/api"
)

type TodoDockSource struct {
	cfg config.TodoDock
}

// task represents the TodoDock task model
type task struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	State         string `json:"state"`
	Color         string `json:"color"`
	Notes         string `json:"notes"`
	Period        int    `json:"period"`
	NextResetDate string `json:"next_reset_date"`
	MuteEmails    int    `json:"mute_reminder_emails"`
}

// fetchTaskResponse represents the HTTP response body returned upon a successful
// GET request to the TodoDock task-fetch API endpoint
type fetchTasksResponse struct {
	Data struct {
		UserId int    `json:"user_id"`
		Tasks  []task `json:"tasks"`
	} `json:"data"`
}

func GetSource(cfg config.TodoDock) TodoDockSource {
	return TodoDockSource{cfg}
}

func (t TodoDockSource) IsEnabled() bool {
	return t.cfg.Enabled
}

func (t TodoDockSource) IsStrict() bool {
	return t.cfg.Strict
}

func (t TodoDockSource) GetName() string {
	return "TodoDock"
}

func (t TodoDockSource) GetLabel() string {
	return t.cfg.Label
}

func (t TodoDockSource) GetPeriod() config.Period {
	return t.cfg.Period
}

func (t TodoDockSource) FetchNewCards() (cards []trello.Card, err error) {
	id, token, err := t.login()
	if err != nil {
		return cards, nil
	}
	tasks, err := t.fetchTasks(id, token)
	return toCards(tasks, t.cfg.Label)
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

// toCards cherry-picks the 'active' and 'due' tasks from a list of tasks,
// then returns a list of cards containing those
func toCards(tasks []task, label string) (cards []trello.Card, err error) {
	cards = make([]trello.Card, 0, len(tasks))
	soon := time.Now().AddDate(0, 0, 2)
	for _, t := range tasks {
		d, ok, err := shouldCreateCard(t, soon)
		if !ok {
			if err != nil {
				return cards, err
			}
			continue
		}

		url := fmt.Sprintf("https://tododock.com/home/%d\n%s", t.Id, t.Notes)
		c, err := trello.NewCard(t.Name, label, url, &d)
		if err != nil {
			return cards, fmt.Errorf("could not create card: %w", err)
		}
		cards = append(cards, c)
	}
	return cards, nil
}

// shouldCreateCard decides if a Trello card should be created from the given TodoDock task
// by looking at the 'status' & 'next reset date' attributes of the task
func shouldCreateCard(t task, ref time.Time) (d time.Time, ok bool, err error) {
	d, err = time.Parse("2006-01-02 15:04:05", t.NextResetDate)
	if err != nil {
		return ref, false, fmt.Errorf("could not parse next reset date '%s': %w", t.NextResetDate, err)
	}

	// only create a card for active taks that are due soon
	if t.State != "active" || d.After(ref) {
		return ref, false, nil
	}
	return d, true, nil
}
