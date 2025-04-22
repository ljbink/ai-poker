package poker

import (
	"math/rand"
	"time"
)

type CardBooleanClosure = func(*Card) bool
type CardCountsClosure = func(val Rank, count int) bool

type Cards []*Card

func (c Cards) Length() int {
	return len(c)
}

func (c *Cards) Append(cards ...*Card) {
	*c = append(*c, cards...)
}

func (c *Cards) Remove(cards ...*Card) int {
	count := 0
	temp := Cards{}
	for _, _c := range *c {
		for _, card := range cards {
			if _c.String() == card.String() {
				count++
				goto NextLoop
			}
		}
		temp = append(temp, _c)
	NextLoop:
	}
	*c = temp
	return count
}

func (c *Cards) Shuffle() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < c.Length(); i++ {
		j := rand.Intn(c.Length())
		if i == j {
			continue
		}
		(*c)[i], (*c)[j] = (*c)[j], (*c)[i]
	}
}

func (c Cards) String() string {
	res := ""
	for i, _c := range c {
		res += _c.String()
		if i < c.Length()-1 {
			res += " "
		}
	}
	return res
}

func NewDeckCards() Cards {
	suits := []Suit{
		SuitHeart,
		SuitDiamond,
		SuitClub,
		SuitSpade,
	}
	ranks := []Rank{
		RankAce,
		RankTwo,
		RankThree,
		RankFour,
		RankFive,
		RankSix,
		RankSeven,
		RankEight,
		RankNine,
		RankTen,
		RankJack,
		RankQueen,
		RankKing,
	}
	cards := Cards{}
	cards.Append(NewCard(SuitNone, RankColoredJoker))
	cards.Append(NewCard(SuitNone, RankJoker))
	for _, suit := range suits {
		for _, rank := range ranks {
			cards.Append(NewCard(suit, rank))
		}
	}
	return cards
}
