package client

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

func (c *Client) UploadFile(folder, filename, localFile string) error {
	fullPath, err := filepath.Abs(localFile)
	if err != nil {
		return err
	}

	// check if file exists
	_, err = os.Stat(fullPath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error: file %q does not exist", fullPath)
	}

	file, err := os.Open(fullPath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("error: could not open file %q", fullPath)
	}

	defer file.Close()

	bufReader := bufio.NewReader(file)
	path := fmt.Sprintf("/v1/fileserver/%s/%s", url.PathEscape(folder), url.PathEscape(filename))
	req, err := http.NewRequest(http.MethodPost, ServerAddress+path, bufReader)
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
