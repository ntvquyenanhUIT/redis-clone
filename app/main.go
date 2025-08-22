package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
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

	// initialize data
	err, store := NewDataStore()

	if err != nil {
		fmt.Println("Error initializing database: ", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn, store)

	}

}

func handleConnection(conn net.Conn, store *Store) {
	defer conn.Close()

	parser := NewRespParser(conn)
	writer := NewRespWriter(conn)

	for {
		value, err := parser.ReadValue()

		if err != nil {

			if err == io.EOF {
				break
			}

			fmt.Println("Error reading from client ", err)
			os.Exit(1)
		}

		command := strings.ToUpper(value.array[0].str)
		args := value.array[1:]

		switch command {
		case "ECHO":
			if len(args) != 1 {
				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'echo' command"}
				writer.Write(errValue)
				continue
			}
			writer.Write(args[0])
		case "PING":
			pong := Value{typ: "string", str: "PONG"}
			writer.Write(pong)
		case "SET":

			argLength := len(args)

			switch argLength {
			case 2:
				store.Set(args[0].str, args[1].str)
				ok := Value{typ: "string", str: "OK"}
				writer.Write(ok)
			case 4:
				store.SetWithTimeOut(args[0].str, args[1].str, args[len(args)-1].str)
				ok := Value{typ: "string", str: "OK"}
				writer.Write(ok)
			default:
				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
				writer.Write(errValue)
				continue
			}

		case "GET":
			if len(args) != 1 {
				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
				writer.Write(errValue)
				continue
			}
			val, ok := store.Get(args[0].str)
			if !ok {
				writer.Write(Value{typ: "null"})
				continue
			}
			writer.Write(Value{typ: "string", str: val})

		case "RPUSH":
			if len(args) != 2 {
				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
				writer.Write(errValue)
				continue
			}

			len, err := store.RPush(args[0].str, args[1].str)

			if err != nil {
				writer.Write(Value{typ: "null"})
				continue
			}

			writer.Write(Value{typ: "int", num: len})
		default:
			unknown := Value{typ: "error", str: fmt.Sprintf("ERR unknown command '%s'", command)}
			writer.Write(unknown)
		}
	}
}
