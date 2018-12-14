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

func Test_matchHandler_request_with_Context_with_incorrect_source(t *testing.T) {
	t.Log("Match a handler fails if context defined in request contains an incorrect source")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	request := &JSONRPCrequest{}
	request.Method = "method"
	request.Context = map[string]interface{}{"source": 123456}

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

func Test_matchHandler_request_with_Context_with_empty_source(t *testing.T) {
	t.Log("Match a handler fails if context defined in request contains an empty source")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	request := &JSONRPCrequest{}
	request.Method = "method"
	request.Context = map[string]interface{}{"source": ""}

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

func Test_matchHandler_request_with_Context_with_source_but_nil_handler(t *testing.T) {
	t.Log("Match a handler fails if context with source defined in request is a nil handler")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}
	request := &JSONRPCrequest{}
	request.Method = "method"
	request.Context = map[string]interface{}{"source": "whatever"}

	s.handlers = map[string]map[string]*JSONRequestHandler{
		"whatever": map[string]*JSONRequestHandler{
			"method": nil,
		},
	}

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
