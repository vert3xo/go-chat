package main

import (
	"log"
	"net"
)

func main() {
	s := newServer()
	go s.run()

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Failed to start the server: %s", err)
	}

	defer listener.Close()

	log.Println("Server started on port 8888")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Unable to accept the connection: %s", err)
			continue
		}

		go s.newClient(conn)
	}
}
