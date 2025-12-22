package main

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/coder/websocket"
	"golang.org/x/time/rate"
)

type FileServer struct {
	clients          map[*Client]struct{}
	mu               sync.Mutex
	broadcastLimiter rate.Limiter
}

type Client struct {
	msgs      chan []byte
	closeSlow func()
}

var fileServer = &FileServer{
	clients:          make(map[*Client]struct{}),
	broadcastLimiter: *rate.NewLimiter(rate.Every(time.Millisecond*100), 8),
}

func (fs *FileServer) broadcast(msg []byte) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.broadcastLimiter.Wait(context.Background())
	for client := range fs.clients {
		select {
		case client.msgs <- msg:
		default:
			go client.closeSlow()
		}
	}
}

func (fs *FileServer) addClient(c *Client) {
	fs.mu.Lock()
	fs.clients[c] = struct{}{}
	fs.mu.Unlock()
}

func (fs *FileServer) removeClient(c *Client) {
	fs.mu.Lock()
	delete(fs.clients, c)
	fs.mu.Unlock()
}

func newClient() (*websocket.Conn, *Client, func(*websocket.Conn) error) {
	const messageBuffer = 16
	var mu sync.Mutex
	var conn *websocket.Conn
	closed := false

	client := &Client{msgs: make(chan []byte, messageBuffer), closeSlow: func() {
		mu.Lock()
		defer mu.Unlock()
		closed = true
		if conn != nil {
			conn.Close(websocket.StatusPolicyViolation, "connection is too slow too keep up")
		}
	}}

	isClosed := func(c *websocket.Conn) error {
		mu.Lock()
		defer mu.Unlock()
		if closed {
			return net.ErrClosed
		}
		conn = c
		return nil
	}

	return conn, client, isClosed
}
