package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
			switch {
			case len(args) == 2:
				len, err := store.RPush(args[0].str, args[1].str)
				fmt.Printf("len: %v, err: %v", len, err)
				if err != nil {
					writer.Write(Value{typ: "null"})
					continue
				}
				writer.Write(Value{typ: "int", num: len})
			case len(args) > 2:
				listLen := 0
				for i := 1; i < len(args); i++ {
					len, err := store.RPush(args[0].str, args[i].str)

					if err != nil {
						writer.Write(Value{typ: "null"})
						continue
					}
					listLen = len
				}
				writer.Write(Value{typ: "int", num: listLen})

			default:

				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'rpush' command"}
				writer.Write(errValue)
				continue

			}
		case "LPUSH":
			switch {
			case len(args) == 2:
				len, err := store.LPush(args[0].str, args[1].str)

				if err != nil {
					writer.Write(Value{typ: "null"})
					continue

				}
				writer.Write(Value{typ: "int", num: len})
			case len(args) > 2:
				listLen := 0
				for i := 1; i < len(args); i++ {
					len, err := store.LPush(args[0].str, args[i].str)

					if err != nil {
						writer.Write(Value{typ: "null"})
						continue
					}
					listLen = len
				}
				writer.Write(Value{typ: "int", num: listLen})

			default:

				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'lpush' command"}
				writer.Write(errValue)
				continue

			}

		case "LRANGE":
			switch {
			case len(args) == 3:
				start, err := strconv.Atoi(args[1].str)
				if err != nil {
					errValue := Value{typ: "error", str: "ERR Cannot parse range"}
					writer.Write(errValue)
					continue
				}
				end, err := strconv.Atoi(args[2].str)
				if err != nil {
					errValue := Value{typ: "error", str: "ERR Cannot parse range"}
					writer.Write(errValue)
					continue
				}

				result, err := store.LRange(args[0].str, start, end)
				if err != nil {
					errValue := Value{typ: "error", str: "ERR error getting list range"}
					writer.Write(errValue)
					continue
				}
				arr := Value{
					array: make([]Value, 0),
				}
				for _, val := range result {
					str := Value{
						typ: "string",
						str: val,
					}
					arr.array = append(arr.array, str)
				}
				writer.Write(Value{typ: "array", array: arr.array})

			default:
				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'lrange' command"}
				writer.Write(errValue)
				continue

			}

		case "LLEN":
			if len(args) != 1 {
				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'llen' command"}
				writer.Write(errValue)
				continue
			}

			result, err := store.LLen(args[0].str)
			if err != nil {
				errValue := Value{typ: "error", str: err.Error()}
				writer.Write(errValue)
				continue
			}

			writer.Write(Value{typ: "int", num: result})

		case "LPOP":

		LPOP_SWITCH:
			switch {
			case len(args) == 1:
				val, hasDeleted, err := store.LPop(args[0].str)
				if err != nil {
					errValue := Value{typ: "error", str: err.Error()}
					writer.Write(errValue)
					continue
				}

				if !hasDeleted {
					writer.Write(Value{typ: "null"})
					continue
				}

				writer.Write(Value{typ: "string", str: val})
			case len(args) == 2:
				timeToExecute, err := strconv.Atoi(args[1].str)

				if err != nil {
					writer.Write(Value{typ: "error", str: fmt.Sprintf("ERR invalid number: %v", args[0].str)})
					continue
				}

				arr := make([]Value, 0)

				for range timeToExecute {
					val, hasDeleted, err := store.LPop(args[0].str)

					if err != nil {
						errValue := Value{typ: "error", str: err.Error()}
						writer.Write(errValue)
						break LPOP_SWITCH
					}

					if !hasDeleted {
						writer.Write(Value{typ: "null"})
						break LPOP_SWITCH
					}
					str := Value{typ: "string", str: val}
					arr = append(arr, str)
				}
				writer.Write(Value{typ: "array", array: arr})

			default:
				errValue := Value{typ: "error", str: "ERR wrong number of arguments for 'LPOP' command"}
				writer.Write(errValue)
				continue
			}
		case "BLPOP":
			timeout, err := strconv.Atoi(args[1].str)

			if err != nil {
				writer.Write(Value{
					typ: "error", str: fmt.Sprintf("ERR converting string: '%s'", args[1].str),
				})
				continue
			}
			result, err := store.BLPop(args[0].str, timeout)
			if err != nil {
				writer.Write(Value{
					typ: "error", str: fmt.Sprintf("ERR performing BLPop: '%s'", args[1].str),
				})
				continue
			}
			fmt.Println(result)
			if result == "" {
				writer.Write(Value{
					typ: "null",
				})
				continue
			}
			writer.Write(Value{
				typ: "array", array: []Value{
					{typ: "string", str: args[0].str},
					{typ: "string", str: result},
				},
			})

		default:
			unknown := Value{typ: "error", str: fmt.Sprintf("ERR unknown command '%s'", command)}
			writer.Write(unknown)
		}
	}
}
