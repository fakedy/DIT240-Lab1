package main

import (
	"bufio"
	"fmt"
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
	}
	port := os.Args[1]

	// Create listener
	// Go support multiple return falues thus the following is valid

	listener, err := net.Listen("tcp", port)

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
		default:
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			<-connectionLimit
			return
		}

		relativePath := file[1:]

		conn, err := net.Dial("tcp", "localhost:8080")

		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			response = "HTTP/1.1 404 Not Found\r\n"
			fmt.Println(response)
			<-connectionLimit
			return
		}

		fmt.Fprintf(conn, "GET"+relativePath+"HTTP/1.0\r\n\r\n")

		//response += fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", contentType, len(data))
		//response += fmt.Sprintf("\r\n%s", data)

	case "POST":

	default:
		conn.Write([]byte("HTTP/1.1 501 Not Implemented\r\n\r\n"))
		<-connectionLimit
		return
	}

	conn.Write([]byte(response))
	<-connectionLimit // remove from channel

}
