package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	nick     string
	room     *Room
	commands chan<- Command
	prompt   string
}

func (c *Client) readInput() {
	for {
		c.sendPrompt(fmt.Sprint(c.prompt))
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		if msg == "" || msg[0] == ' ' {
			continue
		}

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
			case "/help":
				c.commands <- Command{
					id:     CMD_HELP,
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
	c.conn.Write([]byte("\rERR: " + err.Error() + "\n"))
	c.sendPrompt(c.prompt)
}

func (c *Client) warn(warn string) {
	c.conn.Write([]byte("\rWARN: " + warn + "\n"))
	c.sendPrompt(c.prompt)
}

func (c *Client) msg(msg string) {
	c.conn.Write([]byte("\r" + msg + "\n"))
	c.sendPrompt(c.prompt)
}

func (c *Client) sendPrompt(prompt string) {
	c.conn.Write([]byte(c.nick + " " + prompt))
}

func genRandomName(n int) string {
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var s strings.Builder

	for i := 0; i < n; i++ {
		s.WriteByte(letters[rand.Intn(len(letters))])
	}
	return s.String()
}
