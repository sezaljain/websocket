package main

import (
	"log"
	"fmt"
	"net/http"
	"time"
	"math/rand"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)

type ClientStatus struct{
	last_ping_time time.Time
	connected bool
   	client_id string
}

var clients = make(map[*websocket.Conn] ClientStatus) // connected clients\


var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}



func randomString(len int) string {
	rand.Seed(time.Now().UnixNano())
    bytes := make([]byte, len)
    for i := 0; i < len; i++ {
        bytes[i] = byte(rand.Intn(25)+65)
    }
    return string(bytes)
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
    s:=randomString(6)
    log.Println("New Client Connected #"+s)
    err = ws.WriteMessage(1, []byte("Hi Client #"+s+"!"))


    go listen(ws) // listen on the created websocket in a goroutine
	clients[ws]=ClientStatus{time.Time{},true,s}
	
}

func close_client_connection(conn *websocket.Conn){
	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	client_status := clients[conn]
	client_status.connected = false
	clients[conn] = client_status
	// delete(clients,conn)
}

func listen(conn *websocket.Conn){
	
	//continously listens to all incoming messages for the websocket
	//taking it directly from an example
	for { //infinite loop to keep listening until program is terminated
		
		_, p, err := conn.ReadMessage()
		if err!=nil{
			log.Println(err)
			log.Println("Closing client connection")
			return
		}

		client_status := clients[conn]
		if client_status.connected == false {
			conn.Close()
			return

		} else if client_status.last_ping_time.Add(5*time.Second).Before(time.Now()) && !client_status.last_ping_time.IsZero(){
			// more than 5 seconds passed since last ping
			log.Println(string(p),"from client #",client_status.client_id)
			log.Println("Disconnecting as its been more than 5 seconds since ping")
			close_client_connection(conn)
			return

		} else {
			// Only this part of the if-else statements can have no return .. all others should end this goroutine
			log.Println(string(p),"from client! #",client_status.client_id)
			if string(p)=="PONG"{
				client_status.last_ping_time = time.Time{} // resetting timer
				clients[conn] = client_status
			}
		}

	}	
}


func ping_all_clients(){
	//setting up timer interval of 30 seconds
	ticker := time.NewTicker(30*time.Second)

	for {
		select{
			case  t:=<- ticker.C:
				//ticker channel go a new value (30s ticker went off)
				//will save outgoing ping time in client status struct

				for client,status := range clients{
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


func GetClientsEndpoint(w http.ResponseWriter, req *http.Request) {
    s := ""
    for _, item := range clients {
    	if !item.connected{
    		continue
    	} else if s==""{
	    	s += item.client_id
	    } else {
	    	s += ", "+item.client_id	
	    }
    }
    fmt.Fprintf(w, "The connected clients are \n\n%s", s)
}

func main() {
	log.Printf("this is the pinging machine\n")

	router := mux.NewRouter()
	// setting up websocket
	router.HandleFunc("/ws", wsEndpoint)
	// using router to expose only get call on clients end point
	router.HandleFunc("/clients", GetClientsEndpoint).Methods("GET")

	go ping_all_clients()

	log.Fatal(http.ListenAndServe(":8080", router))

}
