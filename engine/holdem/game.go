package holdem

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ljbink/ai-poker/engine/poker"
)

type GamePhase int

const (
	PhasePreflop GamePhase = iota
	PhaseFlop
	PhaseTurn
	PhaseRiver
	PhaseShowdown
)

type SystemActions struct {
	Preflop []Action
	Flop    []Action
	Turn    []Action
	River   []Action
}

type UserActions struct {
	Preflop []Action
	Flop    []Action
	Turn    []Action
	River   []Action
}

type IGame interface {
	GetSmallBlind() int
	GetBigBlind() int

	GetCurrentPhase() GamePhase
	SetCurrentPhase(phase GamePhase)

	GetCommunityCards() poker.Cards

	PlayerSit(player IPlayer, sit int) error
	PlayerLeave(player IPlayer) error
	GetPlayerByID(id int) (IPlayer, error)
	GetPlayerBySit(sit int) (IPlayer, error)
	GetPlayerSitByID(id int) (int, error)
	GetAllPlayers() []IPlayer

	DealHoleCards() error
	DealFlop() error
	DealTurn() error
	DealRiver() error
	ShuffleDeck()

	GetCurrentPlayer() IPlayer

	GetSystemActions() SystemActions
	GetUserActions() UserActions

	TakeAction(action Action) error
}

type Game struct {
	players        [10]IPlayer // Players in the game with sitting number
	deck           poker.Cards // Deck of cards
	communityCards poker.Cards // Community cards
	currentPhase   GamePhase   // Current phase of the game

	smallBlind int // Small blind amount
	bigBlind   int // Big blind amount

	systemActions SystemActions
	userActions   UserActions
}

func (g *Game) PlayerSit(player IPlayer, sit int) error {
	if player == nil {
		return fmt.Errorf("player is nil")
	}
	if sit < 0 || sit >= len(g.players) {
		return fmt.Errorf("invalid sit number: %d", sit)
	}
	if g.players[sit] != nil && g.players[sit].GetID() != player.GetID() {
		return fmt.Errorf("player already sitting at sit: %d", sit)
	}
	g.players[sit] = player
	return nil
}

func (g *Game) PlayerLeave(player IPlayer) error {
	if player == nil {
		return fmt.Errorf("player is nil")
	}
	for i, p := range g.players {
		if p == player {
			g.players[i] = nil
			break
		}
	}
	return nil
}

func (g *Game) GetSmallBlind() int {
	return g.smallBlind
}

func (g *Game) GetBigBlind() int {
	return g.bigBlind
}

func (g *Game) GetCurrentPhase() GamePhase {
	return g.currentPhase
}

func (g *Game) SetCurrentPhase(phase GamePhase) {
	oldPhase := g.currentPhase
	g.currentPhase = phase

	// Log system action for phase change (if it's actually changing)
	if oldPhase != phase {
		g.TakeSystemAction(Action{
			PlayerID: SystemPlayerID,
			Type:     ActionSystemPhaseChange,
			Amount:   int(phase), // Store the new phase as amount
		})
	}
}

func (g *Game) GetCommunityCards() poker.Cards {
	return g.communityCards
}

func (g *Game) GetPlayerByID(id int) (IPlayer, error) {
	for _, player := range g.players {
		if player != nil && player.GetID() == id {
			return player, nil
		}
	}
	return nil, fmt.Errorf("player with ID %d not found", id)
}

func (g *Game) GetPlayerBySit(sit int) (IPlayer, error) {
	if sit < 0 || sit >= len(g.players) {
		return nil, fmt.Errorf("invalid sit number: %d", sit)
	}
	if g.players[sit] == nil {
		return nil, fmt.Errorf("no player at sit %d", sit)
	}
	return g.players[sit], nil
}

func (g *Game) GetPlayerSitByID(id int) (int, error) {
	for i, player := range g.players {
		if player != nil && player.GetID() == id {
			return i, nil
		}
	}
	return -1, fmt.Errorf("player with ID %d not found", id)
}

func (g *Game) GetAllPlayers() []IPlayer {
	var players []IPlayer
	for _, player := range g.players {
		if player != nil {
			players = append(players, player)
		}
	}
	return players
}

func (g *Game) GetCurrentPlayer() IPlayer {
	// Find the first non-nil, non-folded player
	for _, player := range g.players {
		if player != nil && !player.IsFolded() {
			return player
		}
	}
	return nil
}

func (g *Game) GetSystemActions() SystemActions {
	return g.systemActions
}

func (g *Game) GetUserActions() UserActions {
	return g.userActions
}

func (g *Game) TakeAction(action Action) error {
	// Add action to the appropriate phase log in userActions
	switch g.currentPhase {
	case PhasePreflop:
		g.userActions.Preflop = append(g.userActions.Preflop, action)
	case PhaseFlop:
		g.userActions.Flop = append(g.userActions.Flop, action)
	case PhaseTurn:
		g.userActions.Turn = append(g.userActions.Turn, action)
	case PhaseRiver:
		g.userActions.River = append(g.userActions.River, action)
	default:
		return fmt.Errorf("invalid game phase: %d", g.currentPhase)
	}
	return nil
}

// TakeSystemAction logs system actions (like dealing cards, phase changes)
func (g *Game) TakeSystemAction(action Action) error {
	action.PlayerID = SystemPlayerID
	switch g.currentPhase {
	case PhasePreflop:
		g.systemActions.Preflop = append(g.systemActions.Preflop, action)
	case PhaseFlop:
		g.systemActions.Flop = append(g.systemActions.Flop, action)
	case PhaseTurn:
		g.systemActions.Turn = append(g.systemActions.Turn, action)
	case PhaseRiver:
		g.systemActions.River = append(g.systemActions.River, action)
	default:
		return fmt.Errorf("invalid game phase: %d", g.currentPhase)
	}
	return nil
}

// Card dealing methods

func (g *Game) ShuffleDeck() {
	// Shuffle existing deck using Fisher-Yates algorithm
	rand.Seed(time.Now().UnixNano())
	for i := len(g.deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		g.deck[i], g.deck[j] = g.deck[j], g.deck[i]
	}

	// Log system action for deck shuffle
	g.TakeSystemAction(Action{
		PlayerID: SystemPlayerID,
		Type:     ActionSystemShuffle,
		Amount:   0,
	})
}

// newStandardDeck creates a standard 52-card poker deck (no jokers)
func newStandardDeck() poker.Cards {
	suits := []poker.Suit{
		poker.SuitHeart,
		poker.SuitDiamond,
		poker.SuitClub,
		poker.SuitSpade,
	}
	ranks := []poker.Rank{
		poker.RankAce,
		poker.RankTwo,
		poker.RankThree,
		poker.RankFour,
		poker.RankFive,
		poker.RankSix,
		poker.RankSeven,
		poker.RankEight,
		poker.RankNine,
		poker.RankTen,
		poker.RankJack,
		poker.RankQueen,
		poker.RankKing,
	}
	cards := poker.Cards{}
	for _, suit := range suits {
		for _, rank := range ranks {
			cards.Append(poker.NewCard(suit, rank))
		}
	}
	return cards
}

// ResetAndShuffleDeck creates a fresh deck and shuffles it
func (g *Game) ResetAndShuffleDeck() {
	// Reset deck to standard 52 cards (no jokers)
	g.deck = newStandardDeck()
	g.ShuffleDeck()
}

func (g *Game) DealHoleCards() error {
	activePlayers := g.GetAllPlayers()
	if len(activePlayers) < 2 {
		return fmt.Errorf("need at least 2 players to deal cards")
	}

	// Reset and shuffle deck before dealing
	g.ResetAndShuffleDeck()

	// Clear existing cards from players
	for _, player := range activePlayers {
		player.ResetForNewHand()
	}

	// Deal 2 cards to each player
	cardIndex := 0
	for round := 0; round < 2; round++ {
		for _, player := range activePlayers {
			if !player.IsFolded() && cardIndex < len(g.deck) {
				player.DealCard(g.deck[cardIndex])
				cardIndex++
			}
		}
	}

	// Remove dealt cards from deck
	g.deck = g.deck[cardIndex:]

	// Log system action for dealing hole cards
	g.TakeSystemAction(Action{
		PlayerID: SystemPlayerID,
		Type:     ActionSystemDealHole,
		Amount:   len(activePlayers) * 2, // Number of cards dealt
	})

	return nil
}

func (g *Game) DealFlop() error {
	if len(g.deck) < 4 {
		return fmt.Errorf("not enough cards in deck for flop")
	}

	// Burn one card, then deal 3 community cards
	g.deck = g.deck[1:] // Burn card

	// Deal 3 cards to community
	for i := 0; i < 3; i++ {
		g.communityCards = append(g.communityCards, g.deck[i])
	}
	g.deck = g.deck[3:]

	g.currentPhase = PhaseFlop

	// Log system action for dealing flop
	g.TakeSystemAction(Action{
		PlayerID: SystemPlayerID,
		Type:     ActionSystemDealFlop,
		Amount:   3, // Number of community cards dealt
	})

	return nil
}

func (g *Game) DealTurn() error {
	if len(g.deck) < 2 {
		return fmt.Errorf("not enough cards in deck for turn")
	}

	// Burn one card, then deal 1 community card
	g.deck = g.deck[1:] // Burn card

	g.communityCards = append(g.communityCards, g.deck[0])
	g.deck = g.deck[1:]

	g.currentPhase = PhaseTurn

	// Log system action for dealing turn
	g.TakeSystemAction(Action{
		PlayerID: SystemPlayerID,
		Type:     ActionSystemDealTurn,
		Amount:   1, // Number of community cards dealt
	})

	return nil
}

func (g *Game) DealRiver() error {
	if len(g.deck) < 2 {
		return fmt.Errorf("not enough cards in deck for river")
	}

	// Burn one card, then deal 1 community card
	g.deck = g.deck[1:] // Burn card

	g.communityCards = append(g.communityCards, g.deck[0])
	g.deck = g.deck[1:]

	g.currentPhase = PhaseRiver

	// Log system action for dealing river
	g.TakeSystemAction(Action{
		PlayerID: SystemPlayerID,
		Type:     ActionSystemDealRiver,
		Amount:   1, // Number of community cards dealt
	})

	return nil
}

// NewGame creates a new game with specified blinds
func NewGame(smallBlind, bigBlind int) *Game {
	game := &Game{
		players:        [10]IPlayer{},
		deck:           newStandardDeck(), // Use standard 52-card deck
		communityCards: poker.Cards{},
		currentPhase:   PhasePreflop,
		smallBlind:     smallBlind,
		bigBlind:       bigBlind,
		systemActions: SystemActions{
			Preflop: []Action{},
			Flop:    []Action{},
			Turn:    []Action{},
			River:   []Action{},
		},
		userActions: UserActions{
			Preflop: []Action{},
			Flop:    []Action{},
			Turn:    []Action{},
			River:   []Action{},
		},
	}

	// Shuffle deck on creation (without logging since it's initialization)
	rand.Seed(time.Now().UnixNano())
	for i := len(game.deck) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		game.deck[i], game.deck[j] = game.deck[j], game.deck[i]
	}

	return game
}
