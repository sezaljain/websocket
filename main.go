package main

import (
	"log"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)

type ClientStatus struct{
	last_ping_time time.Time
	connected bool
   	client_id int
}
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
var clients = make(map[*websocket.Conn] ClientStatus) // connected clients\

func handler(w http.ResponseWriter, r *http.Request) {
	// Simple http request
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}


func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Request on ws endpoint
    // not worrying about cors etc for now. need to be able to setup the server without looking into security first

    // upgrading http request to websocket
    upgrader.CheckOrigin = func(r *http.Request) bool { return true }

    ws, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
    }
    // helpful log statement to show connections
    log.Println("New Client Connected",ws)
    err = ws.WriteMessage(1, []byte("Hi Client!"))


    go listen(ws) // listen on the created websocket in a goroutine
	
	clients[ws]=ClientStatus{time.Time{},true,1234}
	log.Println(clients)
	
}


func listen(conn *websocket.Conn){
	//continously listens to all incoming messages for the websocket
	//taking it directly from an example
	for { //infinite loop to keep listening until program is terminated
		_, p, err := conn.ReadMessage()
		fmt.Println(string(p))
		if err!=nil{
			log.Println(err)
			conn.Close()
			delete(clients,conn)
			log.Println("Closing client connection",conn)
			return
		}
		client_status := clients[conn]
		if client_status.connected == false {
			conn.Close()
			continue
		} else if client_status.last_ping_time.Add(5*time.Second).Before(time.Now()) && !client_status.last_ping_time.IsZero(){
				client_status.connected = false
				clients[conn] = client_status
				conn.Close()
		} else {
			// if p == "PONG"
			client_status.last_ping_time = time.Time{} // resetting timer
			clients[conn] = client_status
		}

	}	
}


func ping_all_clients(){
	//setting up timer interval of 30 seconds
	ticker := time.NewTicker(1*time.Second)

	for {
		select{
			case  t:=<- ticker.C:
				//ticker channel go a new value (30s ticker went off)
				//will save outgoing ping time in client status struct

				for client,status := range clients{
					fmt.Println("connected",status.connected,status)
					if status.connected{
						if status.last_ping_time.Add(5*time.Second).Before(time.Now()) && !status.last_ping_time.IsZero(){
							status.connected = false
							clients[client]=status
							client.Close()
							continue
						}
						// check connected clients and send ping
						err := client.WriteMessage(websocket.TextMessage, []byte("PING"))
						if err != nil {
				            log.Println(err)
				            return
				        }
				        status.last_ping_time = t
				        clients[client] = status
					}
					
				}
		}
	}
	// ping sent to all connected clients
}

func main() {
	fmt.Printf("this is the pinging machine\n")
	http.HandleFunc("/", handler)
	http.HandleFunc("/ws", wsEndpoint)
	go ping_all_clients()

	log.Fatal(http.ListenAndServe(":8080", nil))

}


// func ping_all_clients(){z
// 	md:=make(chan)
// }

// //a client can join at any time. every 30 s it should receive a ping
// // get request for all connected clients (so http handler should be functional)
// func tick(out chan <- [2]float64){

//     c := time.NewTicker(time.Millisecond *500)
//     for range c.C{
//         out <- mark
//     }
// }

// func main() {

//     fmt.Println("Start")

//     md := make(chan [2]float64)
//     go tick(md)

//     for range <-md{
//         fmt.Println(<-md)
//     }
// }

// setup webscoket, add clients to a list
// keep listening to messages, identify with a client
