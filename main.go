package main

import (
	"log"
	"net/http"
	"time"

	"./arca"
	//grid "github.com/rianby64/arca-grid"
)

func main() {

	Handler := func(w http.ResponseWriter, r *http.Request) {
		conn, err := arca.UpgradeConnection(w, r)

		if err != nil {
			log.Println("connecting", err)
			return
		}

		listenCh := make(chan error)
		go arca.ListenConnection(conn, listenCh)
		go (func() {
			var context interface{} = map[string]string{
				"source": "a source",
			}
			var result interface{} = map[string]string{
				"result": "a result",
			}
			var response = arca.JSONRPCresponse{}
			response.ID = "whatever"
			response.Method = "a method"
			response.Context = context
			response.Result = result

			time.AfterFunc(time.Second*1, func() {
				arca.WriteJSON(conn, &response)
			})
		})()

		log.Println(<-listenCh)
	}
	http.HandleFunc("/ws", Handler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	/*
		grid := grid.Grid{}
		log.Println(grid)
	*/

	log.Println("Serving")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
