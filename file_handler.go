package main

import (
	"html/template"
	"log"
	"net/http"
)

type FSData struct {
	Files []File
}

func fileHandler(files []File, w http.ResponseWriter, _ *http.Request) {
	t, err := template.ParseFS(frontend, "web/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := FSData{Files: files}
	if err := t.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
