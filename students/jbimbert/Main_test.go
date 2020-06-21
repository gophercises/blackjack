package main

import (
	"Gophercizes/deck/students/jbimbert/deck"
	"bufio"
	"strings"
	"testing"
)

// Test a simple score
func TestScore10(t *testing.T) {
	dk := []deck.Card{{Rank: deck.V10}}
	sc := score(dk)
	if len(sc) != 1 {
		t.Errorf("Bad score length: wanted 1 found: %d", len(sc))
	}
	if sc[0] != 10 {
		t.Errorf("Bad score: wanted: 10 found: %d", sc[0])
	}
}

// Test ace, min and max
func TestScoreAce(t *testing.T) {
	dk := []deck.Card{{Rank: deck.VA}}
	sc := score(dk)
	if len(sc) != 2 {
		t.Errorf("Bad score length: wanted 2 found: %d", len(sc))
	}
	if mins(sc) != 1 {
		t.Errorf("Bad score: wanted: 1 found: %d", mins(sc))
	}
	if maxs(sc) != 11 {
		t.Errorf("Bad score: wanted: 11 found: %d", maxs(sc))
	}
}

// Test aces and contain
func TestScoreAces(t *testing.T) {
	dk := []deck.Card{{Rank: deck.VA}, {Rank: deck.VA}}
	sc := score(dk)
	if len(sc) != 3 {
		t.Errorf("Bad score length: wanted 3 found: %d", len(sc))
	}
	if mins(sc) != 2 {
		t.Errorf("Bad score: wanted: 1 found: %d", mins(sc))
	}
	if maxs(sc) != 12 { // Here we should skip 22
		t.Errorf("Bad score: wanted: 22 found: %d", maxs(sc))
	}
	if !contain(sc, 12) {
		t.Errorf("Scores should contain 12")
	}
}

func compareRanks(wanted, found deck.Rank, t *testing.T) {
	if wanted != found {
		t.Errorf("Bad rank. Wanted: %d Found: %d", wanted, found)
	}
}

func TestDistribute(t *testing.T) {
	decks = []deck.Card{{Rank: deck.V5}, {Rank: deck.V8}, {Rank: deck.V6}, {Rank: deck.V9}, {Rank: deck.VA},
		{Rank: deck.V3}}
	players := []Player{{name: "player1"}, {name: dealer}}
	distribute(&players)
	for _, p := range players {
		if len(p.cards) != 2 {
			t.Errorf("Player %s wrong number of cards. Wanted 2 Found %d", p.name, len(p.cards))
		}
	}
	compareRanks(deck.V5, players[0].cards[0].Rank, t)
	compareRanks(deck.V6, players[0].cards[1].Rank, t)
	compareRanks(deck.V8, players[1].cards[0].Rank, t)
	compareRanks(deck.V9, players[1].cards[1].Rank, t)
}

func TestIsSoft17(t *testing.T) {
	sc := []deck.Card{{Rank: deck.V5}, {Rank: deck.V6}, {Rank: deck.V6}} // 17 but no Ace => False
	if isSoft17(sc) {
		t.Errorf("Expected false for %v", sc)
	}
	sc = []deck.Card{{Rank: deck.VA}, {Rank: deck.V6}} // 7 17 and Ace = 11 => True
	if !isSoft17(sc) {
		t.Errorf("Expected true for %v", sc)
	}
	sc = []deck.Card{{Rank: deck.V9}, {Rank: deck.V7}, {Rank: deck.VA}} // 17 but Ace = 1 => False
	if isSoft17(sc) {
		t.Errorf("Expected false for %v", sc)
	}
	sc = []deck.Card{{Rank: deck.VA}, {Rank: deck.V5}, {Rank: deck.VA}} // 7 17 27 and Aces = 1 11 => True
	if !isSoft17(sc) {
		t.Errorf("Expected true for %v", sc)
	}
}

func TestGame1(t *testing.T) {
	cur = 0
	decks = []deck.Card{{Rank: deck.V5}, {Rank: deck.V8}, {Rank: deck.V6}, {Rank: deck.V9}, {Rank: deck.VA}}
	players := []Player{{name: "player1"}, {name: dealer}}
	distribute(&players)

	strategy := "hs"
	reader := bufio.NewReader(strings.NewReader(strategy))

	winner := play(players, decks, reader)

	if (*winner).name != players[1].name {
		t.Error("Bad winner. Expected", players[1], "Got", winner)
	}
}
