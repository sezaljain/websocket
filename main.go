package main

import (
	"log"
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Simple http request
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}


func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Request on ws endpoint
    // fmt.Fprintf(w, "Hello World")
    // not worrying about cors etc for now. need to be able to setup the server without looking into security first
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }
    // helpful log statement to show connections
    log.Println("Client Connected")
    log.Println("Client Connected",ws)

    listen(ws) // listen on the created websocket
}


func listen(conn *websocket.Conn){
	//continously listens to all incoming messages for the websocket
	//taking it directly from an example
	for { //infinite loop to keep listening until program is terminated
		messageType, p, err := conn.ReadMessage()
		fmt.Println(string(p),messageType,err)

	}	
}

func main() {
	fmt.Printf("hello, world\n")
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", wsEndpoint)
	log.Fatal(http.ListenAndServe(":8000", nil))
}