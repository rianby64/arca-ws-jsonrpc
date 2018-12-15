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

func Test_writeJSON_redefinition(t *testing.T) {
	t.Log("Test writeJSON redefinition")

	s := *createServer(t)
	writeJSONCalled := false
	s.transport.writeJSON = func(*websocket.Conn, *JSONRPCresponse) error {
		writeJSONCalled = true
		return nil
	}

	go s.writeJSON(nil, nil)
	<-s.tick

	if !writeJSONCalled {
		t.Error("writeJSON call failed")
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

func Test_tickResponse_1call(t *testing.T) {
	t.Log("Test tickResponse when sending one response")

	s := *createServer(t)
	conn := &websocket.Conn{}
	s.connections[conn] = make(chan *JSONRPCresponse)

	var actualResponse JSONRPCresponse
	expectedResponse := JSONRPCresponse{}
	expectedResponse.Method = "method"
	expectedResponse.ID = "an-id"

	s.transport.writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		actualResponse = *response
		return nil
	}

	go s.tickResponse(conn)
	s.connections[conn] <- &expectedResponse
	s.connections[conn] <- nil

	if expectedResponse.Method != actualResponse.Method ||
		expectedResponse.ID != actualResponse.ID {
		t.Error("expectedResponse differs from actualResponse")
	}

	if <-s.tick {
		t.Error("tick is open")
	}
}

func Test_tickResponse_2call(t *testing.T) {
	t.Log("Test tickResponse when sending two responses")

	s := *createServer(t)
	conn := &websocket.Conn{}
	s.connections[conn] = make(chan *JSONRPCresponse)

	var actualResponse1 JSONRPCresponse
	var actualResponse2 JSONRPCresponse
	expectedResponse1 := JSONRPCresponse{}
	expectedResponse1.Method = "method-1"
	expectedResponse1.ID = "1"

	expectedResponse2 := JSONRPCresponse{}
	expectedResponse2.Method = "method-1"
	expectedResponse2.ID = "2"

	s.transport.writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		ID := (*response).ID
		if ID == "1" {
			actualResponse1 = *response
		}
		if ID == "2" {
			actualResponse2 = *response
		}
		return nil
	}

	go s.tickResponse(conn)
	s.connections[conn] <- &expectedResponse1
	s.connections[conn] <- &expectedResponse2
	s.connections[conn] <- nil

	if expectedResponse1.Method != actualResponse1.Method ||
		expectedResponse1.ID != actualResponse1.ID {
		t.Error("expectedResponse1 differs from actualResponse1")
	}

	if expectedResponse2.Method != actualResponse2.Method ||
		expectedResponse2.ID != actualResponse2.ID {
		t.Error("expectedResponse2 differs from actualResponse2")
	}

	if <-s.tick {
		t.Error("tick is open")
	}
}

func Test_Broadcast(t *testing.T) {
	t.Log("Test Broadcast")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(2)
	conn1 := &websocket.Conn{}
	s.connections[conn1] = make(chan *JSONRPCresponse)
	conn2 := &websocket.Conn{}
	s.connections[conn2] = make(chan *JSONRPCresponse)

	expectedResponse := JSONRPCresponse{}
	expectedResponse.Method = "method"
	expectedResponse.ID = "an-id"

	go s.Broadcast(&expectedResponse)

	var actualResponse1 *JSONRPCresponse
	go (func() {
		actualResponse1 = <-s.connections[conn1]
		wg.Done()
	})()

	var actualResponse2 *JSONRPCresponse
	go (func() {
		actualResponse2 = <-s.connections[conn2]
		wg.Done()
	})()
	wg.Wait()

	if expectedResponse.Method != actualResponse1.Method ||
		expectedResponse.ID != actualResponse1.ID {
		t.Error("expectedResponse differs from actualResponse1")
	}

	if expectedResponse.Method != actualResponse2.Method ||
		expectedResponse.ID != actualResponse2.ID {
		t.Error("expectedResponse differs from actualResponse2")
	}

	go (func() {
		s.closeConnection(conn1)
		s.closeConnection(conn2)
	})()
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
