package arca

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  64,
	WriteBufferSize: 64,
}

// Handle whatever
func (s *JSONRPCServerWS) Handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := s.upgradeConnection(w, r)
	if err != nil {
		log.Println("connecting", err)
		return
	}

	done := make(chan error)
	go s.listenAndResponse(conn, done)
	log.Println(<-done)
}

func (s *JSONRPCServerWS) upgradeConnection(
	w http.ResponseWriter,
	r *http.Request,
) (*websocket.Conn, error) {
	return upgrader.Upgrade(w, r, nil)
}

func (s *JSONRPCServerWS) readJSON(
	conn *websocket.Conn,
	request *JSONRPCrequest,
) error {
	return conn.ReadJSON(&request)
}

func (s *JSONRPCServerWS) writeJSON(
	conn *websocket.Conn,
	response *JSONRPCresponse,
) error {
	err := conn.WriteJSON(response)
	s.tick <- true
	return err
}

func (s *JSONRPCServerWS) closeConnection(
	conn *websocket.Conn,
) error {
	return conn.Close()
}

// Broadcast whatever
func (s *JSONRPCServerWS) Broadcast(response *JSONRPCresponse) {
	for _, conn := range s.connections {
		conn <- response
	}
}

// Response whatever
func (s *JSONRPCServerWS) Response(
	conn *websocket.Conn,
	response *JSONRPCresponse,
) {
	s.connections[conn] <- response
}

func (s *JSONRPCServerWS) tickResponse(
	conn *websocket.Conn,
) {
	for {
		response := <-s.connections[conn]
		go s.writeJSON(conn, response)
		<-s.tick

	}
}
