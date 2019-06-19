package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)


func main(){


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)


	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	log.Printf("connecting to %s", u.String())


	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	//closing websocket if the program is terminated
	defer c.Close()

	//creating a buffer
	done := make(chan struct{})

	go func() {
		//this is listening to websocket
		defer close(done)
		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			if string(message)=="PING"{
				if err := c.WriteMessage(messageType, []byte("PONG")); err != nil {
		            return
		        }
		    }
			log.Printf("recv: %s", message)
		}
	}()


	for {
		select {
		case <-done:
			// websocket listener stopped due to error in readmessage
			return
		case <-interrupt:
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
		}
	}
}