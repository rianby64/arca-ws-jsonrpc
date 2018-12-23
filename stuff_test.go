package arca

import (
	"sync"
	"testing"

	"github.com/gorilla/websocket"
)

func createServer(t *testing.T) *JSONRPCExtensionWS {
	s := JSONRPCExtensionWS{}
	s.Init()

	if s.tick == nil {
		t.Error("Init() should initiate the tick")
	}

	if s.handlers == nil {
		t.Error("Init() should initiate the handler map")
	}

	return &s
}

func Test_readJSON_redefinition(t *testing.T) {
	t.Log("Test readJSON redefinition")

	// setup
	s := *createServer(t)
	readJSONCalled := false // flag

	// redefine the readJSON function
	s.transport.readJSON = func(*websocket.Conn, *JSONRPCrequest) error {
		readJSONCalled = true // expecting this flag to be modified
		return nil
	}

	// excercise
	s.readJSON(nil, nil)

	// verify
	if !readJSONCalled {
		t.Error("readJSON call failed")
	}
}

func Test_closeConnection_redefinition(t *testing.T) {
	t.Log("Test closeConnection redefinition")

	// setup
	s := *createServer(t)
	closeConnectionCalled := false // flag

	// redefine the readJSON function
	s.transport.closeConnection = func(*websocket.Conn) error {
		closeConnectionCalled = true // expecting this flag to be modified
		return nil
	}

	conn := &websocket.Conn{}
	s.connections.Store(conn, make(chan *JSONRPCresponse))

	// excercise
	s.closeConnection(conn)

	// verify
	if !closeConnectionCalled {
		t.Error("closeConnection call failed")
	}

	_, ok := s.connections.Load(conn)
	if ok {
		t.Error("closeConnection should delete the connection conn")
	}
}

func Test_sendResponse_with_ID(t *testing.T) {
	t.Log("Test sendResponse with ID")

	// setup
	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}

	var expectedResult interface{} = "expected result"

	request := JSONRPCrequest{}
	request.Method = "method"
	request.ID = "an-id"

	s.connections.Store(conn, make(chan *JSONRPCresponse))

	// excercise
	go s.sendResponse(conn, &request, &expectedResult)

	// verify
	connChan, ok := s.connections.Load(conn)
	if ok {
		response := <-connChan.(chan *JSONRPCresponse)
		actualResult := response.Result

		// verify
		if actualResult.(string) != expectedResult.(string) {
			t.Error("expected result differs from actual result")
		}
	} else {
		t.Error("unexpected error when loading the response's channel")
	}

	// tear-down
	s.closeConnection(conn)
}

func Test_sendResponse_3_different_responses_with_ID(t *testing.T) {
	t.Log("Test sendResponse 3 different responses with ID")

	// setup
	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	var expectedResult1 interface{} = "expected result 1"
	request1 := JSONRPCrequest{}
	request1.Method = "method"
	request1.ID = "an-id 1"

	var expectedResult2 interface{} = "expected result 2"
	request2 := JSONRPCrequest{}
	request2.Method = "method"
	request2.ID = "an-id 2"

	var expectedResult3 interface{} = "expected result 3"
	request3 := JSONRPCrequest{}
	request3.Method = "method"
	request3.ID = "an-id 3"

	s.connections.Store(conn, make(chan *JSONRPCresponse))

	// excercise
	go (func() {
		s.sendResponse(conn, &request1, &expectedResult1)
		s.sendResponse(conn, &request2, &expectedResult2)
		s.sendResponse(conn, &request3, &expectedResult3)
	})()

	// verify
	connChan, ok := s.connections.Load(conn)
	if ok {
		response := <-connChan.(chan *JSONRPCresponse)
		actualResult1 := response.Result

		response = <-connChan.(chan *JSONRPCresponse)
		actualResult2 := response.Result

		response = <-connChan.(chan *JSONRPCresponse)
		actualResult3 := response.Result

		// verify
		if actualResult1.(string) != expectedResult1.(string) {
			t.Error("expected result differs from actual result")
		}
		if actualResult2.(string) != expectedResult2.(string) {
			t.Error("expected result differs from actual result")
		}
		if actualResult3.(string) != expectedResult3.(string) {
			t.Error("expected result differs from actual result")
		}
	} else {
		t.Error("unexpected error when loading the response's channel")
	}

	// tear-down
	s.closeConnection(conn)
}

func Test_sendResponse_without_ID(t *testing.T) {
	t.Log("Test sendResponse without ID")

	// setup
	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}

	var expectedResult1 interface{} = "expected result 1"
	request1 := JSONRPCrequest{}
	request1.Method = "method"
	request1.ID = ""

	s.connections.Store(conn1, make(chan *JSONRPCresponse))
	s.connections.Store(conn2, make(chan *JSONRPCresponse))

	go s.sendResponse(nil, &request1, &expectedResult1)

	go (func() {
		connChan, ok := s.connections.Load(conn1)
		if ok {
			response := <-connChan.(chan *JSONRPCresponse)
			actualResult1 := response.Result

			if actualResult1.(string) != expectedResult1.(string) {
				t.Error("expected result differs from actual result")
			}
		} else {
			t.Error("unexpected error when loading the response's channel")
		}
		wg.Done()
	})()

	go (func() {
		connChan, ok := s.connections.Load(conn2)
		if ok {
			response := <-connChan.(chan *JSONRPCresponse)
			actualResult1 := response.Result

			if actualResult1.(string) != expectedResult1.(string) {
				t.Error("expected result differs from actual result")
			}
		} else {
			t.Error("unexpected error when loading the response's channel")
		}
		wg.Done()
	})()
	wg.Wait()

	// tear-down
	s.closeConnection(conn1)
	s.closeConnection(conn2)
}

func Test_sendResponse_2conns_3responses_without_ID(t *testing.T) {
	t.Log("Test sendResponse without ID")

	// setup
	s := *createServer(t)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	wg := sync.WaitGroup{}
	wg.Add(2)

	conn1 := &websocket.Conn{}
	conn2 := &websocket.Conn{}

	var expectedResult1 interface{} = "expected result 1"
	request1 := JSONRPCrequest{}
	request1.Method = "method"
	request1.ID = ""

	var expectedResult2 interface{} = "expected result 2"
	request2 := JSONRPCrequest{}
	request2.Method = "method"
	request2.ID = ""

	var expectedResult3 interface{} = "expected result 3"
	request3 := JSONRPCrequest{}
	request3.Method = "method"
	request3.ID = ""

	s.connections.Store(conn1, make(chan *JSONRPCresponse))
	s.connections.Store(conn2, make(chan *JSONRPCresponse))

	go (func() {
		s.sendResponse(nil, &request1, &expectedResult1)
		s.sendResponse(nil, &request2, &expectedResult2)
		s.sendResponse(nil, &request3, &expectedResult3)
	})()

	go (func() {
		connChan, ok := s.connections.Load(conn1)
		if ok {
			response := <-connChan.(chan *JSONRPCresponse)
			actualResult1 := response.Result

			response = <-connChan.(chan *JSONRPCresponse)
			actualResult2 := response.Result

			response = <-connChan.(chan *JSONRPCresponse)
			actualResult3 := response.Result

			if actualResult1.(string) != expectedResult1.(string) {
				t.Error("expected result differs from actual result")
			}
			if actualResult2.(string) != expectedResult2.(string) {
				t.Error("expected result differs from actual result")
			}
			if actualResult3.(string) != expectedResult3.(string) {
				t.Error("expected result differs from actual result")
			}
		} else {
			t.Error("unexpected error when loading the response's channel")
		}
		wg.Done()
	})()

	go (func() {
		connChan, ok := s.connections.Load(conn2)
		if ok {
			response := <-connChan.(chan *JSONRPCresponse)
			actualResult1 := response.Result

			response = <-connChan.(chan *JSONRPCresponse)
			actualResult2 := response.Result

			response = <-connChan.(chan *JSONRPCresponse)
			actualResult3 := response.Result

			if actualResult1.(string) != expectedResult1.(string) {
				t.Error("expected result differs from actual result")
			}
			if actualResult2.(string) != expectedResult2.(string) {
				t.Error("expected result differs from actual result")
			}
			if actualResult3.(string) != expectedResult3.(string) {
				t.Error("expected result differs from actual result")
			}
		} else {
			t.Error("unexpected error when loading the response's channel")
		}
		wg.Done()
	})()
	wg.Wait()

	// tear-down
	s.closeConnection(conn1)
	s.closeConnection(conn2)
}

/*
func Test_call_Init_from_Handle(t *testing.T) {
	t.Log("Test call Init from Handle")

	s := JSONRPCExtensionWS{}
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
}
*/
