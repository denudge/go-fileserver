package fileserver

import (
	"fmt"
	"io"
	"net/http"
)

func (c *Client) DeleteFile(folder, filename string) error {
	resp, err := c.Do(http.MethodDelete, folder, filename, nil)
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
