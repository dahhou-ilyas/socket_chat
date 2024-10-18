package client

import (
	"bufio"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
)

func main() {
	u := "ws://localhost:8080/ws"
	log.Printf("Connected to %s...\n", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)

	if err != nil {
		log.Fatal("Erro in connection:", err)
	}

	defer c.Close()

}

func getClientIDFromInput() string {
	fmt.Print("Enter the ID: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}
