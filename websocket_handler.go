package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/coder/websocket"
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	err := wsHelper(w, r)
	if errors.Is(err, context.Canceled) ||
		websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
		websocket.CloseStatus(err) == websocket.StatusGoingAway {
		return
	}

	if err != nil {
		log.Println(err)
	}
}

func wsHelper(w http.ResponseWriter, r *http.Request) error {
	nc := newClient()

	fileServer.addClient(nc.client)
	defer fileServer.removeClient(nc.client)

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"localhost:5173"}})
	if err != nil {
		return err
	}

	if err := nc.isClosed(c); err != nil {
		return err
	}

	defer nc.conn.CloseNow()
	ctx := nc.conn.CloseRead(context.Background())

	for {
		select {

		case msg := <-nc.client.msgs:
			ctx2, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := nc.conn.Write(ctx2, websocket.MessageText, msg); err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
