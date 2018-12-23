package arca

import (
	"net/http"
)

// Handle whatever
func (s *JSONRPCExtensionWS) Handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	if s.tick == nil {
		s.Init()
	}
	conn, err := s.transport.upgradeConnection(w, r)
	if err != nil {
		return
	}

	s.listenAndResponse(conn)
}
