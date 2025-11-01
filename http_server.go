package main
import("fmt")
import("net")



func main() {
	fmt.Println("Hello World!")

	// Create listener
	// Go support multiple return falues thus the following is valid

	listener, err := net.Listen("tcp", ":8080")

	if(err != nil){	// err will not be nil if there is an error
		fmt.Println("Failed to create listener")
	} else {
		addr := listener.Addr()
		fmt.Printf("Listening to address: %s\n", addr)
	}

	// continiously check for connections
	for{
		connection, err := listener.Accept()
		if(err != nil){
			fmt.Println("Failed to accept connection")
		} else {
			fmt.Printf("Connection established with: %s\n", connection.RemoteAddr())
		}

		// temporary
		connection.Close()

	}

	listener.Close()

}