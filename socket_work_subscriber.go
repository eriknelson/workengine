package main

import (
	"fmt"
	io "github.com/googollee/go-socket.io"
	"log"
)

type SocketWorkSubscriber struct {
	server         *io.Server
	msgBuffer      <-chan string
	clientManifest map[string]io.Socket
}

func NewSocketWorkSubscriber() *SocketWorkSubscriber {
	server, err := io.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	clientManifest := make(map[string]io.Socket)

	s := &SocketWorkSubscriber{
		server:         server,
		clientManifest: clientManifest,
	}

	configureServer(s)
	return s
}

func configureServer(s *SocketWorkSubscriber) {
	s.server.On("connection", func(socket io.Socket) {
		fmt.Printf("Client connected -> %s\n", socket.Id())
		s.clientManifest[socket.Id()] = socket

		socket.On("disconnection", func() {
			fmt.Printf("Client disconnected -> %s\n", socket.Id())
			delete(s.clientManifest, socket.Id())
		})
	})
}

func (s *SocketWorkSubscriber) Subscribe(msgBuffer <-chan string) {
	// Always drain the buffer if there's a message waiting.
	// NOTE: DON'T FORGET TO GOROUTINE THIS, OR WILL YOU CHOKE THE MAIN PROCESSOR
	go func() {
		for {
			msg := <-msgBuffer
			s.broadcastFirehose(msg)
		}
	}()
}

// TODO: Broadcast to job channels instead of the firehose
func (s *SocketWorkSubscriber) broadcastFirehose(msg string) {
	fmt.Printf("E -> %s", msg)
	for _, socket := range s.clientManifest {
		socket.Emit("firehose", msg)
	}
}
