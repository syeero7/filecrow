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
	conn, client, isClosed := newClient()

	fileServer.addClient(client)
	defer fileServer.removeClient(client)

	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		return err
	}

	if err := isClosed(c); err != nil {
		return err
	}
	defer conn.CloseNow()

	ctx := conn.CloseRead(context.Background())
	for {
		select {

		case msg := <-client.msgs:
			ctx2, cancel := context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			return conn.Write(ctx2, websocket.MessageText, msg)

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
