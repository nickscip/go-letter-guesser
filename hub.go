package main

import (
	"fmt"
)

type Players struct {
	defender *Client
	guesser1 *Client
	guesser2 *Client
}

type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	broadcast   chan []byte
	unregister  chan *Client
	players     *Players
	playerQueue []*Client
}

func newHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		broadcast:   make(chan []byte),
		unregister:  make(chan *Client),
		players:     &Players{},
		playerQueue: []*Client{},
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("Client registered")
			h.assignRole(client)
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

func (h *Hub) assignRole(client *Client) {
	if h.players.defender == nil {
		h.players.defender = client
		fmt.Println("Assigned defender role")
		client.send <- []byte("You are the Defender")
	} else if h.players.guesser1 == nil {
		h.players.guesser1 = client
		fmt.Println("Assigned guesser1 role")
		client.send <- []byte("You are Guesser 1")
	} else if h.players.guesser2 == nil {
		h.players.guesser2 = client
		fmt.Println("Assigned guesser2 role")
		client.send <- []byte("You are Guesser 2")
	} else {
		// Add to the player queue if all roles are filled
		h.playerQueue = append(h.playerQueue, client)
		fmt.Println("Added to player queue")
		client.send <- []byte("You are in the queue. Please wait.")
	}
}

func (h *Hub) removePlayer(client *Client) {
	if h.players.defender == client {
		h.players.defender = nil
		fmt.Println("Defender role is now vacant")
	} else if h.players.guesser1 == client {
		h.players.guesser1 = nil
		fmt.Println("Guesser1 role is now vacant")
	} else if h.players.guesser2 == client {
		h.players.guesser2 = nil
		fmt.Println("Guesser2 role is now vacant")
	}

	// Reassign roles from the queue if available
	if len(h.playerQueue) > 0 {
		nextClient := h.playerQueue[0]
		h.playerQueue = h.playerQueue[1:]
		h.assignRole(nextClient)
	}
}
