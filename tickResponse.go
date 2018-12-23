package arca

import (
	"github.com/gorilla/websocket"
)

func (s *JSONRPCExtensionWS) writeJSON(
	conn *websocket.Conn,
	response *JSONRPCresponse,
) error {
	<-s.tick
	err := s.transport.writeJSON(conn, response)
	go (func() { s.tick <- true })()
	return err
}

func (s *JSONRPCExtensionWS) tickResponse(
	conn *websocket.Conn,
) {
	for {
		connChan, ok := s.connections.Load(conn)
		if ok {
			response, ok := <-connChan.(chan *JSONRPCresponse)
			if response == nil || !ok {
				break
			}
			s.writeJSON(conn, response)
		} else {
			break
		}
	}
}

// Broadcast whatever
func (s *JSONRPCExtensionWS) Broadcast(
	response *JSONRPCresponse,
) {
	s.connections.Range(func(_, conn interface{}) bool {
		conn.(chan *JSONRPCresponse) <- response
		return true
	})
}
