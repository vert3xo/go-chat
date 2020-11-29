package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_HELP
	CMD_QUIT
)

type Command struct {
	id     commandID
	client *Client
	args   []string
}
