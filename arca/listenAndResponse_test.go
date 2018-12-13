package arca

import (
	"errors"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_listenAndResponse_readJSON_returning_error(t *testing.T) {
	t.Log("Test listenAndResponse readJSON returning error")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	done := make(chan error)
	expectedDone := errors.New("EOF")

	readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
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

func Test_listenAndResponse_MatchMethod_error(t *testing.T) {
	t.Log("Test listenAndResponse MatchMethod error")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	done := make(chan error)
	expectedDone := errors.New("EOF")

	readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		return nil
	}

	s.MatchMethod = func(*interface{}, *interface{}) (interface{}, error) {
		return nil, expectedDone
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
func Test_listenAndResponse_readJSON_MatchMethod_OK(t *testing.T) {
	t.Log("Test listenAndResponse readJSON MatchMethod OK")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	done := make(chan error)
	expectedDone := errors.New("EOF")
	expectedResult := "my result"
	var actualResult string
	alreadyReadedJSON := false

	readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		if alreadyReadedJSON {
			return expectedDone
		}
		alreadyReadedJSON = true
		request.Method = "method"
		request.ID = "my-id"
		return nil
	}

	writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		actualResult = response.Result.(string)
		return nil
	}

	s.MatchMethod = func(*interface{}, *interface{}) (interface{}, error) {
		var result interface{} = expectedResult
		return result, nil
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
