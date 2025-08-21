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

	parser := NewRespParser(conn)

	for {
		value, err := parser.ReadValue()

		if err != nil {

			if err == io.EOF {
				break
			}

			fmt.Println("Error reading from client ", err)
			os.Exit(1)
		}

		command := value.array[0].str
		args := value.array[1:]

		if command == "echo" || command == "ECHO" {
			if len(args) != 1 {
				conn.Write([]byte("-ERR wrong number of arguments for 'echo' command\r\n"))
				continue
			}

			response := fmt.Sprintf("$%d\r\n%s\r\n", len(args[0].str), args[0].str)
			conn.Write([]byte(response))
		} else {
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}
