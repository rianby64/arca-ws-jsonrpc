package arca

import (
	"github.com/gorilla/websocket"
)

func (s *JSONRPCServerWS) readJSON(
	conn *websocket.Conn,
	request *JSONRPCrequest,
) error {
	return s.transport.readJSON(conn, request)
}

func (s *JSONRPCServerWS) writeJSON(
	conn *websocket.Conn,
	response *JSONRPCresponse,
) error {
	err := s.transport.writeJSON(conn, response)
	s.tick <- true
	return err
}

func (s *JSONRPCServerWS) closeConnection(
	conn *websocket.Conn,
) error {
	go (func() {
		s.connections[conn] <- nil
	})()
	delete(s.connections, conn)
	return s.transport.closeConnection(conn)
}

// Broadcast whatever
func (s *JSONRPCServerWS) Broadcast(
	response *JSONRPCresponse,
) {
	for _, conn := range s.connections {
		conn <- response
	}
}

func (s *JSONRPCServerWS) tickResponse(
	conn *websocket.Conn,
) {
	for {
		response := <-s.connections[conn]
		if response == nil {
			s.tick <- false
			break
		}
		go s.writeJSON(conn, response)
		<-s.tick
	}
}

func (s *JSONRPCServerWS) sendResponse(
	conn *websocket.Conn,
	request *JSONRPCrequest,
	result *interface{},
) {
	var response JSONRPCresponse
	response.Context = &request.Context
	response.Method = request.Method
	response.Result = *result

	if len(request.ID) > 0 {
		response.ID = request.ID
		s.connections[conn] <- &response
	} else {
		s.Broadcast(&response)
	}
}
