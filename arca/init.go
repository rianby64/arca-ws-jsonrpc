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
func (s *JSONRPCServerWS) Init() {
	s.connections = map[*websocket.Conn]chan *JSONRPCresponse{}
	s.handlers = map[string]map[string]*JSONRequestHandler{}
	s.tick = make(chan bool)

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
