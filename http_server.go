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

func handleConn(conn net.Conn) {
	//close conn when we are done
	defer conn.Close()
	defer func() { <-connectionLimit }()
	//read conn
	r := bufio.NewReader(conn)
	//get request from conn
	request, err := http.ReadRequest(r)

	//if request is bad then close the conn
	if err != nil {
		fmt.Println(err)
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
		return
	}

	//identify the http request method
	fmt.Printf("method:%s\n", request.Method)
	//default response to request
	var response = "HTTP/1.1 400 Bad Request\r\n"

	switch request.Method {
	case "GET":
		//get the file path
		file := request.URL.Path

		fmt.Printf("Requested URL: %s\n", file)
		//temporary default response to GET request
		response = "HTTP/1.1 200 OK\r\n"
		filetype := filepath.Ext(file)

		//get the content type from the path
		validType, contentType := validType(filetype)
		if !validType {
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			return
		}

		//check if file exists on server if not respond with 404
		relativePath := file[1:]
		data, err := os.ReadFile(relativePath)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			response = "HTTP/1.1 404 Not Found\r\n"
			fmt.Println(response)
			return
		}

		//add information to response
		response += fmt.Sprintf("Content-Type: %s\r\nContent-Length: %d\r\n", contentType, len(data))
		response += fmt.Sprintf("\r\n%s", data)

	case "POST":
		// get path from post request
		file := request.URL.Path
		// remove the / in the beginning of the path
		relativePath := file[1:]

		//check if file created has a valid file extension
		validFiletype, _ := validType(relativePath)
		if !validFiletype {
			conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
			return
		}

		//create file in directory according to path
		data, err := os.Create(relativePath)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
			return
		}
		//copy contents of POST body into the newly created file
		_, err = io.Copy(data, request.Body)
		if err != nil {
			conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
			return
		}

		response = "HTTP/1.1 200 OK\r\n\r\n"
		defer data.Close()

	//respond with 501 for every http method that isn't GET and POST
	default:
		conn.Write([]byte("HTTP/1.1 501 Not Implemented\r\n\r\n"))
		return
	}

	//write the response to the connection
	conn.Write([]byte(response))

}

// check for valid file extensions and return the content-type
func validType(file string) (bool, string) {
	filetype := filepath.Ext(file)

	switch filetype {
	case ".html":
		return true, "text/html"
	case ".txt":
		return true, "text/plain"
	case ".gif":
		return true, "image/gif"
	case ".jpeg", ".jpg":
		return true, "image/jpeg"
	case ".css":
		return true, "text/css"
	default:
		return false, ""
	}
}
