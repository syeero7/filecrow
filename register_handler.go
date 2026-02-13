package main

import (
	"encoding/json"
	"fmt"
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
		http.Error(w, fmt.Sprintf("failed to read request body.\nerror: %v", err), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(body, &data); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode json.\nerror: %v", err), http.StatusBadRequest)
		return
	}

	if data.Name == "" || data.Size == 0 {
		http.Error(w, "transfer name or size missing", http.StatusBadRequest)
		return
	}

	id, err := generateUUID()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data.ID = id
	pr, pw := io.Pipe()
	s := &Session{reader: pr, writer: pw, done: make(chan struct{})}
	ft := &FileTransfer{name: data.Name, session: s}
	transfers.add(data.ID, ft)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": data.ID})

	msg, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to encode json: %v", err)
		return
	}

	fileServer.broadcast(msg)
}
