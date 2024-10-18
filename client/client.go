package main

import (
	"bufio"
	"chatApp/shared"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	u := "ws://localhost:8080/ws"
	log.Printf("Connected a %s...\n", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal("Error in connection:", err)
	}
	defer c.Close()

	clientID := getClientIDFromInput()

	err = registerClient(c, clientID)

	if err != nil {
		log.Fatal("Error in register:", err)
		return
	}

}

func registerClient(c *websocket.Conn, clientID string) error {
	msg := shared.Message{
		Text:      clientID,
		Sender:    clientID,
		Receiver:  "server",
		Type:      "new_client",
		Timestamp: time.Now(),
	}
	return c.WriteJSON(msg)
}

func readMessages(c *websocket.Conn) {
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error in read:", err)
			return
		}

		var receivedMsg shared.Message
		err = json.Unmarshal(message, &receivedMsg)
		if err != nil {
			log.Println("Error in unmarshal:", err)
			return
		}

		if receivedMsg.Type == "error" {
			log.Fatal("Erro do servidor: %s\n", receivedMsg.Text)
		} else {
			log.Printf("%s: %s\n", receivedMsg.Sender, receivedMsg.Text)
		}
	}
}

func sendMessage(c *websocket.Conn, clientID string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	message := scanner.Text()

	parts := strings.SplitN(message, " ", 2)

	var recipientID string
	if len(parts) > 1 && strings.HasPrefix(parts[0], "to:") {
		recipientID = strings.TrimPrefix(parts[0], "to:")
		message = parts[1]
	}

	if recipientID == "" {
		log.Println("Recipient ID not provided. Please include recipient ID in the message 'to:<id> >message>'.")
		return
	}

	msg := shared.Message{
		Text:      message,
		Sender:    clientID,
		Receiver:  recipientID,
		Type:      "chat",
		Timestamp: time.Now(),
	}

	err := c.WriteJSON(msg)
	if err != nil {
		log.Println("Erro ao enviar mensagem:", err)
		return
	}

}

func getClientIDFromInput() string {
	fmt.Print("Enter the ID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
