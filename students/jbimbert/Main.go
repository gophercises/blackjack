package main

import (
	"Gophercizes/deck/students/jbimbert/deck"
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	dealer = "dealer" // name of the dealer
	BUST   = 21       // threshold value for a blackjack
)

type Player struct {
	name  string
	cards []deck.Card
}

func (p Player) String() string {
	s := score(p.cards)
	return fmt.Sprintf("%s : %s => %v\n", p.name, p.cards, s)
}

var decks []deck.Card
var cur int = 0 // current position (the card to give)

// give cards to the players from the given deck
func distribute(players *[]Player) {
	for i := 0; i < 2; i++ {
		for idx := range *players {
			p := &(*players)[idx]
			hit(p, false)
		}
	}
}

// give a card to player p and display it to screen according to show value
func hit(p *Player, show bool) {
	p.cards = append(p.cards, decks[cur])
	if show {
		fmt.Println(p.name, "got a", decks[cur])
	}
	cur++
}

// display current cards and associated scores
func board(players []Player) {
	for _, p := range players {
		if p.name == dealer {
			fmt.Printf("%s : %s\n", p.name, p.cards[0])
		} else {
			fmt.Printf("%s\n", p)
		}
	}
}

// return the score of a card. If aceId11 is true, then the Ace value is 11 instead of 1
func rank(c deck.Card, aceIs11 bool) int {
	if c.Rank == deck.VA && aceIs11 {
		return 11
	}
	return min(10, int(c.Rank))
}

// compute the different scores (because of Ace) of a slice of cards
func score(cards []deck.Card) []int {
	r := make([]int, 1)
	for _, c := range cards {
		if c.Rank == deck.VA {
			rr := make([]int, 0)
			for i, v := range r {
				r[i] += rank(c, false)
				rr = append(rr, v+rank(c, true))
			}
			r = append(r, rr...)
		} else {
			for i := range r {
				r[i] += rank(c, false)
			}
		}
	}
	//suppress duplicates if any
	d := make([]int, 0)
	for _, v := range r {
		if !contain(d, v) {
			d = append(d, v)
		}
	}
	return d
}

// return true if cs contains v, false otherwise
func contain(cs []int, v int) bool {
	for _, c := range cs {
		if c == v {
			return true
		}
	}
	return false
}

func min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

// return the mins of cs
func mins(cs []int) int {
	m := cs[0]
	for _, c := range cs {
		m = min(m, c)
	}
	return m
}

// return the maxs of cs which is less than BUST (or -1 if none found)
func maxs(cs []int) int {
	if cs[0] > BUST {
		return -1
	}
	m := cs[0]
	for _, c := range cs {
		if c < BUST && m < c {
			m = c
		}
	}
	return m
}

// Get the "dealer"
func getDealer(ps *[]Player) *Player {
	for idx, p := range *ps {
		if p.name == dealer {
			return &(*ps)[idx]
		}
	}
	return nil
}

// Get the player
func getPlayer(ps *[]Player) *Player {
	for idx, p := range *ps {
		if p.name == dealer {
			continue
		}
		return &(*ps)[idx]
	}
	return nil
}

// Read the 'h' or 's' from the reader (allow testing with a file or a string instead of os.Stdin)
func readRune(in io.Reader) string {
	reader := bufio.NewReader(in)
	fmt.Println("\nHit (h) or Stand (s)?")
again:
	r, _, _ := reader.ReadRune()
	if r != 'h' && r != 's' {
		fmt.Println("Bad input. Please enter \"s\" or \"h\"")
		goto again
	}
	return string(r)
}

// A “soft 17” is a score of 17 in which 11 of the points come from an Ace card.
func isSoft17(cards []deck.Card) bool {
	sc := score(cards)
	return contain(sc, 17) && mins(sc) <= 7
}

// Game strategy of the dealer
func dealerStrategy(dealer *Player) int {
	sc := score(dealer.cards)
	for maxs(sc) <= 16 || isSoft17(dealer.cards) {
		if mins(sc) > BUST {
			return -1
		}
		hit(dealer, true)
		sc = score(dealer.cards)
	}
	return maxs(sc)
}

// Player plays as he wants
func playerStrategy(player *Player, players []Player, reader io.Reader) int {
	input := readRune(reader)
	sc := score(player.cards)
	for input == "h" {
		hit(player, true)
		board(players)
		sc = score(player.cards)
		m := mins(sc)
		switch {
		case m == BUST: // Shortcut : player wins
			return BUST
		case m > BUST:
			return -1
		}
		input = readRune(reader)
	}
	return maxs(sc)
}

// Play the game
func play(players []Player, decks []deck.Card, reader io.Reader) *Player {
	player := getPlayer(&players)
	pScore := playerStrategy(player, players, reader)
	if pScore == BUST {
		return player
	}
	dealer := getDealer(&players)
	dealerScore := dealerStrategy(dealer)
	if pScore < dealerScore {
		return dealer
	}
	return player
}

func main() {
	decks = deck.NewDeck(deck.WithShuffle())
	players := []Player{{name: "player1"}, {name: dealer}}
	distribute(&players)
	board(players)

	winner := play(players, decks, os.Stdin)
	fmt.Println("Winner is", winner)

	fmt.Println(players)
}
