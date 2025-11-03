package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
)

type File struct{ Name, Size string }

type fileServer struct {
	Files     []File
	directory string
}

func (f *fileServer) middleware(fn func([]File, http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f.readFSDir(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
			return
		}

		fn(f.Files, w, r)
	}
}

func (f *fileServer) makeFSDir() error {
	dir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	perm := 0o755
	if runtime.GOOS == "windows" {
		perm = 0o777
	}

	fsdir := path.Join(dir, "filecrow/files")
	if err := os.MkdirAll(fsdir, os.FileMode(perm)); err != nil {
		return err
	}
	f.directory = fsdir
	return nil
}

func (f *fileServer) readFSDir() error {
	if len(f.Files) > 0 {
		f.Files = []File{}
	}

	entries, err := os.ReadDir(f.directory)
	if err != nil {
		log.Println(f.directory)

		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return err
		}

		file := File{
			Name: info.Name(),
			Size: humanReadSize(info.Size()),
		}
		f.Files = append(f.Files, file)
	}

	return nil
}

func humanReadSize(s int64) string {
	const unit = 1000
	if s < unit {
		return fmt.Sprintf("%d B", s)
	}

	div, exp := int64(unit), 0
	for n := s / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}

	return fmt.Sprintf("%.2f %cB", float64(s)/float64(div), "kMGTPE"[exp])
}
