package holdem

import "github.com/ljbink/ai-poker/engine/poker"

type IPlayer interface {
	GetID() int
	GetName() string

	DealCard(card *poker.Card) IPlayer
	GetHandCards() []*poker.Card

	GrandChips(amount int) IPlayer
	GetChips() int

	GetBet() int
	GetTotalBet() int
	Bet(amount int) IPlayer
	ResetBet() IPlayer

	IsFolded() bool
	Fold() IPlayer

	ResetForNewHand() IPlayer
}

type Player struct {
	poker.BasePlayer
	cards    []*poker.Card
	chips    int
	bet      int
	totalBet int
	folded   bool
}

func NewPlayer(id int, name string, startingChips int) IPlayer {
	return &Player{
		BasePlayer: poker.BasePlayer{ID: id, Name: name},
		chips:      startingChips,
	}
}

func (p *Player) DealCard(card *poker.Card) IPlayer {
	p.cards = append(p.cards, card)
	return p
}

func (p *Player) GetHandCards() []*poker.Card {
	return p.cards
}

func (p *Player) GetChips() int {
	return p.chips
}

func (p *Player) GetBet() int {
	return p.bet
}

func (p *Player) GetTotalBet() int {
	return p.totalBet
}

func (p *Player) GetID() int {
	return p.ID
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) IsFolded() bool {
	return p.folded
}

func (p *Player) Fold() IPlayer {
	p.folded = true
	return p
}

func (p *Player) GrandChips(amount int) IPlayer {
	p.chips += amount
	return p
}

func (p *Player) Bet(amount int) IPlayer {
	p.bet += amount
	p.totalBet += amount
	p.chips -= amount
	return p
}

func (p *Player) ResetBet() IPlayer {
	p.bet = 0
	return p
}

func (p *Player) ResetForNewHand() IPlayer {
	p.cards = nil
	p.bet = 0
	p.totalBet = 0
	p.folded = false
	return p
}
