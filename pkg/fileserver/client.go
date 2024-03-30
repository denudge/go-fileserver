package fileserver

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	client *http.Client
}

func NewClient() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Do(method, folder, filename string, body io.Reader) (*http.Response, error) {
	addr := ServerAddress + FormatUrlPath(url.PathEscape(folder), url.PathEscape(filename))
	req, err := http.NewRequest(method, addr, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	return c.client.Do(req)
}
