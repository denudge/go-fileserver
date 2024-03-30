package fileserver

import "fmt"

const (
	UrlPathPrefix = "/v1/fileserver"
	UserAgent     = "go-fileserver/1.0"
)

var (
	ServerAddress = "http://localhost:8080"
)

func FormatUrlPath(folder, filename string) string {
	return fmt.Sprintf("%s/%s/%s", UrlPathPrefix, folder, filename)
}
