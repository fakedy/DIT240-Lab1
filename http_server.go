package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
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

func handleConn(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	request, err := http.ReadRequest(r)

	if err != nil {
		fmt.Println(err)
		<-connectionLimit // remove from channel
		return
	}

	fmt.Printf("method:%s\n", request.Method)
	var response = "HTTP/1.1 400 Bad Request\r\n"

	switch request.Method {
	case "GET":
		file := request.URL.Path
		response = "HTTP/1.1 200 OK\r\n"
		filetype := filepath.Ext(file)

		contentType := ""
		switch filetype {
		case ".html":
			contentType = "text/html"
		case ".txt":
			contentType = "text/plain"
		case ".gif":
			contentType = "text/plain"
		case ".jpeg", ".jpg":
			contentType = "image/jpeg"
		case ".css":
			contentType = "text/css"
		}

		response += fmt.Sprintf("Content-Type: %s", contentType)
		fmt.Println(response)

	case "POST":

	}

	conn.Write([]byte(response))
	<-connectionLimit // remove from channel

}
