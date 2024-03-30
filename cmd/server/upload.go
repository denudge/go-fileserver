package main

import (
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	folder := r.PathValue("folder")
	filename := r.PathValue("filename")

	log.Printf("handling upload for folder %q, file %q", folder, filename)

	fullDir := rootDir + string(os.PathSeparator) + folder
	if _, err := os.Stat(fullDir); errors.Is(err, os.ErrNotExist) {
		log.Printf("creating folder %q", folder)
		err := os.Mkdir(fullDir, os.FileMode(0700))
		if err != nil {
			log.Printf("could not create folder %q: %v", folder, err)
			w.WriteHeader(500)
			_, _ = w.Write([]byte("Internal Server Error\n"))
			return
		}
	}

	fullName := rootDir + string(os.PathSeparator) + folder + string(os.PathSeparator) + filename

	// check if file already exists
	stat, err := os.Stat(fullName)
	if stat != nil || err == nil || !errors.Is(err, os.ErrNotExist) {
		log.Printf("error handling upload for folder %q, file %q: file already exists", folder, filename)
		w.WriteHeader(400)
		_, _ = w.Write([]byte("error: file already exists\n"))
		return
	}

	// create file
	file, err := os.Create(fullName)
	if err != nil {
		log.Printf("error handling upload: could not open file %q: %v", fullName, err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Internal Server Error\n"))
		return
	}

	// create buffer
	buffer := make([]byte, bufferSize)

	// write chunk by chunk to file
	hasData := true
	for {
		n, err := r.Body.Read(buffer)
		if errors.Is(err, io.EOF) {
			hasData = false
		} else if err != nil {
			log.Printf("error handling upload for folder %q, file %q: could not read from request body: %v", folder, filename, err)
			w.WriteHeader(500)
			_, _ = w.Write([]byte("Internal Server Error\n"))
			return
		}

		if n > 0 {
			if debug {
				log.Printf("received %d bytes from client", n)
			}

			m, err := file.Write(buffer[:n])
			if err != nil {
				log.Printf("error handling upload for folder %q, file %q: could not write to file: %v", folder, filename, err)
				w.WriteHeader(500)
				_, _ = w.Write([]byte("error processing the upload\n"))
				return
			}

			if debug {
				log.Printf("wrote %d bytes to disk", m)
			}

			if m != n {
				log.Printf("error handling upload for folder %q, file %q: expected %d bytes to write, wrote %d bytes.", folder, filename, n, m)
				w.WriteHeader(500)
				_, _ = w.Write([]byte("error processing the upload\n"))
				return
			}
		}

		if !hasData {
			if debug {
				log.Println("end of file!")
			}
			break
		}
	}

	if err := file.Close(); err != nil {
		log.Printf("error handling upload for folder %q, file %q: could not close file: %v", folder, filename, err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte("error processing the upload\n"))
		return
	}

	log.Printf("successfully handled upload for folder %q, file %q", folder, filename)
	w.WriteHeader(201)
	_, _ = w.Write([]byte("file successfully uploaded\n"))
}
