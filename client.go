package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CmdNick,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CmdJoin,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CmdRooms,
				client: c,
				args:   args,
			}
		case "/msg":
			c.commands <- command{
				id:     CmdMsg,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CmdQuit,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("Unknown Command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte(msg + "\n"))

}
