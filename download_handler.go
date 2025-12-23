package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type TransferState struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ft, exist := transfers.get(id)
	if !exist {
		http.Error(w, "transfer not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", ft.name))
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := io.Copy(w, ft.session.reader); err != nil {
		http.Error(w, fmt.Sprintf("transfer failed.\nerror: %v", err), http.StatusInternalServerError)
		return
	}

	ft.session.reader.Close()
	ft.session.done <- struct{}{}
	transfers.remove(id)

	ts := TransferState{ID: id, Type: "done"}
	msg, err := json.Marshal(ts)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to encode json.\nerror: %v", err), http.StatusInternalServerError)
		return
	}
	fileServer.broadcast(msg)

	w.WriteHeader(http.StatusOK)
}
