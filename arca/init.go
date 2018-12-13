package arca

import (
	"net/http"

	Goods "../examples/goods"
	Users "../examples/users"
	"github.com/gorilla/websocket"
	grid "github.com/rianby64/arca-grid"
)

func init() {
	upgradeConnection = func(
		w http.ResponseWriter,
		r *http.Request,
	) (*websocket.Conn, error) {
		return upgrader.Upgrade(w, r, nil) // HARD-DEPENDENCY
	}
	readJSON = func(
		conn *websocket.Conn,
		request *JSONRPCrequest,
	) error {
		return conn.ReadJSON(&request) // HARD-DEPENDENCY
	}
	writeJSON = func(
		conn *websocket.Conn,
		response *JSONRPCresponse,
	) error {
		return conn.WriteJSON(response) // HARD-DEPENDENCY
	}
}

// Init sets up the
func (s *JSONRPCServerWS) Init() {
	s.connections = map[*websocket.Conn]chan *JSONRPCresponse{}
	s.tick = make(chan bool)

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
		s.Response(&response) // this function MUST be preset in all the grids
	}
	users.RegisterMethod("query", &queryUsers)
	users.Listen(&usersListen)

	s.matchMethod = goods.Query
	// </move-all-this-code-somewhere-else>
}
