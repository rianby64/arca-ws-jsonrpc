package arca

import (
	"testing"
)

func Test_registerMethod_without_handler(t *testing.T) {
	t.Log("RegisterMethod fails if no handler defined in request")

	s := *createServer(t)

	err := s.RegisterMethod("source", "method", nil)
	t.Log(err)
	if err == nil {
		t.Error("a handler function must be provided")
	}
}

func Test_registerMethod_without_source(t *testing.T) {
	t.Log("Registersource fails if no source defined in request")

	s := *createServer(t)
	var handler JSONRequestHandler = func(
		requestParams *interface{},
		context *interface{},
	) (interface{}, error) {
		return nil, nil
	}

	err := s.RegisterMethod("", "method", &handler)
	t.Log(err)
	if err == nil {
		t.Error("a source string must be provided")
	}
}

func Test_registerMethod_without_method(t *testing.T) {
	t.Log("RegisterMethod fails if no method defined in request")

	s := *createServer(t)
	var handler JSONRequestHandler = func(
		requestParams *interface{},
		context *interface{},
	) (interface{}, error) {
		return nil, nil
	}

	err := s.RegisterMethod("source", "", &handler)
	t.Log(err)
	if err == nil {
		t.Error("a method string must be provided")
	}
}

func Test_registerMethod_OK(t *testing.T) {
	t.Log("RegisterMethod fails if no method defined in request")

	s := *createServer(t)
	var handler JSONRequestHandler = func(
		requestParams *interface{},
		context *interface{},
	) (interface{}, error) {
		return nil, nil
	}

	err := s.RegisterMethod("source", "method", &handler)
	t.Log(err)
	if err != nil {
		t.Error("unexpected error")
	}
}
