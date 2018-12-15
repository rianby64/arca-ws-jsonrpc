package main

import (
	"log"
	"net/http"

	arca ".."
	Goods "./examples/goods"
	Users "./examples/users"
	grid "github.com/rianby64/arca-grid"
)

func runServer() {

	ws := arca.JSONRPCServerWS{}

	goods := grid.Grid{}
	users := grid.Grid{}
	var queryGoods grid.RequestHandler = func(
		requestParams *interface{},
		context *interface{},
		notify grid.NotifyCallback,
	) (interface{}, error) {
		result, _ := Goods.Read(requestParams, context)
		var usersContext interface{} = map[string]string{
			"source": "Users",
		}
		users.Update(requestParams, &usersContext)
		return result, nil
	}
	goods.RegisterMethod("query", &queryGoods)

	var updateUsers grid.RequestHandler = func(
		requestParams *interface{},
		context *interface{},
		notify grid.NotifyCallback,
	) (interface{}, error) {
		result, _ := Users.Read(requestParams, context)
		notify(result)
		return result, nil
	}
	var usersListen grid.ListenCallback = func(
		message interface{},
		context interface{},
	) {
		var response arca.JSONRPCresponse
		response.Context = context
		response.Method = "update"
		response.Result = message
		ws.Broadcast(&response)
	}
	users.RegisterMethod("update", &updateUsers)
	users.Listen(&usersListen)

	var queryHandler arca.JSONRequestHandler = goods.Query
	ws.RegisterMethod("Goods", "read", &queryHandler)

	http.HandleFunc("/ws", ws.Handle)

	log.Println("Serving")
	go (func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	})()
}
