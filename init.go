package arca

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  64,
	WriteBufferSize: 64,
}

// Init sets up the
func (s *JSONRPCExtensionWS) Init() {
	if s.tick == nil {
		s.tick = make(chan bool)
		go (func() { s.tick <- true })()
	}

	if s.handlers == nil {
		s.handlers = map[string]map[string]*JSONRequestHandler{}
	}

	if s.transport.upgradeConnection == nil {
		s.transport.upgradeConnection = func(
			w http.ResponseWriter,
			r *http.Request,
		) (*websocket.Conn, error) {
			return upgrader.Upgrade(w, r, nil) // HARD-DEPENDENCY
		}
	}

	if s.transport.readJSON == nil {
		s.transport.readJSON = func(
			conn *websocket.Conn,
			request *JSONRPCrequest,
		) error {
			return conn.ReadJSON(&request) // HARD-DEPENDENCY
		}
	}

	if s.transport.writeJSON == nil {
		s.transport.writeJSON = func(
			conn *websocket.Conn,
			response *JSONRPCresponse,
		) error {
			return conn.WriteJSON(response) // HARD-DEPENDENCY
		}
	}

	if s.transport.closeConnection == nil {
		s.transport.closeConnection = func(
			conn *websocket.Conn,
		) error {
			return conn.Close() // HARD-DEPENDENCY
		}
	}
}
