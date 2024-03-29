package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func (c *Client) DeleteFile(folder, filename string) error {
	path := fmt.Sprintf("/v1/fileserver/%s/%s", url.PathEscape(folder), url.PathEscape(filename))
	req, err := http.NewRequest(http.MethodDelete, ServerAddress+path, nil)
	if err != nil {
		return fmt.Errorf("could not call fileserver: %w", err)
	}
	req.Header.Set("User-Agent", "go-fileserver/1.0")

	resp, err := c.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return err
	}

	message, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse server response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf(string(message))
	}

	return nil
}
