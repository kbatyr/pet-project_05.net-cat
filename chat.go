package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

// list of connected clients
// the key is ip-adress of client, and value is type client
type chat struct {
	members map[net.Addr]*client
}

func NewChat() *chat {
	return &chat{
		members: make(map[net.Addr]*client),
	}
}

// Writes to tmp file client's messages with current time
func (c *client) msgHistory(msg string) {

	file, err := os.OpenFile("files/temp.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	writer := bufio.NewWriter(file)
	if err != nil {
		errLog.Println(err.Error())
		return
	}
	defer file.Close()

	current_time := time.Now().Format("2006-01-02 15:04:05")
	usr := fmt.Sprintf("[%s][%s]: ", current_time, c.nick)

	writer.WriteString(usr + msg + "\n")
	writer.Flush()
}

// Inform other users if new user joined to the chat
func (ch *chat) joinToChat(c *client) error {

	c.inChat = true

	// Uploads chat history
	if err := c.readFile("files/temp.txt"); err != nil {
		return err
	}
	ch.broadcast(c, fmt.Sprintf("\n%s has joined the chat...\n", c.nick))
	return nil
}

// Delete user from chat and close connection with server
func (ch *chat) quitChat(c *client) {

	delete(ch.members, c.conn.RemoteAddr())
	ch.broadcast(c, fmt.Sprintf("\n%s has left our chat...\n", c.nick))
	c.inChat = false
	c.conn.Close()
}

// Writing pattern msg to all users terminal
// F.e: [2021-12-13 16:59:32][nick]:
func selfBroadcast(c *client) {
	current_time := time.Now().Format("2006-01-02 15:04:05")
	usr := fmt.Sprintf("[%s][%s]: ", current_time, c.nick)
	c.printMsg(usr)
}

// Send message from sender to other clients
func (chat *chat) broadcast(sender *client, msg string) {

	// ranging over map with current members in chat
	for addr, member := range chat.members {

		// does not send sender msg to himself
		if addr != sender.conn.RemoteAddr() {

			// send msg only for users in chat and with valid nicks
			// after receiving msg from sender print to terminal the pattern for the next msgs
			if member.inChat && sender.nick != "" {
				member.printMsg(msg)
				selfBroadcast(member)
			}
		} else {
			// after sending msg to users print to terminal the pattern for the next msgs
			selfBroadcast(sender)
		}
	}
}
