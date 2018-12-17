package arca

import (
	"net/http"
)

// Handle whatever
func (s *JSONRPCServerWS) Handle(
	w http.ResponseWriter,
	r *http.Request,
) error {
	if s.connections == nil ||
		s.tick == nil {
		s.Init()
	}
	conn, err := s.transport.upgradeConnection(w, r)
	if err != nil {
		return err
	}

	done := make(chan error)
	go s.listenAndResponse(conn, done)
	return <-done
}
