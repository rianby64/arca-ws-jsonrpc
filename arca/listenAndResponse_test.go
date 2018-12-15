package arca

import (
	"errors"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_listenAndResponse_readJSON_returning_error(t *testing.T) {
	t.Log("Test listenAndResponse readJSON returning error")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	done := make(chan error)
	expectedDone := errors.New("EOF")

	s.transport.readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		return expectedDone
	}

	go s.listenAndResponse(conn, done)
	err := <-done

	_, ok := s.connections[conn]
	if ok {
		t.Error("conn souldn't be present in connections")
	}
	if err != expectedDone {
		t.Error("unexpected done")
	}
}

func Test_listenAndResponse_matchHandler_error(t *testing.T) {
	t.Log("Test listenAndResponse matchHandler error")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	done := make(chan error)
	expectedDone := errors.New("EOF")
	alreadyReadedJSON := false

	s.transport.readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		if alreadyReadedJSON {
			return expectedDone
		}
		alreadyReadedJSON = true
		var context interface{} = map[string]interface{}{"source": "whatever"}
		request.Context = context
		request.Method = "method"
		request.ID = "my-id"
		return nil
	}

	var handler JSONRequestHandler = func(
		*interface{},
		*interface{},
	) (interface{}, error) {
		return nil, expectedDone
	}
	s.handlers = map[string]map[string]*JSONRequestHandler{
		"whatever": map[string]*JSONRequestHandler{
			"method": &handler,
		},
	}

	go s.listenAndResponse(conn, done)
	err := <-done

	_, ok := s.connections[conn]
	if ok {
		t.Error("conn souldn't be present in connections")
	}
	if err != expectedDone {
		t.Error("unexpected done")
	}
}

func Test_listenAndResponse_matchHandler_with_nil_handler_error(t *testing.T) {
	t.Log("Test listenAndResponse matchHandler with nil handler ends up in error")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	done := make(chan error)
	expectedDone := errors.New("EOF")
	alreadyReadedJSON := false

	s.transport.readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		if alreadyReadedJSON {
			return expectedDone
		}
		alreadyReadedJSON = true
		var context interface{} = map[string]interface{}{"source": "whatever"}
		request.Context = context
		request.Method = "method"
		request.ID = "my-id"
		return nil
	}

	s.handlers = map[string]map[string]*JSONRequestHandler{
		"whatever": map[string]*JSONRequestHandler{
			"method": nil,
		},
	}

	go s.listenAndResponse(conn, done)
	err := <-done

	_, ok := s.connections[conn]
	if ok {
		t.Error("conn souldn't be present in connections")
	}
	if err == nil {
		t.Error("unexpected done")
	}
}

func Test_listenAndResponse_readJSON_matchHandler_OK(t *testing.T) {
	t.Log("Test listenAndResponse readJSON matchHandler OK")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	done := make(chan error)
	expectedDone := errors.New("EOF")
	expectedResult := "my result"
	var actualResult string
	alreadyReadedJSON := false

	s.transport.readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		if alreadyReadedJSON {
			return expectedDone
		}
		alreadyReadedJSON = true
		var context interface{} = map[string]interface{}{"source": "whatever"}
		request.Context = context
		request.Method = "method"
		request.ID = "my-id"
		return nil
	}

	s.transport.writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		actualResult = response.Result.(string)
		return nil
	}

	var handler JSONRequestHandler = func(
		*interface{},
		*interface{},
	) (interface{}, error) {
		var result interface{} = expectedResult
		return result, nil
	}
	s.handlers = map[string]map[string]*JSONRequestHandler{
		"whatever": map[string]*JSONRequestHandler{
			"method": &handler,
		},
	}

	go s.listenAndResponse(conn, done)
	err := <-done

	_, ok := s.connections[conn]
	if ok {
		t.Error("conn souldn't be present in connections")
	}
	if err != expectedDone {
		t.Error("unexpected done")
	}
	if actualResult != expectedResult {
		t.Error("expected result differs from actual result")
	}
}
