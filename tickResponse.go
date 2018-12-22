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
		connChan, ok := s.connections.Load(conn)
		if !ok {
			s.tick <- false
			break
		}
		response, ok := <-connChan.(chan *JSONRPCresponse)
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
	s.connections.Range(func(_, conn interface{}) bool {
		conn.(chan *JSONRPCresponse) <- response
		return true
	})
}
