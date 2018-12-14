package arca

func (s *JSONRPCServerWS) matchHandler(
	*JSONRPCrequest,
) (*JSONRequestHandler, error) {
	return &s.handlerMatched, nil
}
