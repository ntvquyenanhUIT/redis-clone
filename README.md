<h1 align="center">Go Redis Server</h1>
<p align="center">
A Redis-compatible server built from scratch in Go.
</p>
<p align="center">
This project is an implementation of Redis as part of the <a href="https://codecrafters.io">CodeCrafters</a> "Build your own Redis" challenge.
</p>

<a name="table-of-contents"></a>

## Table of contents

-   [Table of contents](#table-of-contents)
-   [Description](#description)
-   [Features](#features)
-   [Usage](#usage)
    -   [Requirements](#requirements)
    -   [Running the server](#running-the-server)
-   [Project Structure](#project-structure)
-   [License](#license)

<a name="description"></a>

## Description

This is a custom implementation of the Redis server in Go. It focuses on handling concurrent client connections and parsing the Redis Serialization Protocol (RESP). The server supports basic commands like `PING`, `ECHO`, `SET`, and `GET`, including key expiration.

The primary goal of this project is to learn about networking, protocol parsing, and concurrent programming in Go by building a simplified version of a real-world, high-performance system.

<a name="features"></a>

## Features

*   **TCP Server**: Listens for and accepts concurrent client connections on port `6379`.
*   **RESP Parser**: A custom parser for the Redis Serialization Protocol (RESP) to decode commands from clients.
*   **Core Commands**: Supports `PING`, `ECHO`, `SET`, and `GET`, including key expiration with `PX`.
*   **List Commands**: Implements a list data structure with support for `LPUSH`, `RPUSH`, `LPOP`, `BLPOP`, `LLEN`, and `LRANGE`.
*   **In-Memory Store**: A thread-safe, in-memory key-value store with support for Time-To-Live (TTL) on keys.


<a name="usage"></a>

## Usage

<a name="requirements"></a>

### Requirements

*   Go 1.24 or later.
*   A Unix-like shell (like Git Bash on Windows) to run the helper script.

<a name="running-the-server"></a>

### Running the server

You can build and run the server locally using the provided shell script. This script compiles the Go source files and executes the resulting binary.

```bash
./your_program.sh
```

Once running, the server will listen on `0.0.0.0:6379`. You can connect to it using a Redis client like `redis-cli`.

```bash
# In a new terminal
redis-cli

# Example commands
127.0.0.1:6379> PING
PONG
127.0.0.1:6379> ECHO "hello world"
"hello world"
127.0.0.1:6379> SET name "GitHub Copilot"
OK
127.0.0.1:6379> GET name
"GitHub Copilot"
127.0.0.1:6379> SET key value PX 1000
OK
127.0.0.1:6379> GET key
"value"
# (wait 1 second)
127.0.0.1:6379> GET key
(nil)
```

For a detailed guide on supported list commands, see [list.md](list.md).


<a name="project-structure"></a>

## Project Structure

The project is organized into several files within the `app/` directory:

*   [`app/main.go`](app/main.go): The main entry point for the application. It sets up the TCP listener and handles incoming connections in separate goroutines.
*   [`app/resp.go`](app/resp.go): Contains the logic for parsing and marshalling data according to the RESP protocol. It defines the `Value` struct and the `RespParser` to read from the client stream.
*   [`app/store.go`](app/store.go): Implements the thread-safe in-memory data store using a map and a mutex. It handles `GET`, `SET`, and key expiration logic.
*   [`app/writer.go`](app/writer.go): A helper to abstract away writing RESP-formatted data back to the client.
*   [`your_program.sh`](your_program.sh): A utility script to compile and run the server for local development and testing.

<a name="license"></a>

## License

This project is not licensed for distribution. It is a personal learning project created as part of the CodeCrafters challenge.