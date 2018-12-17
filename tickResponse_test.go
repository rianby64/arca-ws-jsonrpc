package arca

import (
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

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
