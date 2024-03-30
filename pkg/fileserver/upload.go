package fileserver

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
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

	resp, err := c.Do(http.MethodPost, folder, filename, bufio.NewReader(file))
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
