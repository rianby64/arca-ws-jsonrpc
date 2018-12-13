package main

import (
	"log"
	"net/http"

	"./arca"
	Goods "./examples/goods"
	Users "./examples/users"
	grid "github.com/rianby64/arca-grid"
)

func main() {

	ws := arca.JSONRPCServerWS{}
	ws.Init()

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
		users.Query(requestParams, &usersContext)
		return result, nil
	}
	goods.RegisterMethod("query", &queryGoods)

	var queryUsers grid.RequestHandler = func(
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
		response.Method = "read"
		response.Result = message
		ws.Broadcast(&response)
	}
	users.RegisterMethod("query", &queryUsers)
	users.Listen(&usersListen)

	ws.MatchMethod = goods.Query

	http.HandleFunc("/ws", ws.Handle)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Println("Serving")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
