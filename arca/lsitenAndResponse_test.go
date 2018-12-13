package arca

import (
	"errors"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_sendResponse_without_ID(t *testing.T) {
	t.Log("Test sendResponse with ID")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}

	var expectedResult interface{} = "expected result"

	request := JSONRPCrequest{}
	request.Method = "method"
	request.ID = ""

	s.connections[conn1] = make(chan *JSONRPCresponse)
	s.connections[conn2] = make(chan *JSONRPCresponse)

	go s.sendResponse(nil, &request, &expectedResult)

	response1 := <-s.connections[conn1]
	actualResult1 := response1.Result

	response2 := <-s.connections[conn2]
	actualResult2 := response2.Result

	if actualResult1.(string) != expectedResult.(string) {
		t.Error("expected result differs from actual result")
	}
	if actualResult2.(string) != expectedResult.(string) {
		t.Error("expected result differs from actual result")
	}

	go (func() {
		s.closeConnection(conn1)
		s.closeConnection(conn2)
	})()
}

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
