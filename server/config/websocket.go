package config

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	websocketWriteWait  = 10 * time.Second
	websocketPongWait   = 60 * time.Second
	websocketPingPeriod = 45 * time.Second
)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *Application) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	events := app.Hub.Subscribe(ctx)
	app.Hub.Publish("connection.opened", map[string]string{"transport": "websocket"})

	conn.SetReadLimit(512)
	conn.SetReadDeadline(time.Now().Add(websocketPongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(websocketPongWait))
		return nil
	})

	go func() {
		defer cancel()
		for {
			if _, _, err := conn.NextReader(); err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(websocketPingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			app.Hub.Publish("connection.closed", map[string]string{"transport": "websocket"})
			return
		case event, ok := <-events:
			if !ok {
				return
			}
			if err := writeWebSocketJSON(conn, event); err != nil {
				log.Printf("websocket write failed: %v", err)
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(websocketWriteWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func writeWebSocketJSON(conn *websocket.Conn, payload any) error {
	conn.SetWriteDeadline(time.Now().Add(websocketWriteWait))
	return conn.WriteJSON(payload)
}
