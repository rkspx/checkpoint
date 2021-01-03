package checkpoint

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// Session is representation of user session.
type Session struct {
	ID        string `json:"id"`
	server    string
	client    *http.Client
	Published bool
}

// Object represents Check Point Object.
type Object struct {
	Name string `json:"name"`
	UID  string `json:"uid"`
	Type string `json:"type"`
}

// Task represents Check Point Task.
type TaskResult struct {
	ID string `json:"task-id"`
}

const (
	taskStatusInProgress = "in progress"
)

type Task struct {
	Name       string `json:"task-name"`
	ID         string `json:"task-id"`
	Progress   int    `json:"progress-percentage"`
	Status     string `json:"status"`
	Suppressed bool   `json:"suppressed"`
}

func (t *Task) IsDone() bool {
	return t.Status != taskStatusInProgress
}

type Tasks struct {
	Tasks []*Task `json:"tasks"`
}
type DiscardResult struct {
	Message string `json:"message"`
	N       int    `json:"number-of-discarded-changes"`
}

type ExitResult struct {
	Message string `json:"message"`
}

// NewSession login to Web Management API Server, and create a new session.
func (c *Client) NewSession(username, password, APIKey string) (*Session, error) {
	if APIKey != "" {
		return nil, fmt.Errorf("API Key login is not supported yet")
	}

	url := c.APIURL + "login"

	body, err := createRequestBody(map[string]string{
		"user":     username,
		"password": password,
	})

	if err != nil {
		return nil, err
	}

	req, err := makeRequest(url, body, "")
	if err != nil {
		return nil, err
	}

	res, err := c.Execute(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, parseErrorResponse(res.Body)
	}

	var out LoginResponse
	if err := json.NewDecoder(res.Body).Decode(&out); err != nil {
		return nil, err
	}

	return &Session{
		ID:     out.SID,
		server: c.APIURL,
		client: c.client,
	}, nil
}

func (c *Session) Do(action string, payload map[string]interface{}, result interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	body := bytes.NewReader(b)
	url := fmt.Sprintf("%s/%s", c.server, action)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-chkp-sid", c.ID)

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// TODO: add retry here

	if res.StatusCode != http.StatusOK {
		r := &FailResponse{}
		r.fromReader(res.Body)

		err := parseError(*r)
		if err != nil {
			return err
		}

		return err
	}

	be, err := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(be, result); err != nil {
		return err
	}

	return nil
}

func (c *Session) Publish() error {
	task := new(TaskResult)
	if err := c.Do("publish", map[string]interface{}{}, task); err != nil {
		return err
	}

	if err := c.waitTask(task); err != nil {
		return err
	}

	c.Published = true

	return nil
}

func (c *Session) Discard() (*DiscardResult, error) {
	res := new(DiscardResult)
	if err := c.Do("discard", map[string]interface{}{}, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Session) Exit() (*ExitResult, error) {
	e := new(ExitResult)
	if err := c.Do("logout", map[string]interface{}{}, e); err != nil {
		return nil, err
	}

	return e, nil
}

func (c *Session) Wait() error {
	t, err := randInt64(10)
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(t) * time.Second)
	return nil
}

func (c *Session) waitTask(t *TaskResult) error {
	sleepTime := 500 * time.Millisecond
	tasks, err := c.showTask(t.ID)
	if err != nil {
		return err
	}

	task := tasks.Tasks[0]

	timeout := time.NewTimer(10 * time.Minute)

	for !task.IsDone() {
		tasks, err := c.showTask(t.ID)
		if err != nil {
			return err
		}

		task = tasks.Tasks[0]
		time.Sleep(sleepTime)

		select {
		case <-timeout.C:
			return fmt.Errorf("task timeout exceeded")
		default:
		}
	}

	return nil
}

func (c *Session) showTask(id string) (*Tasks, error) {
	tasks := new(Tasks)
	if err := c.Do("show-task", map[string]interface{}{
		"task-id": id,
	}, tasks); err != nil {
		return nil, err
	}

	if len(tasks.Tasks) == 0 {
		return nil, fmt.Errorf("no task")
	}

	return tasks, nil
}

func createAPIURL(slug string) string {
	return ""
}
