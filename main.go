package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"strconv"
)

// to save logs in files, run:
// go run ./ >>logs/info.log 2>>logs/err.log
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	s := NewServer()
	ch := NewChat()

	// Clear chat history
	if err := os.Truncate("files/temp.txt", 0); err != nil {
		errLog.Fatalf("Failed to truncate: %v", err)
	}

	// Running incoming commands from users
	go s.runCommands(ch)

	var port string
	switch len(os.Args) {
	case 1:
		port = ":8989"
	case 2:
		if _, err := strconv.Atoi(os.Args[1]); err != nil {
			fmt.Println("[USAGE]: ./TCPChat $port")
			return
		}
		port = ":" + os.Args[1]
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	// Running TCP server
	listener, err := net.Listen("tcp", port)
	if err != nil {
		errLog.Fatalf("unable to start server %s\n", err.Error())
	}

	defer listener.Close()
	infoLog.Printf("listening on the port: %s\n", port)

	// Listening connections to the server
	for {
		conn, err := listener.Accept()
		if err != nil {
			errLog.Printf("unable to accept connection %s\n", err.Error())
			continue
		}

		if len(ch.members) == 10 {
			fmt.Fprintln(conn, "max num of connections is 10. try later")
			conn.Close()
			continue
		}

		var mutex sync.Mutex

		// Creating new client
		go s.NewClient(conn, ch, &mutex)
	}
}
