package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// creates a buffered channel
var connectionLimit = make(chan struct{}, 10) // channel of structs

func main() {
	fmt.Println("Hello World!")

	// Create listener
	// Go support multiple return falues thus the following is valid

	listener, err := net.Listen("tcp", ":8080")

	if err != nil { // err will not be nil if there is an error
		fmt.Println("Failed to create listener")
	} else {
		addr := listener.Addr()
		fmt.Printf("Listening to address: %s\n", addr)
	}

	// continiously check for connections
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept connection")
		} else {
			fmt.Printf("Connection established with: %s\n", connection.RemoteAddr())
			continue // should do a new loop iteration
		}

		// if the channel is full it will block right here
		connectionLimit <- struct{}{} // sends struct{}{} to the buffer

		// temporary
		go handleConn(connection)

	}

	listener.Close()

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	request, err := http.ReadRequest(r)
	if err != nil {
		fmt.Println("uh oh")
		return
	} else if request.Method == "GET" {
		fmt.Printf("URL Requested: %s\n", request.URL.Path)

		// respond

	}

	<-connectionLimit

}
