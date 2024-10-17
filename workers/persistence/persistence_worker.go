package persistence

import (
	"chatApp/shared"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
)

type MessageRPCServer string

func (t *MessageRPCServer) PersistMessage(msg shared.Message, reply *string) error {

	log.Printf("Received message: %s from sender: %s, receiver: %s\n", msg.Text, msg.Sender, msg.Receiver)
	folderName := "persistence/" + msg.Receiver
	if _, err := os.Stat(folderName); os.IsNotExist(err) {
		if err := os.MkdirAll(folderName, 0755); err != nil {
			return err
		}
	}

	filename := folderName + "/message_" + msg.Sender + "_" + msg.Timestamp.Format("2006_01_02_15_04_05") + ".json"
	file, err := os.Create(filename)

	if err != nil {
		log.Printf("Error creating file: %v", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(msg); err != nil {
		return err
	}

	*reply = "Message Persisted"
	return nil
}

func main() {

	messageRPC := new(MessageRPCServer)

	rpc.Register(messageRPC)
	rpc.HandleHTTP()

	port := ":1123"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	http.Serve(listener, nil)
}
