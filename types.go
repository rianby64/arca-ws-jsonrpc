package arca

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// JSONRPCBase is the base for both request and response structures
type JSONRPCBase struct {
	ID      string
	Method  string
	Context interface{}
}

// JSONRPCerror is the structure of JSON-RPC response
type JSONRPCerror struct {
	Code    int
	Message string
	Data    interface{}
}

// JSONRPCrequest is the structure of JSON-RPC request
type JSONRPCrequest struct {
	JSONRPCBase
	Params interface{}
}

// JSONRPCresponse is the structure of JSON-RPC response
type JSONRPCresponse struct {
	JSONRPCBase
	Result interface{}
	Error  interface{}
}

// JSONRequestHandler whatever
type JSONRequestHandler func(
	requestParams *interface{},
	context *interface{},
) (interface{}, error)

// JSONRPCServerWS whatever
type JSONRPCServerWS struct {
	connections sync.Map
	tick        chan bool
	handlers    map[string]map[string]*JSONRequestHandler
	transport   connectWithWS
}

type connectWithWS struct {
	upgradeConnection func(
		w http.ResponseWriter,
		r *http.Request,
	) (*websocket.Conn, error)
	readJSON func(
		conn *websocket.Conn,
		request *JSONRPCrequest,
	) error
	writeJSON func(
		conn *websocket.Conn,
		response *JSONRPCresponse,
	) error
	closeConnection func(
		conn *websocket.Conn,
	) error
}
