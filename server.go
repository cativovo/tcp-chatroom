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
		c.sendMessage("")
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
		default:
			fmt.Fprintln(cmd.client.conn, "eyyyy")
		}

		cmd.client.sendMessage("")
	}
}

func (s *Server) join(c *client, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(c.conn, "room name is required")
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

	r.addMember(c)
	c.joinRoom(r)

	s.rooms[roomName] = r

	r.broadCast(c, fmt.Sprintf("%s joined the room", c.username))
	c.sendMessage(fmt.Sprintf("welcome to %s\n", roomName))
}

func (s *Server) sendMessage(c *client, args []string) {
	if len(args) == 0 {
		c.sendMessage("message is required\n")
		return
	}

	if c.room == nil {
		c.sendMessage("join a room before you can message\n")
		return
	}

	c.room.broadCast(c, fmt.Sprintf("%s says: %s", c.username, strings.Join(args, " ")))
}
