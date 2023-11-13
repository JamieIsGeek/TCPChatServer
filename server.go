package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func (s *server) run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CmdNick:
			s.nick(cmd.client, cmd.args)
		case CmdMsg:
			s.msg(cmd.client, cmd.args)
		case CmdJoin:
			s.join(cmd.client, cmd.args)
		case CmdQuit:
			s.quit(cmd.client)
		case CmdRooms:
			s.listRooms(cmd.client)
		}
	}
}

func (s *server) newClient(conn net.Conn) {
	log.Printf("New client connected: %s", conn.RemoteAddr().String())

	c := &client{
		conn:     conn,
		nick:     "Anonymous",
		commands: s.commands,
	}

	c.readInput()
}

func newServer() *server {
	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("Your nick is now: %s", c.nick))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]

	r, ok := s.rooms[roomName]

	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c
	s.quitCurrentRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("You must join a room first!"))
		return
	}

	c.room.broadcast(c, c.nick+" > "+strings.Join(args[1:len(args)], " "))
}

func (s *server) quit(c *client) {
	s.quitCurrentRoom(c)
	c.msg("Sad to see you go!")
	c.conn.Close()
}

func (s *server) listRooms(c *client) {
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}

	c.msg(fmt.Sprintf("Available rooms: %s", strings.Join(rooms, ", ")))
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left", c.nick))
	}
}
