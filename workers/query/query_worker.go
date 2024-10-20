package main

import (
	"chatApp/shared"
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
)

type MessageRPCServer string

func (t *MessageRPCServer) ReadAllMessages(args struct{ Receiver string }, reply *[]shared.Message) error {
	folderPath := filepath.Join("persistence", args.Receiver)
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Folder does not exist: %s\n", folderPath)
			*reply = []shared.Message{}
			return nil
		}
		return err
	}

	var messages []shared.Message
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			data, err := ioutil.ReadFile(filepath.Join(folderPath, file.Name()))
			if err != nil {
				return err
			}
			var msg shared.Message
			if err := json.Unmarshal(data, &msg); err != nil {
				return err
			}
			messages = append(messages, msg)
			log.Printf("Message read: %+v\n", msg) // Print the message read from the file
		}
	}

	*reply = messages
	return nil
}

func main() {
	messageRPC := new(MessageRPCServer)
	rpc.Register(messageRPC)
	rpc.HandleHTTP()
	port := ":1122"

	// listen for requests on 1122
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal("listen error: ", err)
	}

	http.Serve(listener, nil)
}
