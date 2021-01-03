package checkpoint

import (
	"crypto/tls"
	"io"
	"net/http"
	"time"
)

func newHTTPClient(timeout time.Duration) *http.Client {
	tr := http.DefaultTransport.(*http.Transport)
	tr.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}
}

func sendHTTP(c *http.Client, url string, body io.Reader, header map[string]string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	for k, v := range header {
		req.Header.Add(k, v)
	}

	return c.Do(req)
}
