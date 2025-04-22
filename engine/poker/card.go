package poker

import (
	"fmt"
)

type Suit uint8

const (
	SuitNone Suit = iota
	SuitHeart
	SuitDiamond
	SuitClub
	SuitSpade
)

var (
	suitMap = map[Suit]string{
		SuitNone:    "",
		SuitHeart:   "Heart",
		SuitDiamond: "Diamond",
		SuitClub:    "Club",
		SuitSpade:   "Spade",
	}
)

type Rank uint8

const (
	RankNone Rank = iota
	RankAce
	RankTwo
	RankThree
	RankFour
	RankFive
	RankSix
	RankSeven
	RankEight
	RankNine
	RankTen
	RankJack
	RankQueen
	RankKing
	RankJoker
	RankColoredJoker
)

var (
	RankMap = map[Rank]string{
		RankNone:         "",
		RankAce:          "A",
		RankTwo:          "2",
		RankThree:        "3",
		RankFour:         "4",
		RankFive:         "5",
		RankSix:          "6",
		RankSeven:        "7",
		RankEight:        "8",
		RankNine:         "9",
		RankTen:          "10",
		RankJack:         "J",
		RankQueen:        "Q",
		RankKing:         "K",
		RankJoker:        "Joker",
		RankColoredJoker: "ColoredJoker",
	}
)

var (
	cardUnicodeMap = map[string]string{
		fmt.Sprintf("%d%d", SuitNone, RankNone): "🂠",

		fmt.Sprintf("%d%d", SuitNone, RankColoredJoker): "🃏",
		fmt.Sprintf("%d%d", SuitNone, RankJoker):        "🃟",

		fmt.Sprintf("%d%d", SuitHeart, RankAce):   "🂱",
		fmt.Sprintf("%d%d", SuitHeart, RankTwo):   "🂲",
		fmt.Sprintf("%d%d", SuitHeart, RankThree): "🂳",
		fmt.Sprintf("%d%d", SuitHeart, RankFour):  "🂴",
		fmt.Sprintf("%d%d", SuitHeart, RankFive):  "🂵",
		fmt.Sprintf("%d%d", SuitHeart, RankSix):   "🂶",
		fmt.Sprintf("%d%d", SuitHeart, RankSeven): "🂷",
		fmt.Sprintf("%d%d", SuitHeart, RankEight): "🂸",
		fmt.Sprintf("%d%d", SuitHeart, RankNine):  "🂹",
		fmt.Sprintf("%d%d", SuitHeart, RankTen):   "🂺",
		fmt.Sprintf("%d%d", SuitHeart, RankJack):  "🂫",
		fmt.Sprintf("%d%d", SuitHeart, RankQueen): "🂭",
		fmt.Sprintf("%d%d", SuitHeart, RankKing):  "🂮",

		fmt.Sprintf("%d%d", SuitDiamond, RankAce):   "🃁",
		fmt.Sprintf("%d%d", SuitDiamond, RankTwo):   "🃂",
		fmt.Sprintf("%d%d", SuitDiamond, RankThree): "🃃",
		fmt.Sprintf("%d%d", SuitDiamond, RankFour):  "🃄",
		fmt.Sprintf("%d%d", SuitDiamond, RankFive):  "🃅",
		fmt.Sprintf("%d%d", SuitDiamond, RankSix):   "🃆",
		fmt.Sprintf("%d%d", SuitDiamond, RankSeven): "🃇",
		fmt.Sprintf("%d%d", SuitDiamond, RankEight): "🃈",
		fmt.Sprintf("%d%d", SuitDiamond, RankNine):  "🃉",
		fmt.Sprintf("%d%d", SuitDiamond, RankTen):   "🃊",
		fmt.Sprintf("%d%d", SuitDiamond, RankJack):  "🃋",
		fmt.Sprintf("%d%d", SuitDiamond, RankQueen): "🃍",
		fmt.Sprintf("%d%d", SuitDiamond, RankKing):  "🃎",

		fmt.Sprintf("%d%d", SuitClub, RankAce):   "🃑",
		fmt.Sprintf("%d%d", SuitClub, RankTwo):   "🃒",
		fmt.Sprintf("%d%d", SuitClub, RankThree): "🃓",
		fmt.Sprintf("%d%d", SuitClub, RankFour):  "🃔",
		fmt.Sprintf("%d%d", SuitClub, RankFive):  "🃕",
		fmt.Sprintf("%d%d", SuitClub, RankSix):   "🃖",
		fmt.Sprintf("%d%d", SuitClub, RankSeven): "🃗",
		fmt.Sprintf("%d%d", SuitClub, RankEight): "🃘",
		fmt.Sprintf("%d%d", SuitClub, RankNine):  "🃙",
		fmt.Sprintf("%d%d", SuitClub, RankTen):   "🃚",
		fmt.Sprintf("%d%d", SuitClub, RankJack):  "🃛",
		fmt.Sprintf("%d%d", SuitClub, RankQueen): "🃝",
		fmt.Sprintf("%d%d", SuitClub, RankKing):  "🃞",

		fmt.Sprintf("%d%d", SuitSpade, RankAce):   "🂡",
		fmt.Sprintf("%d%d", SuitSpade, RankTwo):   "🂢",
		fmt.Sprintf("%d%d", SuitSpade, RankThree): "🂣",
		fmt.Sprintf("%d%d", SuitSpade, RankFour):  "🂤",
		fmt.Sprintf("%d%d", SuitSpade, RankFive):  "🂥",
		fmt.Sprintf("%d%d", SuitSpade, RankSix):   "🂦",
		fmt.Sprintf("%d%d", SuitSpade, RankSeven): "🂧",
		fmt.Sprintf("%d%d", SuitSpade, RankEight): "🂨",
		fmt.Sprintf("%d%d", SuitSpade, RankNine):  "🂩",
		fmt.Sprintf("%d%d", SuitSpade, RankTen):   "🂪",
		fmt.Sprintf("%d%d", SuitSpade, RankJack):  "🂻",
		fmt.Sprintf("%d%d", SuitSpade, RankQueen): "🂽",
		fmt.Sprintf("%d%d", SuitSpade, RankKing):  "🂾",
	}
)

type Card struct {
	Suit Suit
	Rank Rank
}

func (r Card) String() string {
	display, ok := cardUnicodeMap[fmt.Sprintf("%d%d", r.Suit, r.Rank)]
	if ok {
		return display
	}
	return fmt.Sprintf("%s-%s", suitMap[r.Suit], RankMap[r.Rank])
}

func NewCard(suit Suit, rank Rank) *Card {
	return &Card{
		Suit: suit,
		Rank: rank,
	}
}
