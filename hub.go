package main

import (
	"fmt"
)

type Players struct {
	defender *Client
	guesser1 *Client
	guesser2 *Client
}

func (p *Players) isFull() bool {
	return p.defender != nil && p.guesser1 != nil && p.guesser2 != nil
}

type MessageType string

const (
	MessageBroadcast MessageType = "broadcast"
	MessagePrivate   MessageType = "private"
)

type Message struct {
	sender      *Client
	recipient   *Client
	message     []byte
	messageType MessageType
}

type Hub struct {
	clients     map[*Client]bool
	register    chan *Client
	broadcast   chan Message
	unregister  chan *Client
	players     *Players
	playerQueue []*Client
	game        *Game
	startGame   chan bool
}

func newHub() *Hub {
	return &Hub{
		clients:     make(map[*Client]bool),
		register:    make(chan *Client),
		broadcast:   make(chan Message, 1),
		unregister:  make(chan *Client),
		players:     &Players{},
		playerQueue: []*Client{},
		game:        &Game{isGameOn: false, gameMessages: make(chan Message, 1)},
		startGame:   make(chan bool, 1),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("Client registered")
			h.assignRole(client)

		case m := <-h.broadcast:
			if m.messageType == MessageBroadcast {
				// Broadcast to all clients except the sender
				for client := range h.clients {
					if client != m.sender {
						select {
						case client.send <- m.message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			} else if m.messageType == MessagePrivate && m.recipient != nil {
				// Send only to the specified recipient
				select {
				case m.recipient.send <- m.message:
				default:
					close(m.recipient.send)
					delete(h.clients, m.recipient)
				}
			}
		case <-h.startGame:
			h.game.players = h.players
			h.game.isGameOn = true
			go h.game.run()

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.removePlayer(client)
		}
	}
}

func (h *Hub) assignRole(client *Client) {
	if h.players.isFull() {
		h.playerQueue = append(h.playerQueue, client)
		fmt.Println("Added to player queue")
		client.send <- []byte("You are in the queue. Please wait.")
		return
	}

	if h.players.defender == nil {
		h.players.defender = client
		fmt.Println("Assigned defender role")
		client.send <- []byte("You are the Defender")
	} else if h.players.guesser1 == nil {
		h.players.guesser1 = client
		fmt.Println("Assigned guesser1 role")
		client.send <- []byte("You are Guesser 1")
	} else {
		h.players.guesser2 = client
		fmt.Println("Assigned guesser2 role")
		client.send <- []byte("You are Guesser 2")
	}

	if h.players.isFull() && !h.game.isGameOn {
                h.broadcast <- Message{messageType: MessageBroadcast, message: []byte("All roles have been assigned. The game will begin shortly.")}
		h.startGame <- true
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

	if len(h.playerQueue) > 0 {
		fmt.Println("Reassigning roles from the queue")
		nextClient := h.playerQueue[0]
		nextClient.send <- []byte("A player has left. You will be assigned their role in the game.")
		h.playerQueue = h.playerQueue[1:]
		h.assignRole(nextClient)
	}
}
