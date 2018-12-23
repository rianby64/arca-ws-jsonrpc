package arca

import (
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_writeJSON_redefinition(t *testing.T) {
	t.Log("Test writeJSON redefinition")

	s := *createServer(t)
	w := sync.WaitGroup{}
	w.Add(1)
	writeJSONCalled := false
	s.transport.writeJSON = func(*websocket.Conn, *JSONRPCresponse) error {
		writeJSONCalled = true
		w.Done()
		return nil
	}

	go s.writeJSON(nil, nil)
	w.Wait()

	if !writeJSONCalled {
		t.Error("writeJSON call failed")
	}

}

func Test_writeJSON_3x_responses(t *testing.T) {
	t.Log("Test writeJSON sending 3 responses")

	s := *createServer(t)
	w := sync.WaitGroup{}
	w.Add(3)

	var expectedResult1 interface{} = "expected result 1"
	response1 := JSONRPCresponse{}
	response1.Method = "method"
	response1.ID = ""
	response1.Result = expectedResult1

	var expectedResult2 interface{} = "expected result 2"
	response2 := JSONRPCresponse{}
	response2.Method = "method"
	response2.ID = ""
	response2.Result = expectedResult2

	var expectedResult3 interface{} = "expected result 3"
	response3 := JSONRPCresponse{}
	response3.Method = "method"
	response3.ID = ""
	response3.Result = expectedResult3

	writeJSONCalled := false
	i := 0
	s.transport.writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		i++
		writeJSONCalled = true
		actualResponse := (*response).Result
		if i == 1 {
			if actualResponse.(string) != expectedResult1.(string) {
				t.Error("expected result differs from actual result")
			}
		} else if i == 2 {
			if actualResponse.(string) != expectedResult2.(string) {
				t.Error("expected result differs from actual result")
			}
		} else if i == 3 {
			if actualResponse.(string) != expectedResult3.(string) {
				t.Error("expected result differs from actual result")
			}
		}
		w.Done()
		return nil
	}

	go (func() {
		s.writeJSON(nil, &response1)
		s.writeJSON(nil, &response2)
		s.writeJSON(nil, &response3)
	})()
	w.Wait()

	if !writeJSONCalled {
		t.Error("writeJSON call failed")
	}
}

func Test_tickResponse_1call(t *testing.T) {
	t.Log("Test tickResponse when sending one response")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	w := sync.WaitGroup{}
	w.Add(1)
	s.connections.Store(conn, make(chan *JSONRPCresponse))

	var actualResponse JSONRPCresponse
	expectedResponse := JSONRPCresponse{}
	expectedResponse.Method = "method"
	expectedResponse.ID = "an-id"

	s.transport.writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		actualResponse = *response
		w.Done()
		return nil
	}

	go s.tickResponse(conn)
	connChan, _ := s.connections.Load(conn)
	connChan.(chan *JSONRPCresponse) <- &expectedResponse
	w.Wait()

	if expectedResponse.Method != actualResponse.Method ||
		expectedResponse.ID != actualResponse.ID {
		t.Error("expectedResponse differs from actualResponse")
	}

	s.closeConnection(conn)
}

func Test_tickResponse_3call(t *testing.T) {
	t.Log("Test tickResponse when sending 3 responses")

	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	w := sync.WaitGroup{}
	w.Add(3)

	conn := &websocket.Conn{}
	connChan := make(chan *JSONRPCresponse)
	s.connections.Store(conn, connChan)

	var actualResponse1 JSONRPCresponse
	var actualResponse2 JSONRPCresponse
	var actualResponse3 JSONRPCresponse
	expectedResponse1 := JSONRPCresponse{}
	expectedResponse1.Method = "method-1"
	expectedResponse1.ID = "1"

	expectedResponse2 := JSONRPCresponse{}
	expectedResponse2.Method = "method-1"
	expectedResponse2.ID = "2"

	expectedResponse3 := JSONRPCresponse{}
	expectedResponse3.Method = "method-1"
	expectedResponse3.ID = "3"

	s.transport.writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		ID := (*response).ID
		if ID == "1" {
			actualResponse1 = *response
		}
		if ID == "2" {
			actualResponse2 = *response
		}
		if ID == "3" {
			actualResponse3 = *response
		}
		w.Done()
		return nil
	}

	go s.tickResponse(conn)
	connChan <- &expectedResponse1
	connChan <- &expectedResponse2
	connChan <- &expectedResponse3
	w.Wait()

	if expectedResponse1.Method != actualResponse1.Method ||
		expectedResponse1.ID != actualResponse1.ID {
		t.Error("expectedResponse1 differs from actualResponse1")
	}

	if expectedResponse2.Method != actualResponse2.Method ||
		expectedResponse2.ID != actualResponse2.ID {
		t.Error("expectedResponse2 differs from actualResponse2")
	}

	if expectedResponse3.Method != actualResponse3.Method ||
		expectedResponse3.ID != actualResponse3.ID {
		t.Error("expectedResponse3 differs from actualResponse3")
	}

	s.closeConnection(conn)
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
	s.connections.Store(conn1, make(chan *JSONRPCresponse))
	conn2 := &websocket.Conn{}
	s.connections.Store(conn2, make(chan *JSONRPCresponse))

	expectedResponse := JSONRPCresponse{}
	expectedResponse.Method = "method"
	expectedResponse.ID = "an-id"

	go s.Broadcast(&expectedResponse)

	var actualResponse1 *JSONRPCresponse
	go (func() {
		connChan, _ := s.connections.Load(conn1)
		actualResponse1 = <-connChan.(chan *JSONRPCresponse)
		wg.Done()
	})()

	var actualResponse2 *JSONRPCresponse
	go (func() {
		connChan, _ := s.connections.Load(conn2)
		actualResponse2 = <-connChan.(chan *JSONRPCresponse)
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

	s.closeConnection(conn1)
	s.closeConnection(conn2)
}
