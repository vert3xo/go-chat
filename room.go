package main

import "net"

type Room struct {
	name    string
	members map[net.Addr]*Client
}

func (r *Room) broadcast(sender *Client, message string) {
	for addr, member := range r.members {
		if addr != sender.conn.RemoteAddr() {
			member.msg(message)
		}
	}
}
