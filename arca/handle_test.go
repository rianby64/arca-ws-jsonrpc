package arca

import (
	"errors"
	"net/http"
	"testing"

	"github.com/gorilla/websocket"
)

func Test_Handle_upgradeConnection_OK(t *testing.T) {
	t.Log("Test Handle function")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	done := make(chan bool)
	expectedDone := errors.New("EOF")
	expectedResult := "my result"
	var actualResult string
	alreadyReadedJSON := false

	upgradeConnection = func(
		http.ResponseWriter,
		*http.Request,
	) (*websocket.Conn, error) {
		return &websocket.Conn{}, nil
	}

	readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		if alreadyReadedJSON {
			return expectedDone
		}
		alreadyReadedJSON = true
		request.Method = "method"
		request.ID = "my-id"
		return nil
	}

	writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		actualResult = response.Result.(string)
		done <- true
		return nil
	}

	s.matchHandler = func(request *JSONRPCrequest) (*JSONRequestHandler, error) {
		var handler JSONRequestHandler = func(
			*interface{},
			*interface{},
		) (interface{}, error) {
			var result interface{} = expectedResult
			return result, nil
		}
		return &handler, nil
	}

	go s.Handle(nil, nil)
	<-done

	if actualResult != expectedResult {
		t.Error("expected result differs from actual result")
	}
}

func Test_Handle_upgradeConnection_error(t *testing.T) {
	t.Log("Test Handle function")

	s := *createServer(t)
	closeConnection = func(conn *websocket.Conn) error {
		return nil
	}

	expectedDone := errors.New("EOF")

	upgradeConnection = func(
		http.ResponseWriter,
		*http.Request,
	) (*websocket.Conn, error) {
		return nil, expectedDone
	}

	readJSON = func(_ *websocket.Conn, request *JSONRPCrequest) error {
		t.Error("readJSON must be unreachable")
		return nil
	}

	writeJSON = func(_ *websocket.Conn, response *JSONRPCresponse) error {
		t.Error("writeJSON must be unreachable")
		return nil
	}

	s.matchHandler = func(request *JSONRPCrequest) (*JSONRequestHandler, error) {
		t.Error("matchHandler must be unreachable")
		var handler JSONRequestHandler = func(
			*interface{},
			*interface{},
		) (interface{}, error) {
			t.Error("handler method must be unreachable")
			return nil, nil
		}
		return &handler, nil
	}

	s.Handle(nil, nil)
}
