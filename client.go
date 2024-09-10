package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn        net.Conn
	username    string
	commandChan chan<- command
	room        *room
}

func newClient(conn net.Conn, commandChan chan<- command) *client {
	return &client{
		conn:        conn,
		username:    "anonymous",
		commandChan: commandChan,
		room:        nil,
	}
}

func (c *client) readInput() {
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() {
		input := strings.Split(scanner.Text(), " ")
		if len(input) == 0 {
			continue
		}

		c.commandChan <- command{
			commandType: commandType(input[0]),
			args:        input[1:],
			client:      c,
		}
	}
}

func (c *client) joinRoom(r *room) {
	c.room = r
}

func (c *client) sendMessage(m string) {
	fmt.Fprintf(c.conn, "> %s", m)
}
