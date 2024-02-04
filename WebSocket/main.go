package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Creates an instance of the Upgrader struct from the Gorilla WebSocket library. This struct is used to upgrade an HTTP connection to a WebSocket connection.
var upgrader = websocket.Upgrader{
	//Sets up a function to allow connections from any origin. In a real application, you might want to implement a security check here.
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Defines a function to handle WebSocket connections. This function is called when a client requests an upgrade to WebSocket.
func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//shows recived message from the client
		log.Printf("receivd message is %s", p)

		//send message to the client
		if err := conn.WriteMessage(messageType, p); err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	http.HandleFunc("/ws", webSocketHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
