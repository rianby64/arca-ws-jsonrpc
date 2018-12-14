package arca

import (
	"testing"

	"github.com/gorilla/websocket"
)

func Test_matchHandler_request_without_Method(t *testing.T) {
	t.Log("Match a handler fails if no method defined in request")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	request := &JSONRPCrequest{}

	handler, err := s.matchHandler(request)
	t.Log(err)
	if handler == nil {
		if err == nil {
			t.Error("nil handler must lead to an error")
		}
	} else {
		t.Error("a Method must be defined at Handler")
	}
}

func Test_matchHandler_request_without_Context(t *testing.T) {
	t.Log("Match a handler fails if no context defined in request")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	request := &JSONRPCrequest{}
	request.Method = "method"

	handler, err := s.matchHandler(request)
	t.Log(err)
	if handler == nil {
		if err == nil {
			t.Error("nil handler must lead to an error")
		}
	} else {
		t.Error("a Context must be defined at Handler")
	}
}

func Test_matchHandler_request_with_incorrect_Context(t *testing.T) {
	t.Log("Match a handler fails if context defined in request is not an object")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	request := &JSONRPCrequest{}
	request.Method = "method"
	request.Context = []string{}

	handler, err := s.matchHandler(request)
	t.Log(err)
	if handler == nil {
		if err == nil {
			t.Error("nil handler must lead to an error")
		}
	} else {
		t.Error("a Context must be defined at Handler")
	}
}

func Test_matchHandler_request_with_Context_without_source(t *testing.T) {
	t.Log("Match a handler fails if context defined in request doesn't contain a source")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	request := &JSONRPCrequest{}
	request.Method = "method"
	request.Context = map[string]interface{}{"whatever": ""}

	handler, err := s.matchHandler(request)
	t.Log(err)
	if handler == nil {
		if err == nil {
			t.Error("nil handler must lead to an error")
		}
	} else {
		t.Error("The Context must contain a source")
	}
}

/*
func Test_matchHandler_request_with_Context_with_source(t *testing.T) {
	t.Log("Match a handler if context defined in request contains a source")
	conn := websocket.Conn{}
	source := "source-defined"
	method := "read"
	methods := DIRUD{
		Read: func(requestParams *interface{},
			context *interface{}, response chan interface{}) error {
			response <- nil
			return nil
		},
	}
	RegisterSource(source, &methods)
	request := JSONRPCrequest{}
	request.Method = method
	request.Context = map[string]interface{}{"source": source}

	handler, err := matchHandler(&request, &conn)
	if err != nil {
		t.Error("Unexpected error", err)
	}
	if handler == nil {
		t.Errorf("The Context must match the handler [%s][%s]", source, method)
		if err == nil {
			t.Error("nil handler must lead to an error")
		}
	} else {
		var ptHandler = reflect.ValueOf(*handler).Pointer()
		var ptMethod = reflect.ValueOf(methods.Read).Pointer()
		if ptHandler == ptMethod {
			t.Log("Matched handler and given method are the same")
		} else {
			t.Error("Matched handler differs from given method")
		}
	}
	setupGlobals()
}
*/
