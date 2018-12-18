package arca

import "github.com/gorilla/websocket"

func (s *JSONRPCServerWS) writeJSON(
	conn *websocket.Conn,
	response *JSONRPCresponse,
) error {
	err := s.transport.writeJSON(conn, response)
	s.tick <- true
	return err
}

func (s *JSONRPCServerWS) tickResponse(
	conn *websocket.Conn,
) {
	for {
		response, ok := <-s.connections[conn]
		if response == nil || !ok {
			s.tick <- false
			break
		}
		go s.writeJSON(conn, response)
		<-s.tick
	}
}

// Broadcast whatever
func (s *JSONRPCServerWS) Broadcast(
	response *JSONRPCresponse,
) {
	for _, conn := range s.connections {
		conn <- response
	}
}
