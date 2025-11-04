package main

import (
	"net/http"
	"os"
	"path"
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

func deleteAllHandler(fsvr *fileServer, w http.ResponseWriter, r *http.Request) {
	if err := removeContent(fsvr.directory); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func removeContent(dir string) error {
	files, err := filepath.Glob(path.Join(dir, "*"))
	if err != nil {
		return err
	}

	for _, f := range files {
		if err := os.RemoveAll(f); err != nil {
			return err
		}
	}

	return nil
}
