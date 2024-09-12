package main

import (
	"fmt"
	"net"
	"slices"
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

func (r *room) listMembers() []string {
	members := make([]string, 0, len(r.members))
	for _, c := range r.members {
		members = append(members, c.username)
	}
	slices.Sort(members)
	return members
}
