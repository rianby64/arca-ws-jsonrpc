package arca

import (
	"errors"
	"net/http"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

func createServer(t *testing.T) *JSONRPCServerWS {
	s := JSONRPCServerWS{}
	s.Init()

	if s.tick == nil {
		t.Error("Init() should initiate tick channel")
	}

	if s.connections == nil {
		t.Error("Init() should initiate connections map")
	}
	return &s
}

func Test_readJSON_redefinition(t *testing.T) {
	t.Log("Test readJSON redefinition")

	s := *createServer(t)
	readJSONCalled := false
	s.transport.readJSON = func(*websocket.Conn, *JSONRPCrequest) error {
		readJSONCalled = true
		return nil
	}

	s.readJSON(nil, nil)

	if !readJSONCalled {
		t.Error("readJSON call failed")
	}
}

func Test_closeConnection_redefinition(t *testing.T) {
	t.Log("Test closeConnection redefinition")

	s := *createServer(t)
	closeConnectionCalled := false
	s.transport.closeConnection = func(*websocket.Conn) error {
		closeConnectionCalled = true
		return nil
	}

	conn := &websocket.Conn{}
	s.connections[conn] = make(chan *JSONRPCresponse)
	s.closeConnection(conn)

	if !closeConnectionCalled {
		t.Error("closeConnection call failed")
	}

	if s.connections[conn] != nil {
		t.Error("closeConnection should delete the connection conn")
	}
}

func Test_sendResponse_with_ID(t *testing.T) {
	t.Log("Test sendResponse with ID")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}

	var expectedResult interface{} = "expected result"

	request := JSONRPCrequest{}
	request.Method = "method"
	request.ID = "an-id"

	s.connections[conn] = make(chan *JSONRPCresponse)
	go s.sendResponse(conn, &request, &expectedResult)
	response := <-s.connections[conn]
	actualResult := response.Result

	if actualResult.(string) != expectedResult.(string) {
		t.Error("expected result differs from actual result")
	}

	go (func() {
		s.closeConnection(conn)
	})()
}

func Test_sendResponse_without_ID(t *testing.T) {
	t.Log("Test sendResponse with ID")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}

	var expectedResult interface{} = "expected result"

	request := JSONRPCrequest{}
	request.Method = "method"
	request.ID = ""

	s.connections[conn1] = make(chan *JSONRPCresponse)
	s.connections[conn2] = make(chan *JSONRPCresponse)

	go s.sendResponse(nil, &request, &expectedResult)

	var response1 *JSONRPCresponse
	var actualResult1 interface{}
	go (func() {
		response1 = <-s.connections[conn1]
		actualResult1 = response1.Result
		wg.Done()
	})()

	var response2 *JSONRPCresponse
	var actualResult2 interface{}
	go (func() {
		response2 = <-s.connections[conn2]
		actualResult2 = response2.Result
		wg.Done()
	})()
	wg.Wait()

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

func Test_call_Init_from_Handle(t *testing.T) {
	t.Log("Test call Init from Handle")

	s := JSONRPCServerWS{}
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	expectedDone := errors.New("EOF")
	done := make(chan bool)

	s.transport.upgradeConnection = func(
		http.ResponseWriter,
		*http.Request,
	) (*websocket.Conn, error) {
		done <- true
		return nil, expectedDone
	}

	go s.Handle(nil, nil)
	<-done

	if s.tick == nil {
		t.Error("Init() should initiate tick channel")
	}

	if s.connections == nil {
		t.Error("Init() should initiate connections map")
	}
}
