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

		go handleConn(connection)

	}

}

func handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	request, err := http.ReadRequest(r)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("method:%s\n", request.Method)
	switch request.Method {
	case "GET":
		file := request.URL.Path
		fmt.Println(file)
	case "POST":
		fmt.Println("its the nutshack")
	}

}
