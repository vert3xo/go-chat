package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type Server struct {
	rooms    map[string]*Room
	commands chan Command
}

func newServer() *Server {
	return &Server{
		rooms:    make(map[string]*Room),
		commands: make(chan Command),
	}
}

func (s *Server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_HELP:
			s.help(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}

func (s *Server) newClient(conn net.Conn) {
	log.Printf("New client has connected: %s", conn.RemoteAddr().String())

	c := &Client{
		conn:     conn,
		nick:     "anonymous",
		commands: s.commands,
		prompt:   "> ",
	}

	c.readInput()
}

func (s *Server) nick(c *Client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("From now on, you shall be know as %s", c.nick))
}

func (s *Server) join(c *Client, args []string) {
	roomName := args[1]
	r, ok := s.rooms[roomName]
	if !ok {
		r = &Room{
			name:    roomName,
			members: make(map[net.Addr]*Client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c

	s.quitCurrent(c)

	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room.", c.nick))

	c.msg(fmt.Sprintf("Welcome to %s, %s", r.name, c.nick))
}

func (s *Server) listRooms(c *Client, args []string) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("Available rooms are: %s", strings.Join(rooms, ", ")))
}

func (s *Server) msg(c *Client, args []string) {
	if c.room == nil {
		c.err(errors.New("You must join a room first!"))
		return
	}

	c.room.broadcast(c, c.nick+": "+strings.Join(args[1:], " "))
}

func (s *Server) help(c *Client, args []string) {
	c.msg(fmt.Sprintf("Help:\n/nick\tSet a nickname\n/join\tJoin a room\n/msg\tAnother way to talk in chat\n/quit\tDisconnect"))
}

func (s *Server) quit(c *Client, args []string) {
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitCurrent(c)

	c.msg("See you later :)")
	c.conn.Close()
}

func (s *Server) quitCurrent(c *Client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room.", c.nick))
	}
}
