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

var upgradeConnection func(
	w http.ResponseWriter,
	r *http.Request,
) (*websocket.Conn, error)

var readJSON func(
	conn *websocket.Conn,
	request *JSONRPCrequest,
) error

var writeJSON func(
	conn *websocket.Conn,
	response *JSONRPCresponse,
) error

// Handle whatever
func (s *JSONRPCServerWS) Handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := upgradeConnection(w, r)
	if err != nil {
		log.Println("connecting", err)
		return
	}

	done := make(chan error)
	go s.listenAndResponse(conn, done)
	log.Println(<-done)
}

func (s *JSONRPCServerWS) readJSON(
	conn *websocket.Conn,
	request *JSONRPCrequest,
) error {
	return readJSON(conn, request)
}

func (s *JSONRPCServerWS) writeJSON(
	conn *websocket.Conn,
	response *JSONRPCresponse,
) error {
	err := writeJSON(conn, response)
	s.tick <- true
	return err
}

func (s *JSONRPCServerWS) closeConnection(
	conn *websocket.Conn,
) error {
	delete(s.connections, conn)
	return conn.Close()
}

// Response whatever
func (s *JSONRPCServerWS) Response(
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
	response.Result = result

	if len(request.ID) > 0 {
		response.ID = request.ID
		s.connections[conn] <- &response
	} else {
		s.Response(&response)
	}
}

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

		result, err := s.MatchMethod(&request.Params, &request.Context)
		if err != nil {
			done <- err
			s.closeConnection(conn)
			return
		}

		s.sendResponse(conn, &request, &result)
	}
}
