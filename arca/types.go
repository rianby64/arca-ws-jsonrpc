package arca

import (
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
type JSONRequestHandler func(requestParams *interface{},
	context *interface{}) (interface{}, error)

// DIRUD whatever
// for a while, this struct will be something rigid. But later I'll
// try to turn it into
type DIRUD struct {
	Describe JSONRequestHandler
	Insert   JSONRequestHandler
	Read     JSONRequestHandler
	Update   JSONRequestHandler
	Delete   JSONRequestHandler
}

// JSONRPCServerWS whatever
type JSONRPCServerWS struct {
	connections       map[*websocket.Conn]chan *JSONRPCresponse
	tick              chan bool
	listenAndResponse func(conn *websocket.Conn, done chan error)
}
