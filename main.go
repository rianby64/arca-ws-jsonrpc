package main

import (
	"log"
	"net/http"

	"./arca"
)

func main() {

	http.HandleFunc("/ws", arca.Handle)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Serving")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
