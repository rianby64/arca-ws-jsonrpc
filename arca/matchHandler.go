package arca

import "fmt"

func (s *JSONRPCServerWS) matchHandler(
	request *JSONRPCrequest,
) (*JSONRequestHandler, error) {
	method := request.Method
	if method == "" {
		return nil, fmt.Errorf("Method must be present in request")
	}
	if request.Context == nil {
		return nil, fmt.Errorf("Context must be present in request")
	}
	contextRequest, ok := request.Context.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Context must be an Object")
	}
	if contextRequest["source"] == nil {
		return nil, fmt.Errorf("Context must define a source")
	}
	source, ok := contextRequest["source"].(string)
	if !ok {
		return nil, fmt.Errorf(
			"Context has an incorrect source expecting an string")
	}
	if source == "" {
		return nil, fmt.Errorf(
			"Source '%s' in Context must be a defined string", source)
	}
	handler := s.handlers[source][method]
	if handler == nil {
		return nil, fmt.Errorf(
			"handler for source '%s' and method '%s' is nil", source, method)
	}
	return handler, nil
}
