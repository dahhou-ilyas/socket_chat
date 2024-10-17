package main

import (
	"chatApp/shared"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
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

func removeClient(clientID string) {
	delete(clients, clientID)
	log.Printf("Client '%s' removed\n", clientID)
}

func main() {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.
	s := "gopher"
	fmt.Println("Hello and welcome, %s!", s)

	for i := 1; i <= 5; i++ {
		//TIP You can try debugging your code. We have set one <icon src="AllIcons.Debugger.Db_set_breakpoint"/> breakpoint
		// for you, but you can always add more by pressing <shortcut actionId="ToggleLineBreakpoint"/>. To start your debugging session,
		// right-click your code in the editor and select the <b>Debug</b> option.
		fmt.Println("i =", 100/i)
	}
}

//TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
// Also, you can try interactive lessons for GoLand by selecting 'Help | Learn IDE Features' from the main menu.
