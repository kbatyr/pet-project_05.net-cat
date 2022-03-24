package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type server struct {
	commands chan command
}

type command struct {
	*client
	msg string
}

func NewServer() *server {
	return &server{
		commands: make(chan command),
	}
}

// Reading and running commands of users from the channel
func (s *server) runCommands(ch *chat) {
	for cmd := range s.commands {

		current_time := time.Now().Format("2006-01-02 15:04:05")
		ch.broadcast(cmd.client, fmt.Sprintf("\n[%s][%s]: %s", current_time, cmd.nick, cmd.msg))
	}
}

func (s *server) NewClient(conn net.Conn, ch *chat, mutex *sync.Mutex) {

	infoLog.Printf("new client has connected %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "",
		commands: s.commands,
	}

	mutex.Lock()
	ch.members[c.conn.RemoteAddr()] = c
	mutex.Unlock()

	if err := c.readFile("files/linux_logo.txt"); err != nil {
		mutex.Lock()
		errLog.Println(err.Error())
		mutex.Unlock()
		return
	}

	if err := c.readName(ch); err != nil {
		mutex.Lock()
		if err != io.EOF {
			errLog.Println(err.Error())
		}
		mutex.Unlock()
		return
	}

	if err := ch.joinToChat(c); err != nil {
		mutex.Lock()
		errLog.Println(err.Error())
		mutex.Unlock()
		return
	}
	c.readInput(ch)
}
