package arca

import "fmt"

// RegisterMethod whatever
func (s *JSONRPCServerWS) RegisterMethod(
	source string,
	method string,
	handler *JSONRequestHandler,
) error {
	if handler == nil {
		return fmt.Errorf("A handler must be a defined function")
	}
	if source == "" {
		return fmt.Errorf(
			"A Source must be a defined string")
	}
	if method == "" {
		return fmt.Errorf(
			"A method must be a defined string")
	}
	if s.handlers == nil {
		s.handlers = map[string]map[string]*JSONRequestHandler{}
	}
	if s.handlers[source] == nil {
		s.handlers[source] = map[string]*JSONRequestHandler{}
	}
	s.handlers[source][method] = handler
	return nil
}
