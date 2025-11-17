package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// creates a buffered channel
var connectionLimit = make(chan struct{}, 10) // channel of structs

func main() {
	//take command line arguments as port
	if len(os.Args) != 2 {
		fmt.Println("Usage:./http_server <port>")
		os.Exit(1)
	}
	port := os.Args[1]

	// Create listener
	// Go support multiple return falues thus the following is valid

	listener, err := net.Listen("tcp", ":"+port)

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
			continue // should do a new loop iteration
		} else {
			fmt.Printf("Connection established with: %s\n", connection.RemoteAddr())
		}

		// if the channel is full it will block right here
		connectionLimit <- struct{}{} // sends struct{}{} to the buffer (0 memory)

		// temporary
		go handleConn(connection)

	}

}

func handleConn(clientConn net.Conn) {
	defer clientConn.Close()
	defer func() { <-connectionLimit }()
	r := bufio.NewReader(clientConn)
	request, err := http.ReadRequest(r)

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("method:%s\n", request.Method)

	switch request.Method {
	case "GET":
		file := request.URL.Path
		filetype := filepath.Ext(file)

		// useless atm except for checking if its valid filetype
		switch filetype {
		case ".html":
		case ".txt":
		case ".gif":
		case ".jpeg", ".jpg":
		case ".css":
		default:
			clientConn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			return
		}

		// connects to the actual webserver
		webserverConn, err := net.Dial("tcp", "localhost:8080")

		// if we fail to dial webserver
		if err != nil {
			// return response to the client (not the webserver)
			clientConn.Write([]byte("HTTP/1.1 502 Bad Gateway\r\n\r\n"))
			return
		}

		// ask webserver for this page
		query := fmt.Sprintf("GET %s HTTP/1.1\r\n\r\n", file)
		webserverConn.Write([]byte(query))

		io.Copy(clientConn, webserverConn)

	default: // if its not a GET request
		clientConn.Write([]byte("HTTP/1.1 501 Not Implemented\r\n\r\n"))
		return
	}

}
