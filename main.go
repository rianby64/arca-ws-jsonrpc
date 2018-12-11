package main

import (
	"log"
	"net/http"

	"./arca"
)

func main() {

	Handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := arca.UpgradeConnection(w, r)
		if err != nil {
			log.Println("connecting", err)
			return
		}

		done := make(chan error)
		go arca.ListenAndResponse(conn, done)
		log.Println(<-done)
	}
	http.HandleFunc("/ws", Handler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Serving")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
