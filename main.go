package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed web
var frontend embed.FS

func main() {
	mux := http.NewServeMux()
	fsvr := fileServer{}
	if err := fsvr.makeFSDir(); err != nil {
		log.Fatal(err)
	}

	mux.Handle("GET /web/", http.FileServer(http.FS(frontend)))
	mux.HandleFunc("POST /upload", fsvr.middleware(uploadHandler))
	mux.HandleFunc("GET /files/{file}", fsvr.middleware(downloadHandler))
	mux.HandleFunc("POST /delete/{file}", fsvr.middleware(deleteFileHandler))
	mux.HandleFunc("GET /", fsvr.middleware(fileHandler))

	server := &http.Server{Addr: fmt.Sprintf(":%s", "8090"), Handler: mux}
	serverErr := make(chan error, 1)

	go func() {
		log.Printf("file server is running on http://localhost%s\n", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			serverErr <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		log.Printf("server error: %v\n", err)
	case <-stop:
	}

	log.Println("server is shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v\n", err)
	}
}
