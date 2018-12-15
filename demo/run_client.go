package main

import (
	"log"

	arca ".."
	"github.com/gorilla/websocket"
)

func runClient() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("Something went wrong", err)
		return
	}

	var context interface{} = map[string]string{
		"source": "Goods",
	}
	request := arca.JSONRPCrequest{}
	request.Method = "read"
	request.Context = context
	request.ID = "an-ID"

	c.WriteJSON(request)

	readFromServer := func() {
		_, data, err := c.ReadMessage()
		if err != nil {
			log.Fatal("Error at reading", err)
		}
		log.Println(string(data))
	}

	// Here is expected to receive two messages...
	readFromServer()
	readFromServer()
}
