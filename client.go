package checkpoint

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Client core logic of Check Point API client.
type Client struct {
	server string
	port   string
	client *http.Client

	sessionID string

	// APIURL is url for Check Point Web API.
	APIURL string
}

// New creates a new configurator client.
func New(server string, port string) *Client {

	c := &Client{
		server: server,
		port:   port,
		client: newHTTPClient(DefaultTimeout),
	}

	c.APIURL = fmt.Sprintf("https://%s:%s/web_api/", c.server, c.port)
	return c
}

// HasSessionID returns true if there is session ID stored.
func (c *Client) HasSessionID() bool {
	return c.sessionID != ""
}

// SessionID returns current stored session ID.
func (c *Client) SessionID() string {
	return c.sessionID
}

// SetSessionID update session id information.
func (c *Client) SetSessionID(id string) {
	c.sessionID = id
}

// ClearSessionID clear any stored session ID.
func (c *Client) ClearSessionID() {
	if c.sessionID != "" {
		c.sessionID = ""
	}
}

// Execute execute generated request.
func (c *Client) Execute(req *http.Request) (*http.Response, error) {
	return c.client.Do(req)
}

func (c *Client) call(slug string, body interface{}, checkSession bool) (io.Reader, error) {
	b, err := createRequestBody(body)
	if err != nil {
		return nil, err
	}

	if checkSession {
		if !c.HasSessionID() {
			return nil, ErrNoSessionID
		}
	}

	req, err := makeRequest(slug, b, c.SessionID())
	if err != nil {
		return nil, err
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	// fix keep-alive bug
	_, err = ioutil.ReadAll(res.Body)

	if res.StatusCode != http.StatusOK {
		err = parseErrorResponse(res.Body)
		return nil, err
	}

	return res.Body, nil
}

// CallObject do the API call that returns back object response.
func (c *Client) CallObject(slug string, body interface{}) (*Object, error) {
	res, err := c.call(slug, body, true)
	if err != nil {
		return nil, err
	}

	return parseObjectResponse(res)
}

// CallTask do the API call that returns back task.
func (c *Client) CallTask(slug string, body interface{}) error {
	_, err := c.call(slug, body, true)
	if err != nil {
		return err
	}

	// TODO: parse task object here

	return nil
}
