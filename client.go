package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	nick     string
	room     *Room
	commands chan<- Command
}

func (c *Client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		if cmd[0] == '/' {
			switch cmd {
			case "/nick":
				c.commands <- Command{
					id:     CMD_NICK,
					client: c,
					args:   args,
				}
			case "/join":
				c.commands <- Command{
					id:     CMD_JOIN,
					client: c,
					args:   args,
				}
			case "/rooms":
				c.commands <- Command{
					id:     CMD_ROOMS,
					client: c,
					args:   args,
				}
			case "/msg":
				c.commands <- Command{
					id:     CMD_MSG,
					client: c,
					args:   args,
				}
			case "/quit":
				c.commands <- Command{
					id:     CMD_QUIT,
					client: c,
					args:   args,
				}
			default:
				c.err(fmt.Errorf("Unknown command: %s", cmd))
			}
		} else {
			c.commands <- Command{
				id:     CMD_MSG,
				client: c,
				args:   strings.Split("m "+msg, " "),
			}
		}
	}
}

func (c *Client) err(err error) {
	c.conn.Write([]byte("ERR: " + err.Error() + "\n"))
}

func (c *Client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
