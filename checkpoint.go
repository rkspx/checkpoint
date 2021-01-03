// Package checkpoint contains all required codes for checkpoint actions.
package checkpoint

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// DefaultTimeout is default timeout for every request.
	DefaultTimeout = 120 * time.Second
)

var (
	// ErrNoSessionID is error thrown when there's no session ID stored.
	ErrNoSessionID = errors.New("no session ID")
)

func env(name, fallback string) string {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}

	return value
}

func makeRequest(url string, body io.Reader, sessionID string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	if sessionID != "" {
		req.Header.Add("X-chkp-sid", sessionID)
	}

	return req, nil
}

func parseErrorResponse(body io.Reader) error {
	var out FailResponse
	if err := json.NewDecoder(body).Decode(&out); err != nil {
		return fmt.Errorf("can not parse error: %s", err.Error())
	}

	return parseError(out)
}

func parseObjectResponse(body io.Reader) (*Object, error) {
	var o Object
	if err := json.NewDecoder(body).Decode(&o); err != nil && err != io.EOF {
		return nil, fmt.Errorf("can not parse object response: %s", err.Error())
	}

	return &o, nil

}

func createRequestBody(body interface{}) (io.Reader, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
