package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {

	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)

	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)

	for {
		_, err := conn.Read(buf)

		if err != nil {

			if err == io.EOF {
				break
			}

			fmt.Println("Error reading from connection ", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error responding to client")
			break
		}

	}
}
