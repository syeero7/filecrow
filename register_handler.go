package main

import (
	"io"
	"net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	name := r.URL.Query().Get("name")
	size := r.URL.Query().Get("size")
	if id == "" || name == "" || size == "" {
		http.Error(w, "transfer id or name or size missing", http.StatusBadRequest)
		return
	}

	pr, pw := io.Pipe()
	s := &Session{reader: pr, writer: pw, done: make(chan struct{})}
	ft := &FileTransfer{name: name, size: size, session: s}
	transfers.add(id, ft)

	// TODO: broadcast file info availability
	w.WriteHeader(http.StatusCreated)
}
