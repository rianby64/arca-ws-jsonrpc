package arca

import (
	"net/http"

	Goods "../examples/goods"
	Users "../examples/users"
	"github.com/gorilla/websocket"
	grid "github.com/rianby64/arca-grid"
)

// UpgradeConnection whatever
var UpgradeConnection func(w http.ResponseWriter,
	r *http.Request) (*websocket.Conn, error)

// ReadJSON whatever
var ReadJSON func(conn *websocket.Conn, request *JSONRPCrequest) error

// WriteJSON whatever
var WriteJSON func(conn *websocket.Conn, response *JSONRPCresponse) error

// CloseConnection whatever
var CloseConnection func(conn *websocket.Conn) error

// ListenAndResponse whatever
var ListenAndResponse func(conn *websocket.Conn, done chan error)

// Broadcast whatever
var Broadcast func(response *JSONRPCresponse)

// Answer whatever
var Answer func(conn *websocket.Conn, response *JSONRPCresponse)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  64,
	WriteBufferSize: 64,
}

func init() {
	connections := map[*websocket.Conn]chan *JSONRPCresponse{}
	tick := make(chan bool)

	// </move-all-this-code-somewhere-else>
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
		// if return something, it means it becomes the response
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
		var response JSONRPCresponse
		response.Context = context
		response.Method = "read"
		response.Result = message
		Broadcast(&response) // this function MUST be preset in all the grids
	}
	users.RegisterMethod("query", &queryUsers)
	users.Listen(&usersListen)

	// </move-all-this-code-somewhere-else>

	// all the functions below MUST be private
	UpgradeConnection = func(w http.ResponseWriter,
		r *http.Request) (*websocket.Conn, error) {
		return upgrader.Upgrade(w, r, nil)
	}
	ReadJSON = func(conn *websocket.Conn, request *JSONRPCrequest) error {
		return conn.ReadJSON(&request)
	}
	WriteJSON = func(conn *websocket.Conn, response *JSONRPCresponse) error {
		err := conn.WriteJSON(response)
		tick <- true
		return err
	}
	CloseConnection = func(conn *websocket.Conn) error {
		return conn.Close()
	}
	// but not this one. This MUST be public as we want the grids to use it
	Broadcast = func(response *JSONRPCresponse) {
		for _, conn := range connections {
			conn <- response
		}
	}
	Answer = func(conn *websocket.Conn, response *JSONRPCresponse) {
		connections[conn] <- response
	}
	// and also this one becase we want to bind it with the http handler
	ListenAndResponse = func(conn *websocket.Conn, done chan error) {
		connections[conn] = make(chan *JSONRPCresponse)
		go (func() {
			for {
				response := <-connections[conn]
				go WriteJSON(conn, response)
				<-tick
			}
		})()
		for {
			var request JSONRPCrequest
			if err := ReadJSON(conn, &request); err != nil {
				done <- err
				delete(connections, conn)
				return
			}

			result, err := goods.Query(&request.Params, &request.Context)
			if err != nil {
				done <- err
				delete(connections, conn)
				return
			}

			var response JSONRPCresponse
			response.Context = &request.Context
			response.Method = request.Method
			response.Result = result

			if len(request.ID) > 0 {
				response.ID = request.ID
			}

			Answer(conn, &response) // how to tie up here the grids with methods?
		}
	}
}
