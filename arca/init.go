package arca

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

// ListenConnection whatever
var ListenConnection func(conn *websocket.Conn, done chan error)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  64,
	WriteBufferSize: 64,
}

func init() {
	UpgradeConnection = func(w http.ResponseWriter,
		r *http.Request) (*websocket.Conn, error) {
		return upgrader.Upgrade(w, r, nil)
	}
	ReadJSON = func(conn *websocket.Conn, request *JSONRPCrequest) error {
		return conn.ReadJSON(&request)
	}
	WriteJSON = func(conn *websocket.Conn, response *JSONRPCresponse) error {
		return conn.WriteJSON(*response)
	}
	CloseConnection = func(conn *websocket.Conn) error {
		return conn.Close()
	}
	ListenConnection = func(conn *websocket.Conn, done chan error) {
		for {
			var request JSONRPCrequest
			if err := ReadJSON(conn, &request); err != nil {
				done <- err
				return
			}
			log.Println(request)
		}
	}
}

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
