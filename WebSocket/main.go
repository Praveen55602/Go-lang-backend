package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Creates an instance of the Upgrader struct from the Gorilla WebSocket library. This struct is used to upgrade an HTTP connection to a WebSocket connection.
var (
	upgrader = websocket.Upgrader{
		//Sets up a function to allow connections from any origin. In a real application, you might want to implement a security check here.
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients = make(map[*websocket.Conn]bool)
)

// Defines a function to handle WebSocket connections. This function is called when a client requests an upgrade to WebSocket.
// whenever a new clients calls this function go creates a new separate concurrent instance of this function for that perticular client
func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		delete(clients, conn)
		conn.Close()
	}()

	clients[conn] = true
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//shows recived message from the client
		log.Printf("receivd message is %s", p)

		//iterate over all the connected client and sends them the same message except self
		for client := range clients {
			if client != conn {
				err := client.WriteMessage(messageType, p)
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/chat", webSocketHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
