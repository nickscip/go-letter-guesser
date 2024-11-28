package main

type Game struct {
	secretWord   string
        guessWord    string
	players      *Players
	isGameOn     bool
	gameMessages chan Message
}

func (g *Game) run() {
	g.players.defender.currentTask = "enter_secret_word"
	g.players.defender.send <- []byte("Please enter a secret word.")

	for {
		select {
		case message := <-g.gameMessages:
			if message.sender == g.players.defender {
				if g.players.defender.currentTask == "enter_secret_word" {
					g.secretWord = string(message.message)
                                        g.players.defender.currentTask = ""
                                        g.players.defender.hub.broadcast <- Message{message: []byte("Secret word has been set."), messageType: MessageBroadcast, sender: g.players.defender}
                                        g.players.defender.hub.broadcast <- Message{message: []byte("The first letter of the secret word is: " + string(g.secretWord[0])), messageType: MessageBroadcast, sender: g.players.defender}
                                        g.players.guesser1.currentTask = "enter_guess_word"
                                        g.players.guesser1.send <- []byte("Please enter a guess word.")
				}
			}
                        if message.sender == g.players.guesser1 {
                                if g.players.guesser1.currentTask == "enter_guess_word" {
                                        if string(message.message)[0] != g.secretWord[0] {
                                                g.players.guesser1.send <- []byte("Please enter a guess word starting with the letter " + string(g.secretWord[0]))
                                        } else {
                                                g.guessWord = string(message.message)
                                                g.players.guesser1.currentTask = ""
                                                g.players.guesser1.hub.broadcast <- Message{message: []byte("Guess word has been set. The clue is: "), messageType: MessageBroadcast, sender: g.players.guesser1}
                                                g.players.guesser1.send <- []byte("Please enter a clue.")
                                        }
                                }
                        }
                default:
                        break
		}
	}
}
