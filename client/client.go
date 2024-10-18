package main

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
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
}

func getClientIDFromInput() string {
	fmt.Print("Enter the ID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
