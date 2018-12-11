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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  64,
	WriteBufferSize: 64,
}

func init() {
	connections := map[*websocket.Conn]chan *JSONRPCresponse{}
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
		var response JSONRPCresponse
		response.Context = context
		response.Method = "read"
		response.Result = message

		for _, chConn := range connections {
			chConn <- &response
		}
	}
	users.RegisterMethod("query", &queryUsers)
	users.Listen(&usersListen)

	UpgradeConnection = func(w http.ResponseWriter,
		r *http.Request) (*websocket.Conn, error) {
		return upgrader.Upgrade(w, r, nil)
	}
	ReadJSON = func(conn *websocket.Conn, request *JSONRPCrequest) error {
		return conn.ReadJSON(&request)
	}
	WriteJSON = func(conn *websocket.Conn, response *JSONRPCresponse) error {
		return conn.WriteJSON(response)
	}
	CloseConnection = func(conn *websocket.Conn) error {
		return conn.Close()
	}
	ListenAndResponse = func(conn *websocket.Conn, done chan error) {
		connections[conn] = make(chan *JSONRPCresponse)
		go (func() {
			for {
				response := <-connections[conn]
				go WriteJSON(conn, response)
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

			// response
			if len(request.ID) > 0 {
				response.ID = request.ID
			}
			WriteJSON(conn, &response)
		}
	}
}
