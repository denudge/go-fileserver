package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/denudge/go-fileserver/pkg/fileserver"
	"golang.org/x/sys/unix"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	timeout        = 10 * time.Second
	maxHeaderBytes = 4096
)

func main() {
	flag.IntVar(&port, "port", 8080, "the port to listen on")
	flag.StringVar(&rootDir, "root-dir", "/tmp", "the root directory")
	flag.IntVar(&bufferSize, "buffer-size", 8092, "the used buffer size")
	flag.BoolVar(&debug, "debug", false, "print more debug logs")
	flag.Parse()

	// check root rootDir existence
	stat, err := os.Stat(rootDir)
	if stat == nil || errors.Is(err, os.ErrNotExist) {
		log.Fatalf("error: root directory %q does not exist", rootDir)
		return
	}

	if !stat.IsDir() {
		log.Fatalf("error: root directory %q is no directory", rootDir)
		return
	}

	// TODO: does this only work for *NIX systems?
	if unix.Access(rootDir, unix.W_OK) != nil {
		log.Fatalf("error: root directory %q is not writeable", rootDir)
		return
	}

	urlFilePath := fileserver.FormatUrlPath("{folder}", "{filename}")

	mux := http.NewServeMux()
	mux.HandleFunc(http.MethodPost+" "+urlFilePath, checkFolderAndFile(uploadFile))
	mux.HandleFunc(http.MethodDelete+" "+urlFilePath, checkFolderAndFile(deleteFile))
	mux.HandleFunc(http.MethodGet+" "+urlFilePath, checkFolderAndFile(downloadFile))

	addr := fmt.Sprintf(":%d", port)
	server := http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       timeout,
		ReadHeaderTimeout: timeout,
		WriteTimeout:      timeout,
		IdleTimeout:       timeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}

	log.Printf("starting HTTP server on %s ...\n", addr)
	log.Printf("buffer size: %d, debug: %v\n", bufferSize, debug)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error running HTTP server: %s\n", err)
		}
	}
}
