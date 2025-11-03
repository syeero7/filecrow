package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func deleteFileHandler(fsvr *fileServer, w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("file")
	if strings.ContainsRune(filename, filepath.Separator) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := os.Remove(filepath.Join(fsvr.directory, filename)); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
