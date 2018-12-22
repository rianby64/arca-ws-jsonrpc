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

func (s *JSONRPCServerWS) closeConnection(
	conn *websocket.Conn,
) error {
	for {
		connChan, ok := s.connections.Load(conn)
		if !ok {
			break
		}
		close(connChan.(chan *JSONRPCresponse))
		break
	}
	s.connections.Delete(conn)
	return s.transport.closeConnection(conn)
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
		for {
			connChan, ok := s.connections.Load(conn)
			if !ok {
				break
			}
			connChan.(chan *JSONRPCresponse) <- &response
			break
		}
	} else {
		s.Broadcast(&response)
	}
}
