package arca

import (
	Goods "../examples/goods"
	Users "../examples/users"
	"github.com/gorilla/websocket"
	grid "github.com/rianby64/arca-grid"
)

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
		s.Broadcast(&response) // this function MUST be preset in all the grids
	}
	users.RegisterMethod("query", &queryUsers)
	users.Listen(&usersListen)

	// </move-all-this-code-somewhere-else>

	s.listenAndResponse = func(conn *websocket.Conn, done chan error) {
		s.connections[conn] = make(chan *JSONRPCresponse)
		go s.tickResponse(conn)
		for {
			var request JSONRPCrequest
			if err := s.readJSON(conn, &request); err != nil {
				done <- err
				delete(s.connections, conn)
				return
			}

			result, err := goods.Query(&request.Params, &request.Context)
			if err != nil {
				done <- err
				delete(s.connections, conn)
				return
			}

			var response JSONRPCresponse
			response.Context = &request.Context
			response.Method = request.Method
			response.Result = result

			if len(request.ID) > 0 {
				response.ID = request.ID
			}

			s.Response(conn, &response) // how to tie up here the grids with methods?
		}
	}
}
