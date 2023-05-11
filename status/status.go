package status

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Client struct {
	etag   *string // etag to reduce API bandwidth usage
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		etag:   nil,
		client: http.DefaultClient,
	}
}

func (c *Client) Poll() (*SystemStatus, error) {
	resp, err := c.getData()
	if err != nil {
		return nil, err
	}
	switch resp.StatusCode {
	case http.StatusOK:
		newEtag := resp.Header.Get("etag")
		c.etag = &newEtag
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		result := new(SystemStatus)
		err = json.Unmarshal(body, result)
		if err != nil {
			return nil, err
		}
		return result, nil
	case http.StatusNotModified:
		// No updates, just return nil
		return nil, nil
	default:
		return nil, fmt.Errorf("unexpected http status code, expected 200 or 304, but got %d", resp.StatusCode)
	}
}

func (c *Client) getData() (*http.Response, error) {
	url := "https://www.githubstatus.com/api/v2/summary.json"
	reader := strings.Reader{}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, &reader)
	if err != nil {
		return nil, err
	}
	if c.etag != nil {
		req.Header.Add("If-None-Match", *c.etag)
	}
	req.Header.Add("User-Agent", fmt.Sprintf("gh-status/%s", strings.TrimLeft(Version, "v")))
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
