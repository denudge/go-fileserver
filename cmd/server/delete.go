package main

import (
	"errors"
	"log"
	"net/http"
	"os"
)

func deleteFile(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	folder := r.PathValue("folder")
	filename := r.PathValue("filename")

	if !folderRegexp.MatchString(folder) {
		w.WriteHeader(422)
		_, _ = w.Write([]byte("error: invalid folder name\n"))
		return
	}

	if !fileRegexp.MatchString(filename) {
		w.WriteHeader(422)
		_, _ = w.Write([]byte("error: invalid file name\n"))
		return
	}

	log.Printf("handling delete for folder %q, file %q", folder, filename)

	fullName := rootDir + string(os.PathSeparator) + folder + string(os.PathSeparator) + filename

	// check if file exists
	_, err := os.Stat(fullName)
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("error handling delete for folder %q, file %q: file does not exist", folder, filename)
		w.WriteHeader(404)
		_, _ = w.Write([]byte("error: file does not exist\n"))
		return
	}

	// delete file
	err = os.Remove(fullName)
	if err != nil {
		log.Printf("error: could not delete file %q: %v", fullName, err)
		w.WriteHeader(500)
		_, _ = w.Write([]byte("Internal Server Error\n"))
		return
	}

	log.Printf("successfully handled delete for folder %q, file %q", folder, filename)
	w.WriteHeader(200)
	_, err = w.Write([]byte("file successfully deleted\n"))
	if err != nil {
		log.Printf("error: could not write delete response: %v", err)
	}

	return
}
