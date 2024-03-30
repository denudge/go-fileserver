package main

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func DownloadFile(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	folder := r.PathValue("folder")
	filename := r.PathValue("filename")

	log.Printf("handling download for folder %q, file %q", folder, filename)

	fullName := rootDir + string(os.PathSeparator) + folder + string(os.PathSeparator) + filename
	// check if file exists
	_, err := os.Stat(fullName)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("error handling delete for folder %q, file %q: file does not exist", folder, filename)
		w.WriteHeader(404)
		_, _ = w.Write([]byte("error: file does not exist\n"))
		return
	}

	// read file
	file, err := os.Open(fullName)
	if err != nil {
		log.Printf("error handling download: could not open file %q: %v", fullName, err)
		w.WriteHeader(500)
		// _, _ = w.Write([]byte("Internal Server Error\n"))
		return
	}

	defer file.Close()

	bufReader := bufio.NewReader(file)

	n, err := io.Copy(w, bufReader)
	if err != nil {
		log.Printf("error download upload for folder %q, file %q: could not send data to client: %v", folder, filename, err)
		w.WriteHeader(500)
		// write no message to client because the body is used for the content
		return
	}

	if debug {
		log.Printf("sent %d bytes to client", n)
	}

	log.Printf("successfully handled download for folder %q, file %q", folder, filename)
}
