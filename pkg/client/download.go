package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func (c *Client) DownloadFile(folder, filename, localFile string) error {
	fullPath, err := filepath.Abs(localFile)
	if err != nil {
		return err
	}

	// check if file already exists
	stat, err := os.Stat(fullPath)
	if stat != nil || err == nil || !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error: file %q already exists", fullPath)
	}

	file, err := os.Create(fullPath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error: could not create local file %q", fullPath)
	}
	defer file.Close()

	path := fmt.Sprintf("/v1/fileserver/%s/%s", url.PathEscape(folder), url.PathEscape(filename))
	req, err := http.NewRequest(http.MethodGet, ServerAddress+path, nil)
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

	if resp.StatusCode >= 400 && resp.StatusCode < 500 {
		message, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("could not parse server response: %w", err)
		}

		return fmt.Errorf(string(message))
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("could not stream data from server to disk: %w", err)
	}

	return nil
}
