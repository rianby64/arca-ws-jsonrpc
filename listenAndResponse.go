package arca

import "github.com/gorilla/websocket"

func (s *JSONRPCServerWS) listenAndResponse(
	conn *websocket.Conn,
	done chan error,
) {
	s.connections[conn] = make(chan *JSONRPCresponse)
	go s.tickResponse(conn)
	for {
		var request JSONRPCrequest
		if err := s.readJSON(conn, &request); err != nil {
			done <- err
			s.closeConnection(conn)
			return
		}

		handler, err := s.matchHandler(&request)
		if err != nil {
			done <- err
			s.closeConnection(conn)
			return
		}
		result, err := (*handler)(&request.Params, &request.Context)
		if err != nil {
			done <- err
			s.closeConnection(conn)
			return
		}

		s.sendResponse(conn, &request, &result)
	}
}
