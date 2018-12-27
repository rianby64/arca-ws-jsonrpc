package arca

import (
	"github.com/gorilla/websocket"
)

func (s *JSONRPCExtensionWS) listenAndResponse(
	conn *websocket.Conn,
) error {
	s.connections.Store(conn, make(chan *JSONRPCresponse))
	go s.tickResponse(conn)
	for {
		var request JSONRPCrequest
		if err := s.readJSON(conn, &request); err != nil {
			s.closeConnection(conn)
			return err
		}

		if err := s.processRequest(&request, conn); err != nil {
			return err
		}
	}
}

func (s *JSONRPCExtensionWS) processRequest(
	request *JSONRPCrequest,
	conn *websocket.Conn,
) error {
	handler, err := s.matchHandler(request)
	if err != nil {
		s.closeConnection(conn)
		return err
	}
	result, err := (*handler)(&request.Params, &request.Context)
	if err != nil {
		s.closeConnection(conn)
		return err
	}

	if result != nil {
		s.sendResponse(conn, request, &result)
	}
	return nil
}

// ProcessRequest whatever
func (s *JSONRPCExtensionWS) ProcessRequest(
	request *JSONRPCrequest,
) error {
	return s.processRequest(request, nil)
}
