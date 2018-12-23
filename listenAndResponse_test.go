package arca

import (
	"errors"
	"fmt"
	"sync"
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
	expectedDone := errors.New("EOF")

	s.transport.readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		return expectedDone
	}

	err := s.listenAndResponse(conn)

	_, ok := s.connections.Load(conn)
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

	err := s.listenAndResponse(conn)

	_, ok := s.connections.Load(conn)
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

	err := s.listenAndResponse(conn)

	_, ok := s.connections.Load(conn)
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
	w := sync.WaitGroup{}
	w.Add(1)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
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
		w.Done()
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

	err := s.listenAndResponse(conn)
	_, ok := s.connections.Load(conn)

	w.Wait()

	if ok {
		t.Error("conn souldn't be present in connections")
	}
	if actualResult != expectedResult {
		t.Errorf("expected '%s' result differs from actual '%s' result",
			expectedResult, actualResult)
	}
	if err != expectedDone {
		t.Error("unexpected done")
	}
}

func Test_listenAndResponse_readJSON_matchHandler_3x_responses_OK(t *testing.T) {
	t.Log("Test listenAndResponse readJSON matchHandler and send 3 responses OK")

	s := *createServer(t)
	w := sync.WaitGroup{}
	w.Add(3)
	s.transport.closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	conn := &websocket.Conn{}
	expectedDone := errors.New("EOF")
	expectedResult1 := "my result 1"
	expectedResult2 := "my result 2"
	expectedResult3 := "my result 3"
	var actualResult1 string
	var actualResult2 string
	var actualResult3 string

	i := 0

	s.transport.readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		i++
		if i > 3 {
			return expectedDone
		}
		var context interface{} = map[string]interface{}{"source": "whatever"}
		request.Context = context
		request.Method = "method"
		request.ID = fmt.Sprintf("my-id-%v", i)
		request.Params = i
		return nil
	}

	s.transport.writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		if response.ID == "my-id-1" {
			actualResult1 = response.Result.(string)
		} else if response.ID == "my-id-2" {
			actualResult2 = response.Result.(string)
		} else if response.ID == "my-id-3" {
			actualResult3 = response.Result.(string)
		}
		w.Done()
		return nil
	}

	var handler JSONRequestHandler = func(
		requestParams *interface{},
		_ *interface{},
	) (interface{}, error) {
		var result interface{}
		if (*requestParams).(int) == 1 {
			result = expectedResult1
		} else if (*requestParams).(int) == 2 {
			result = expectedResult2
		} else if (*requestParams).(int) == 3 {
			result = expectedResult3
		}
		return result, nil
	}
	s.handlers = map[string]map[string]*JSONRequestHandler{
		"whatever": map[string]*JSONRequestHandler{
			"method": &handler,
		},
	}

	err := s.listenAndResponse(conn)
	_, ok := s.connections.Load(conn)

	w.Wait()

	if ok {
		t.Error("conn souldn't be present in connections")
	}
	if actualResult1 != expectedResult1 {
		t.Errorf("expected '%s' result differs from actual '%s' result",
			expectedResult1, actualResult1)
	}
	if actualResult2 != expectedResult2 {
		t.Errorf("expected '%s' result differs from actual '%s' result",
			expectedResult2, actualResult2)
	}
	if actualResult3 != expectedResult3 {
		t.Errorf("expected '%s' result differs from actual '%s' result",
			expectedResult3, actualResult3)
	}
	if err != expectedDone {
		t.Error("unexpected done")
	}
}
