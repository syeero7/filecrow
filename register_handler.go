package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type RegisterFile struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
	Size int    `json:"size"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var data RegisterFile

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, "failed to decode json", http.StatusBadRequest)
		log.Println(err)
		return
	}

	if data.ID == "" || data.Name == "" || data.Size == 0 {
		http.Error(w, "transfer id or name or size missing", http.StatusBadRequest)
		return
	}

	pr, pw := io.Pipe()
	s := &Session{reader: pr, writer: pw, done: make(chan struct{})}
	ft := &FileTransfer{name: data.Name, size: data.Size, session: s}
	transfers.add(data.ID, ft)

	fileServer.broadcast(body)
	w.WriteHeader(http.StatusCreated)
}
