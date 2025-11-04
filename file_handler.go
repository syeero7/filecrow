package main

import (
	"html/template"
	"log"
	"net/http"
)

type FSData struct {
	Files []File
	Usage *DiskUsage
}

func fileHandler(fsvr *fileServer, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	t, err := template.ParseFS(frontend, "web/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if err := fsvr.readFSDir(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	data := FSData{Files: fsvr.Files, Usage: getDiskUsage()}
	if err := t.Execute(w, data); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
}
