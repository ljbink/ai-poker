package holdem

import (
	"github.com/ljbink/ai-poker/engine/poker"
	"github.com/samber/lo"
)

type GamePhase int

const (
	PhasePreflop GamePhase = iota
	PhaseFlop
	PhaseTurn
	PhaseRiver
	PhaseShowdown
)

type Game struct {
	Players               []IPlayer   // Players in the game
	Deck                  poker.Cards // Deck of cards
	Pot                   int         // Current pot
	CurrentPhase          GamePhase   // Current phase of the game
	CurrentPlayerPosition int         // Current player
	DealerPosition        int         // Position of the dealer
	SmallBlind            int         // Small blind amount
	BigBlind              int         // Big blind amount
	CurrentBet            int         // Current bet amount
	CommunityCards        poker.Cards // Community cards
	ActionsThisRound      int         // Actions taken this round
}

// NewGame creates a new game
func NewGame(playerNames []string, smallBlind, bigBlind int) *Game {
	var players = make([]IPlayer, len(playerNames))
	for i, name := range playerNames {
		players[i] = NewPlayer(i, name, 1000) // Starting with 1000 chips
	}
	return &Game{
		Players:               []IPlayer(players),
		SmallBlind:            smallBlind,
		BigBlind:              bigBlind,
		DealerPosition:        0,
		CurrentPlayerPosition: 0,
		CurrentPhase:          PhasePreflop,
	}
}

// GetCurrentPlayer returns the current player
func (g *Game) GetCurrentPlayer() IPlayer {
	return g.Players[g.CurrentPlayerPosition]
}

// GetDealer returns the dealer
func (g *Game) GetDealer() IPlayer {
	return g.Players[g.DealerPosition]
}

// StartHand starts a new hand
func (g *Game) StartHand() {
	// Reset for new hand
	g.resetForNewHand()

	// Create and shuffle deck
	g.createDeck()

	// Deal hole cards
	g.dealHoleCards()

	// Post blinds
	g.postBlinds()

	// Set first player to act
	g.setFirstPlayerToAct()
}

func (g *Game) resetForNewHand() {
	g.Pot = 0
	g.CurrentBet = 0
	g.CurrentPhase = PhasePreflop
	g.CommunityCards = poker.Cards{}
	g.ActionsThisRound = 0

	// Reset all players
	for _, player := range g.Players {
		player.ResetForNewHand()
	}
}

func (g *Game) createDeck() {
	g.Deck = poker.NewDeckCards()

	// Remove jokers for Hold'em
	joker1 := poker.NewCard(poker.SuitNone, poker.RankJoker)
	joker2 := poker.NewCard(poker.SuitNone, poker.RankColoredJoker)
	g.Deck.Remove(joker1, joker2)

	g.Deck.Shuffle()
}

func (g *Game) dealHoleCards() {
	// Deal 2 cards to each player
	for i := 0; i < 2; i++ {
		for _, player := range g.Players {
			if len(g.Deck) > 0 {
				player.DealCard(g.Deck[0])
				g.Deck = g.Deck[1:]
			}
		}
	}
}

func (g *Game) postBlinds() {
	numPlayers := len(g.Players)

	// Small blind position
	sbPos := (g.DealerPosition + 1) % numPlayers
	sbPlayer := g.Players[sbPos]

	// Big blind position
	bbPos := (g.DealerPosition + 2) % numPlayers
	bbPlayer := g.Players[bbPos]

	// Post small blind
	sbAmount := min(g.SmallBlind, sbPlayer.GetChips())
	sbPlayer.Bet(sbAmount)
	g.Pot += sbAmount

	// Post big blind
	bbAmount := min(g.BigBlind, bbPlayer.GetChips())
	bbPlayer.Bet(bbAmount)
	g.Pot += bbAmount
	g.CurrentBet = bbAmount
}

func (g *Game) setFirstPlayerToAct() {
	numPlayers := len(g.Players)
	// First to act is after big blind (UTG)
	g.CurrentPlayerPosition = (g.DealerPosition + 3) % numPlayers

	// If heads-up, small blind acts first preflop
	if numPlayers == 2 {
		g.CurrentPlayerPosition = (g.DealerPosition + 1) % numPlayers
	}
}

// GetWinners returns the winners of the current hand (only valid after showdown)
func (g *Game) GetWinners() []IPlayer {
	if g.CurrentPhase != PhaseShowdown {
		return nil
	}

	activePlayers := g.GetActivePlayers()
	if len(activePlayers) <= 1 {
		return activePlayers
	}

	// Evaluate all active hands
	type playerResult struct {
		player IPlayer
		hand   *HandResult
	}

	results := lo.Map(activePlayers, func(player IPlayer, _ int) playerResult {
		hand := EvaluatePlayerHand(player, g.CommunityCards)
		return playerResult{player, hand}
	})

	// Find the best hand value
	bestValue := lo.MaxBy(results, func(a, b playerResult) bool {
		return a.hand.Value > b.hand.Value
	}).hand.Value

	// Return all winners
	return lo.FilterMap(results, func(result playerResult, _ int) (IPlayer, bool) {
		if result.hand.Value == bestValue {
			return result.player, true
		}
		return nil, false
	})
}

// IsHandComplete returns true if the hand is finished
func (g *Game) IsHandComplete() bool {
	return g.CurrentPhase == PhaseShowdown
}

// Utility functions

// GetActivePlayers returns non-folded players
func (g *Game) GetActivePlayers() []IPlayer {
	var active []IPlayer
	for _, player := range g.Players {
		if !player.IsFolded() {
			active = append(active, player)
		}
	}
	return active
}

// AI Helper Methods

// GetPlayerPosition returns the position of a specific player
func (g *Game) GetPlayerPosition(player IPlayer) int {
	for i, p := range g.Players {
		if p.GetID() == player.GetID() {
			return i
		}
	}
	return -1
}

// GetNumActivePlayers returns the number of active (non-folded) players
func (g *Game) GetNumActivePlayers() int {
	activePlayers := g.GetActivePlayers()
	return len(activePlayers)
}

// GetTotalPlayers returns the total number of players
func (g *Game) GetTotalPlayers() int {
	return len(g.Players)
}

// CalculatePotOdds calculates pot odds for a given player
func (g *Game) CalculatePotOdds(player IPlayer) float64 {
	callAmount := g.CurrentBet - player.GetBet()
	if callAmount <= 0 {
		return 0.0
	}
	return float64(callAmount) / float64(g.Pot+callAmount)
}

// IsPlayerTurn checks if it's a specific player's turn
func (g *Game) IsPlayerTurn(player IPlayer) bool {
	currentPlayer := g.GetCurrentPlayer()
	return currentPlayer.GetID() == player.GetID()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
