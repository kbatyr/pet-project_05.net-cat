package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	inChat   bool
	commands chan<- command
}

func (c *client) readName(ch *chat) error {

	for c.nick == "" {

		fmt.Fprintf(c.conn, "[ENTER YOUR NAME]: ")

		// reading input data from connection
		name, err := bufio.NewReader(c.conn).ReadString('\n')
		if err == io.EOF {
			infoLog.Printf("client has disconnected %s", c.conn.RemoteAddr().String())
			ch.quitChat(c)
			return err
		}

		if err != nil {
			return err
		}

		name = strings.Trim(name, "\r\n")
		if name == "" {
			continue
		}
		c.nick = name
		c.printMsg(fmt.Sprintf("all right, I will call you %s\n", c.nick))
	}
	return nil
}

func (c *client) readInput(ch *chat) {

	for {
		// reading input data from connection
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err == io.EOF {
			infoLog.Printf("client has disconnected %s", c.conn.RemoteAddr().String())
			ch.quitChat(c)
			return
		}

		if err != nil {
			errLog.Println(err)
			c.printMsg(fmt.Sprint("invalid input, try again"))
			continue
		}

		// putting input data to the channel
		c.commands <- command{
			client: c,
			msg:    msg,
		}

		msg = strings.Trim(msg, "\r\n")

		// save user message to the file if not empty
		if msg != "" {
			c.msgHistory(msg)
		}
	}
}

// Prints (writes) string exactly to the client's terminal
func (c *client) printMsg(msg string) {
	c.conn.Write([]byte(msg))
	// 	alternative way:
	// fmt.Fprintf(c.conn, msg))
}

func (c *client) readFile(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		c.printMsg(fmt.Sprint("internal server error: try later\n"))
		c.conn.Close()
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			c.printMsg(fmt.Sprint("internal server error: try later\n"))
			c.conn.Close()
			return err
		}
		c.printMsg(line)
	}
	return nil
}
