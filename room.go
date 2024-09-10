package main

import (
	"fmt"
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadCast(c *client, message string) {
	for k, v := range r.members {
		if k != c.conn.RemoteAddr() {
			fmt.Fprintf(v.conn, "%s\n> ", message)
		}
	}
}

func (r *room) addMember(c *client) {
	r.members[c.conn.RemoteAddr()] = c
}
