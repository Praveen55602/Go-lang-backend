package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket/chatApp/user"

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

	clients = make(map[*user.User]bool)
)

// Defines a function to handle WebSocket connections. This function is called when a client requests an upgrade to WebSocket.
// whenever a new clients calls this function go creates a new separate concurrent instance of this function for that perticular client
func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	conn.WriteMessage(websocket.TextMessage, []byte("Enter your username: "))
	_, username, err := conn.ReadMessage()
	println("new user added is ", string(username))

	if err != nil {
		println("error while reading", err.Error())
		return
	}
	newUser := &user.User{Name: string(username), Conn: conn}
	defer func() {
		delete(clients, newUser)
		conn.Close()
	}()

	clients[newUser] = true
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		//shows recived message from the client
		//log.Printf("receivd message is %s", p)

		//iterate over all the connected client and sends them the same message except self
		for client := range clients {
			if client.Conn != conn {
				message := fmt.Sprintf("%s: %s", newUser.Name, string(p))
				err := client.Conn.WriteMessage(messageType, []byte(message))
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
