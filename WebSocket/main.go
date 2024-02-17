package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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
	//conn.ReadMessage pauses the processing of the function until client sends some message
	_, username, err := conn.ReadMessage()
	println("new user added is ", string(username))

	if err != nil {
		println("error while reading", err.Error())
		return
	}
	user := &user.User{Name: string(username), Conn: conn}
	defer func() {
		delete(clients, user)
		conn.Close()
	}()

	clients[user] = true
	for {
		//here also in the for loop it will wait for the client to send some request then only process the code further
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}

		if name, msg := isPrivateMessage(string(p)); name != "" && msg != "" {
			reciever := findUserByName(name)
			fmt.Println("message to received ", name)

			if reciever == nil {
				fmt.Println("no user with this username found")
				continue
			}

			privateMessageHandler(user, reciever, msg)
			continue
		}

		//iterate over all the connected client and sends them the same message except self
		for client := range clients {
			if client != user {
				message := fmt.Sprintf("%s: %s", user.Name, string(p))
				err := client.Conn.WriteMessage(messageType, []byte(message))
				if err != nil {
					fmt.Println(err)
					return
				}
			}
		}
	}
}

func privateMessageHandler(sender *user.User, reciever *user.User, messageToSend string) {
	message := fmt.Sprintf("%s: %s", sender.Name, messageToSend)
	reciever.RecieveMessage(message)
}

func isPrivateMessage(message string) (string, string) {
	userAndMessage := strings.Split(message, ":")
	if len(userAndMessage) != 2 {
		return "", ""
	}
	return userAndMessage[0], userAndMessage[1]
}

func findUserByName(name string) *user.User {
	for client := range clients {
		if client.Name == name {
			return client
		}
	}
	return nil
}

func main() {
	http.HandleFunc("/chat", webSocketHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
