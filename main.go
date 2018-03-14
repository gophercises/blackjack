package main

import (
	"strings"

	"github.com/gophercises/deck"
)

type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			// ace is currently worth 1, and we are changing it to be worth 11
			// 11 - 1 = 10
			return minScore + 10
		}
	}
	return minScore
}

func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
		score += min(int(c.Rank), 10)
	}
	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	// cards := deck.New(deck.Deck(3), deck.Shuffle)
	// var card deck.Card
	// var player, dealer Hand
	// for i := 0; i < 2; i++ {
	// 	for _, hand := range []*Hand{&player, &dealer} {
	// 		card, cards = draw(cards)
	// 		*hand = append(*hand, card)
	// 	}
	// }
	// var input string
	// for input != "s" {
	// 	fmt.Println("Player:", player)
	// 	fmt.Println("Dealer:", dealer.DealerString())
	// 	fmt.Println("What will you do? (h)it, (s)tand")
	// 	fmt.Scanf("%s\n", &input)
	// 	switch input {
	// 	case "h":
	// 		card, cards = draw(cards)
	// 		player = append(player, card)
	// 	}
	// }
	// // If dealer score <= 16, we hit
	// // If dealer has a soft 17, then we hit.
	// for dealer.Score() <= 16 || (dealer.Score() == 17 && dealer.MinScore() != 17) {
	// 	card, cards = draw(cards)
	// 	dealer = append(dealer, card)
	// }
	// pScore, dScore := player.Score(), dealer.Score()
	// fmt.Println("==FINAL HANDS==")
	// fmt.Println("Player:", player, "\nScore:", pScore)
	// fmt.Println("Dealer:", dealer, "\nScore:", dScore)
	// switch {
	// case pScore > 21:
	// 	fmt.Println("You busted")
	// case dScore > 21:
	// 	fmt.Println("Dealer busted")
	// case pScore > dScore:
	// 	fmt.Println("You win!")
	// case dScore > pScore:
	// 	fmt.Println("You lose")
	// case dScore == pScore:
	// 	fmt.Println("Draw")
	// }
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

type State int8

const (
	StatePlayerTurn State = iota
	StateDealerTurn
	StateHandOver
)

type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it isn't currently any player's turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}
