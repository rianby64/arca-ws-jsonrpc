package arca

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func init() {
	upgradeConnection = func(
		w http.ResponseWriter,
		r *http.Request,
	) (*websocket.Conn, error) {
		return upgrader.Upgrade(w, r, nil) // HARD-DEPENDENCY
	}
	readJSON = func(
		conn *websocket.Conn,
		request *JSONRPCrequest,
	) error {
		return conn.ReadJSON(&request) // HARD-DEPENDENCY
	}
	writeJSON = func(
		conn *websocket.Conn,
		response *JSONRPCresponse,
	) error {
		return conn.WriteJSON(response) // HARD-DEPENDENCY
	}
}

// Init sets up the
func (s *JSONRPCServerWS) Init() {
	s.connections = map[*websocket.Conn]chan *JSONRPCresponse{}
	s.tick = make(chan bool)
}
