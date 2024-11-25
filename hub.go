package main

import (
	"fmt"
)

type Hub struct {
	clients   map[*Client]bool
	register  chan *Client
	broadcast chan []byte
        unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		register:  make(chan *Client),
		broadcast: make(chan []byte),
                unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("Client registered")
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
                case client := <-h.unregister:
                        if _, ok := h.clients[client]; ok {
                                delete(h.clients, client)
                                close(client.send)
                        }
		}
	}
}
