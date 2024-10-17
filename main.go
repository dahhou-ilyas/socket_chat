package main

import (
	"chatApp/shared"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/rpc"
	"time"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	chat_broadcast       = make(chan shared.Message)
	persist_broadcast    = make(chan shared.Message)
	historical_broadcast = make(chan shared.Message)
	clients              = make(map[string]*websocket.Conn)
)

//TIP To run your code, right-click the code and select <b>Run</b>. Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.

func handleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)

		return
	}

	defer conn.Close()

	for {
		var msg shared.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			return
		}
		if msg.Type == "new_client" {

		} else if msg.Type == "chat" {
			chat_broadcast <- msg
			persist_broadcast <- msg
		} else {
			log.Println("Unknown message:", msg.Type)
		}
	}
}

func addNewUser(msg shared.Message, conn *websocket.Conn) bool {
	clientID := msg.Text

	if _, exists := clients[clientID]; exists {
		log.Println("Already exists client:", clientID)

		err := conn.WriteJSON(shared.Message{
			Text:      "User already exists",
			Sender:    "server",
			Receiver:  clientID,
			Type:      "error",
			Timestamp: time.Now(),
		})
		if err != nil {
			log.Println(err)
		}
		return false
	} else {
		clients[clientID] = conn
		log.Println("Added new client:", clientID)
		return true
	}
}

func handleChatMessages() {
	for {
		msg := <-chat_broadcast

		log.Println("message sender : %s , message receveire : %s , contenue : %s", msg.Sender, msg.Receiver, msg.Text)

		if conn, ok := clients[msg.Receiver]; ok {

			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println(err)
				conn.Close()
				delete(clients, msg.Receiver)
			}
		}
	}
}

func handleHistoricMessage() {
	for {
		msg := <-historical_broadcast

		if conn, ok := clients[msg.Receiver]; ok {
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Println(err)
				conn.Close()
				delete(clients, msg.Receiver)
			}
		}
	}
}

func forwardMessagesToRPC() {
	for {
		msg := <-persist_broadcast
		client, err := rpc.DialHTTP("tcp", "localhost:1123")

		if err != nil {
			log.Println("Failed to connect to RPC server:", err)
			continue
		}
		defer client.Close()

		var reply string
		err = client.Call("MessageRPCServer.PersistMessage", msg, &reply)
		if err != nil {
			log.Println("Failed to persist message:", err)
			continue
		}
		log.Println("RPC service response:", reply)

	}
}

func readAllMessagesFromRPC(receiver string) {
	client, err := rpc.DialHTTP("tcp", "localhost:1122")
	if err != nil {
		log.Println("Failed to connect to RPC server:", err)
		return
	}

	defer client.Close()

	var reply []shared.Message

	args := struct {
		Receiver string
	}{
		Receiver: receiver,
	}

	err = client.Call("MessageRPCServer.ReadAllMessages", args, &reply)
	if err != nil {
		log.Println("Failed to read message:", err)
	}

	for _, msg := range reply {
		historical_broadcast <- msg
	}
}

func removeClient(clientID string) {
	delete(clients, clientID)
	log.Printf("Client '%s' removed\n", clientID)
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	log.Println("server listen in port 8080")

	go handleHistoricMessage()
	go handleChatMessages()
	go forwardMessagesToRPC()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error in listen server 8080", err)
	}
}

func test1() {
	for i := 0; i < 1111; i++ {
		fmt.Println("test1 :", i)
	}
}

func test2() {
	for i := 0; i < 222; i++ {
		fmt.Println("test2 :", i)
	}
}
