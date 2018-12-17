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
	close(s.connections[conn])
	delete(s.connections, conn)
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
		s.connections[conn] <- &response
	} else {
		s.Broadcast(&response)
	}
}
