package main

import (
	"io"
	"net/http"
)

func testUpload(w http.ResponseWriter, r *http.Request) {
	if r.Body != nil {
		defer r.Body.Close()
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Serer Error"))
	}

	w.WriteHeader(200)
	_, _ = w.Write(content)
}
