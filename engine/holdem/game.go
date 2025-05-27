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

// Player Actions

// Call - current player calls
func (g *Game) Call() bool {
	player := g.GetCurrentPlayer()
	callAmount := g.CurrentBet - player.GetBet()

	if callAmount > 0 {
		actualAmount := min(callAmount, player.GetChips())
		player.Bet(actualAmount)
		g.Pot += actualAmount
	}

	g.ActionsThisRound++
	g.nextPlayer()
	return g.checkAndAdvanceIfBettingComplete()
}

// Raise - current player raises
func (g *Game) Raise(raiseAmount int) bool {
	player := g.GetCurrentPlayer()
	totalBet := g.CurrentBet + raiseAmount
	betNeeded := totalBet - player.GetBet()

	actualAmount := min(betNeeded, player.GetChips())
	player.Bet(actualAmount)
	g.Pot += actualAmount

	if player.GetBet() > g.CurrentBet {
		g.CurrentBet = player.GetBet()
	}

	g.ActionsThisRound++
	g.nextPlayer()
	return g.checkAndAdvanceIfBettingComplete()
}

// Check - current player checks
func (g *Game) Check() bool {
	g.ActionsThisRound++
	g.nextPlayer()
	return g.checkAndAdvanceIfBettingComplete()
}

// Fold - current player folds
func (g *Game) Fold() bool {
	player := g.GetCurrentPlayer()
	player.Fold()

	g.ActionsThisRound++
	g.nextPlayer()

	// Check if only one player remains
	activePlayers := g.GetActivePlayers()
	if len(activePlayers) <= 1 {
		g.endHandEarly()
		return true
	}

	return g.checkAndAdvanceIfBettingComplete()
}

// checkAndAdvanceIfBettingComplete checks if betting round is complete and advances if so
func (g *Game) checkAndAdvanceIfBettingComplete() bool {
	if g.IsBettingRoundComplete() {
		g.advanceToNextPhase()
		return true // Betting round was completed
	}
	return false
}

// advanceToNextPhase moves to next phase or ends hand
func (g *Game) advanceToNextPhase() bool {
	// Reset bets and action count for next round
	for _, player := range g.Players {
		player.ResetBet()
	}
	g.CurrentBet = 0
	g.ActionsThisRound = 0

	switch g.CurrentPhase {
	case PhasePreflop:
		g.CurrentPhase = PhaseFlop
		g.dealFlop()
		g.setFirstPlayerPostFlop()
		return false
	case PhaseFlop:
		g.CurrentPhase = PhaseTurn
		g.dealTurn()
		g.setFirstPlayerPostFlop()
		return false
	case PhaseTurn:
		g.CurrentPhase = PhaseRiver
		g.dealRiver()
		g.setFirstPlayerPostFlop()
		return false
	case PhaseRiver:
		g.CurrentPhase = PhaseShowdown
		g.showdown()
		return true // Hand is over
	case PhaseShowdown:
		return true // Already at showdown
	}
	return false
}

// endHandEarly ends hand when only one player remains
func (g *Game) endHandEarly() {
	activePlayers := g.GetActivePlayers()
	if len(activePlayers) == 1 {
		winner := activePlayers[0]
		winner.GrandChips(g.Pot)
		g.CurrentPhase = PhaseShowdown
	}
}

// showdown handles the showdown phase
func (g *Game) showdown() {
	activePlayers := g.GetActivePlayers()
	if len(activePlayers) == 0 {
		return
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

	// Find all winners (in case of tie)
	winners := lo.FilterMap(results, func(result playerResult, _ int) (IPlayer, bool) {
		if result.hand.Value == bestValue {
			return result.player, true
		}
		return nil, false
	})

	// Distribute pot among winners
	potShare := g.Pot / len(winners)
	remainder := g.Pot % len(winners)

	for i, winner := range winners {
		share := potShare
		if i == 0 {
			share += remainder // Give remainder to first winner
		}
		winner.GrandChips(share)
	}

	g.CurrentPhase = PhaseShowdown
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

// GetHandResult returns the hand result for a player (only valid after river)
func (g *Game) GetHandResult(player IPlayer) *HandResult {
	if g.CurrentPhase < PhaseRiver {
		return nil
	}
	return EvaluatePlayerHand(player, g.CommunityCards)
}

// NextPhase advances to next phase
func (g *Game) NextPhase() {
	// Reset bets for next round
	for _, player := range g.Players {
		player.ResetBet()
	}
	g.CurrentBet = 0

	switch g.CurrentPhase {
	case PhasePreflop:
		g.CurrentPhase = PhaseFlop
		g.dealFlop()
	case PhaseFlop:
		g.CurrentPhase = PhaseTurn
		g.dealTurn()
	case PhaseTurn:
		g.CurrentPhase = PhaseRiver
		g.dealRiver()
	case PhaseRiver:
		g.CurrentPhase = PhaseShowdown
	}

	// Set first player to act post-flop (small blind position)
	if g.CurrentPhase != PhaseShowdown {
		g.setFirstPlayerPostFlop()
	}
}

func (g *Game) setFirstPlayerPostFlop() {
	g.CurrentPlayerPosition = (g.DealerPosition + 1) % len(g.Players)

	// Find first active player
	original := g.CurrentPlayerPosition
	for g.Players[g.CurrentPlayerPosition].IsFolded() {
		g.CurrentPlayerPosition = (g.CurrentPlayerPosition + 1) % len(g.Players)
		if g.CurrentPlayerPosition == original {
			break
		}
	}
}

func (g *Game) dealFlop() {
	// Burn one card
	if len(g.Deck) > 0 {
		g.Deck = g.Deck[1:]
	}

	// Deal 3 cards
	for i := 0; i < 3 && len(g.Deck) > 0; i++ {
		g.CommunityCards.Append(g.Deck[0])
		g.Deck = g.Deck[1:]
	}
}

func (g *Game) dealTurn() {
	// Burn one card
	if len(g.Deck) > 0 {
		g.Deck = g.Deck[1:]
	}

	// Deal 1 card
	if len(g.Deck) > 0 {
		g.CommunityCards.Append(g.Deck[0])
		g.Deck = g.Deck[1:]
	}
}

func (g *Game) dealRiver() {
	// Burn one card
	if len(g.Deck) > 0 {
		g.Deck = g.Deck[1:]
	}

	// Deal 1 card
	if len(g.Deck) > 0 {
		g.CommunityCards.Append(g.Deck[0])
		g.Deck = g.Deck[1:]
	}
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

// IsBettingRoundComplete checks if betting round is done
func (g *Game) IsBettingRoundComplete() bool {
	activePlayers := g.GetActivePlayers()
	if len(activePlayers) <= 1 {
		return true
	}

	// In preflop, we need special handling because of blinds
	if g.CurrentPhase == PhasePreflop {
		// Everyone needs to match the current bet or be all-in
		// AND we need at least as many actions as players
		if g.ActionsThisRound < len(activePlayers) {
			return false
		}

		for _, player := range activePlayers {
			if player.GetBet() != g.CurrentBet && player.GetChips() > 0 {
				return false
			}
		}
		return true
	}

	// Post-flop: everyone needs to act at least once
	// and all have same bet (usually 0) or are all-in
	if g.ActionsThisRound < len(activePlayers) {
		return false
	}

	for _, player := range activePlayers {
		if player.GetBet() != g.CurrentBet && player.GetChips() > 0 {
			return false
		}
	}

	return true
}

// NextHand prepares for next hand
func (g *Game) NextHand() {
	g.DealerPosition = (g.DealerPosition + 1) % len(g.Players)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *Game) nextPlayer() {
	original := g.CurrentPlayerPosition
	for {
		g.CurrentPlayerPosition = (g.CurrentPlayerPosition + 1) % len(g.Players)
		player := g.Players[g.CurrentPlayerPosition]

		// Found active player or completed full circle
		if !player.IsFolded() || g.CurrentPlayerPosition == original {
			break
		}
	}
}
