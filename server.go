package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	rooms       map[string]*room
	commandChan chan command
}

func NewServer() Server {
	return Server{
		rooms:       make(map[string]*room),
		commandChan: make(chan command),
	}
}

func (s *Server) ListenAndServe(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("running at %s\n", addr)

	go s.runCommand()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %s", err)
			continue
		}

		log.Printf("%s connected", conn.RemoteAddr())

		c := newClient(conn, s.commandChan)
		s.listRooms(c)
		go c.readInput()
	}
}

func (s *Server) runCommand() {
	log.Println("listening to commands")
	for cmd := range s.commandChan {
		switch cmd.commandType {
		case CmdJoinRoom:
			s.join(cmd.client, cmd.args)
		case CmdSendMessage:
			s.sendMessage(cmd.client, cmd.args)
		case CmdListRooms:
			s.listRooms(cmd.client)
		case CmdListMembers:
			s.listMembers(cmd.client)
		case CmdSetUsername:
			s.setUsername(cmd.client, cmd.args)
		case CmdQuit:
			s.quit(cmd.client)
		default:
			cmd.client.sendMessage(fmt.Sprintf("invalid command %s", cmd.commandType))
		}
	}
}

func (s *Server) join(c *client, args []string) {
	if len(args) == 0 {
		c.sendMessage("room name is required")
		return
	}

	roomName := args[0]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
	}

	if c.room != nil {
		s.quitRoom(c)
	}

	r.addMember(c)
	c.joinRoom(r)

	s.rooms[roomName] = r

	r.broadCast(c, fmt.Sprintf("(%s) joined the room", c.username))
	c.sendMessage(fmt.Sprintf("welcome to %s", roomName))
}

func (s *Server) sendMessage(c *client, args []string) {
	if len(args) == 0 {
		c.sendMessage("message is required")
		return
	}

	if c.room == nil {
		c.sendMessage("join a room before you can message")
		return
	}

	c.room.broadCast(c, fmt.Sprintf("(%s) says: %s", c.username, strings.Join(args, " ")))
	fmt.Fprint(c.conn, "> ")
}

func (s *Server) listRooms(c *client) {
	if len(s.rooms) == 0 {
		c.sendMessage("no rooms found, create one using /join ROOM_NAME")
		return
	}

	var rooms strings.Builder
	for room := range s.rooms {
		rooms.WriteString(fmt.Sprintf("  - %s\n", room))
	}

	c.sendMessage(fmt.Sprintf("available rooms:\n%s", strings.TrimSuffix(rooms.String(), "\n")))
}

func (s *Server) listMembers(c *client) {
	if c.room == nil {
		c.sendMessage("must be inside of a room to list the members")
		return
	}
	c.sendMessage(fmt.Sprintf("members:\n  - %s", strings.Join(c.room.listMembers(), "\n  - ")))
}

func (s *Server) setUsername(c *client, args []string) {
	if len(args) == 0 || args[0] == "" {
		c.sendMessage("username cannot be empty")
		return
	}

	oldUsername := c.username
	c.username = args[0]

	c.sendMessage(fmt.Sprintf("changed username from (%s) to (%s)", oldUsername, c.username))

	if c.room != nil {
		c.room.broadCast(c, fmt.Sprintf("(%s) changed their username to (%s)", oldUsername, c.username))
	}
}

func (s *Server) quit(c *client) {
	s.quitRoom(c)
	c.sendMessage("it was nice having you here, take care!")
	c.conn.Close()
}

func (s *Server) quitRoom(c *client) {
	delete(s.rooms[c.room.name].members, c.conn.RemoteAddr())

	c.room.broadCast(c, fmt.Sprintf("(%s) left the room", c.username))
	c.quitRoom()
}
