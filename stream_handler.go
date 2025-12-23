package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ProgressWriter struct {
	writer     io.Writer
	total      int
	written    int
	onProgress func(int, int)
}

func (pw *ProgressWriter) Write(d []byte) (int, error) {
	n, err := pw.writer.Write(d)
	if n > 0 {
		pw.written += n
		pw.onProgress(pw.written, pw.total)
	}
	return n, err
}

type FileProgress struct {
	ID      string `json:"id"`
	Total   int    `json:"total"`
	Current int    `json:"current"`
	Type    string `json:"type"`
}

func streamHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	ft, exist := transfers.get(id)
	if !exist {
		http.Error(w, "transfer not found", http.StatusNotFound)
		return
	}

	pw := &ProgressWriter{
		total:  int(r.ContentLength),
		writer: ft.session.writer,
		onProgress: func(current, total int) {
			p := FileProgress{ID: id, Total: total, Current: current, Type: "progress"}
			msg, err := json.Marshal(p)
			if err != nil {
				log.Println(err)
				return
			}
			fileServer.broadcast(msg)
		},
	}

	ts := TransferState{ID: id, Type: "ready"}
	msg, err := json.Marshal(ts)
	if err != nil {
		log.Println(err)
		return
	}
	fileServer.broadcast(msg)

	if _, err := io.Copy(pw, r.Body); err != nil {
		log.Printf("transfer failed: %v", err)
	}

	ft.session.writer.Close()
	<-ft.session.done
	w.WriteHeader(http.StatusOK)
}
