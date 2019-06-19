package main

import (
	"log"
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/websocket"
)
type Client struct {
    conn *websocket.Conn
   	client_id int
   	// last_ping_time time.Time
   	// last_pong_time time.Time


}

type ClientStatus struct{
	last_ping_time time.Time
	// ping_timer time.NewTimer
	// last_pong_time time.Time
	connected bool
}
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}
var clients = make(map[Client] ClientStatus) // connected clients\

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
    err = ws.WriteMessage(1, []byte("Hi Client!"))

    go listen(ws) // listen on the created websocket in a goroutine
	client := Client{ws,1234}
	
	clients[client]=ClientStatus{time.Time{},true}
	// clients_ping_time[client]=
	log.Println(clients)
	// log.Println(clients_ping_time)
	log.Println("-----------------")
	
}


func listen(conn *websocket.Conn){
	//continously listens to all incoming messages for the websocket
	//taking it directly from an example
	for { //infinite loop to keep listening until program is terminated
		messageType, p, err := conn.ReadMessage()
		fmt.Println(string(p))
		if err!=nil{
			log.Println(err)
			// err
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
            // log.Println(err)
            return
        }

	}	
}

// func ping_all_clients() {
	
// }
func ping_all_clients(){
	fmt.Println("starting")
	//setting up timer interval of 30 seconds
	ticker := time.NewTicker(10*time.Second)

	for {
		select{
			case  t:=<- ticker.C:
				//ticker channel go a new value (30s ticker went off)
				// can save outgoing ping time in struct, or can setup a 5 sec timer
				// which can be stopped (in case of pong) or expire(resulting in disconnection)

				// for _,time :=range clients_ping_time{
				// 	fmt.Println("------",time,t)
				// }
				fmt.Println(t.String())
				fmt.Println(clients)
				for client,status := range clients{
					fmt.Println("connected",status.connected,status)
					if status.connected{
						if status.last_ping_time.Add(5*time.Second).Before(time.Now()) && !status.last_ping_time.IsZero(){
							status.connected = false
							clients[client]=status
							client.conn.Close()
							continue
						}
						// check connected clients and send ping
						// fmt.Println(client,connected)	
						err := client.conn.WriteMessage(websocket.TextMessage, []byte("PING"))
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
	// send ping to all connected clients
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
