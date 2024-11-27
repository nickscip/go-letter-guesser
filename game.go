package main

type Game struct {
        secretWord string
        players    *Players
        isGameOn   bool
}

func (g *Game) run() {
        g.players.defender.send <- []byte("Please enter a secret word.")
}
        
