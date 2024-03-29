package main

import (
	"errors"
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net/http"
	"os"
	"time"
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

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/fileserver/{folder}/{filename}", uploadFile)
	mux.HandleFunc("DELETE /v1/fileserver/{folder}/{filename}", deleteFile)
	mux.HandleFunc("GET /v1/fileserver/{folder}/{filename}", DownloadFile)

	addr := fmt.Sprintf(":%d", port)
	server := http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       10 * time.Second,
		MaxHeaderBytes:    4096,
	}

	log.Printf("starting HTTP server on %s ...\n", addr)
	log.Printf("buffer size: %d, debug: %v\n", bufferSize, debug)
	if err := server.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			log.Printf("error running HTTP server: %s\n", err)
		}
	}
}
