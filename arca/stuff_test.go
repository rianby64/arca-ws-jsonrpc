package arca

import (
	"testing"

	"github.com/gorilla/websocket"
)

func callInit(t *testing.T) *JSONRPCServerWS {
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
	readJSONCalled := false

	readJSON = func(*websocket.Conn, *JSONRPCrequest) error {
		readJSONCalled = true
		return nil
	}

	s := *callInit(t)
	s.readJSON(nil, nil)

	if !readJSONCalled {
		t.Error("readJSON call failed")
	}
}

func Test_writeJSON_redefinition(t *testing.T) {
	t.Log("Test writeJSON redefinition")
	writeJSONCalled := false

	writeJSON = func(*websocket.Conn, *JSONRPCresponse) error {
		writeJSONCalled = true
		return nil
	}

	s := *callInit(t)
	go s.writeJSON(nil, nil)
	<-s.tick

	if !writeJSONCalled {
		t.Error("writeJSON call failed")
	}
}

func Test_closeConnection_redefinition(t *testing.T) {
	t.Log("Test closeConnection redefinition")
	closeConnectionCalled := false

	closeConnection = func(*websocket.Conn) error {
		closeConnectionCalled = true
		return nil
	}

	s := *callInit(t)
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
	s := *callInit(t)
	conn := &websocket.Conn{}
	s.connections[conn] = make(chan *JSONRPCresponse)

	var actualResponse JSONRPCresponse
	expectedResponse := JSONRPCresponse{}
	expectedResponse.Method = "method"
	expectedResponse.ID = "an-id"

	writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
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
	s := *callInit(t)
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

	writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
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
