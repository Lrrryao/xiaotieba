package api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	connections = make(map[*websocket.Conn]bool)
	broadcast   = make(chan []byte)
)
