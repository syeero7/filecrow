package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//go:embed dist
var frontend embed.FS

func main() {
	port := flag.Int("port", 8080, "The port number to run the server on.")
	flag.Parse()

	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		flag.PrintDefaults()
		os.Exit(0)
	}

	distFS, err := fs.Sub(frontend, "dist")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /register", registerHandler)
	mux.HandleFunc("POST /stream", streamHandler)
	mux.HandleFunc("GET /download", downloadHandler)
	mux.HandleFunc("/ws", websocketHandler)
	mux.Handle("/", http.FileServer(http.FS(distFS)))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", *port),
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 5,
	}
	serverErr := make(chan error, 1)

	go func() {
		fmt.Printf("access main interface at http://localhost%s\n", server.Addr)
		printWebInterfaceAddr(server.Addr)
		fmt.Println("---------------------------------------------------")
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v\n", err)
	}
}

func printWebInterfaceAddr(port string) {
	conn, err := net.Dial("udp", "1.1.1.1:80")
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()
	if addr, ok := conn.LocalAddr().(*net.UDPAddr); ok {
		fmt.Printf("access other interfaces at http://%s%s\n", addr.IP.To4().String(), port)
	}
}
