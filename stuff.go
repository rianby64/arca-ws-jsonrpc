package arca

import (
	"github.com/gorilla/websocket"
)

func (s *JSONRPCExtensionWS) readJSON(
	conn *websocket.Conn,
	request *JSONRPCrequest,
) error {
	return s.transport.readJSON(conn, request)
}

func (s *JSONRPCExtensionWS) closeConnection(
	conn *websocket.Conn,
) error {
	err := s.transport.closeConnection(conn)
	connChan, ok := s.connections.Load(conn)
	if ok {
		close(connChan.(chan *JSONRPCresponse))
	}
	s.connections.Delete(conn)
	return err
}

func (s *JSONRPCExtensionWS) sendResponse(
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
		connChan, ok := s.connections.Load(conn)
		if ok {
			connChan.(chan *JSONRPCresponse) <- &response
		}
	} else {
		s.Broadcast(&response)
	}
}
