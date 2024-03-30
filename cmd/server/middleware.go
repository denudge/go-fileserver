package main

import "net/http"

func checkFolderAndFile(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		next(w, r)
	}
}
