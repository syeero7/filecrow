package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

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
		log.Printf("transfer failed: %v", err)
	}

	// TODO: broadcast file transfer speed

	ft.session.reader.Close()
	ft.session.done <- struct{}{}

	transfers.remove(id)
	w.WriteHeader(http.StatusOK)
}
