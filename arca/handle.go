package arca

import (
	"log"
	"net/http"
)

// Handle whatever
func (s *JSONRPCServerWS) Handle(
	w http.ResponseWriter,
	r *http.Request,
) {
	if s.connections == nil ||
		s.tick == nil {
		s.Init()
	}
	conn, err := s.transport.upgradeConnection(w, r)
	if err != nil {
		log.Println("connecting", err)
		return
	}

	done := make(chan error)
	go s.listenAndResponse(conn, done)
	log.Println(<-done)
}
