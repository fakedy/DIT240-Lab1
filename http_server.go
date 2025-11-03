package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

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
		}

		// temporary
		go handleConn(connection)

	}

	listener.Close()

}

func handleConn(conn net.Conn) {

	tmp := make([]byte, 256)
	conn.Read(tmp)

	r := bufio.NewReader(conn)
	request, err := http.ReadRequest(r)
	fmt.Printf("fail:%s\n", request.Method)
	if err != nil {
		fmt.Println("uh oh")
	} else if request.Method == "GET" {
		fmt.Printf("its the nutshack")

	}

}
