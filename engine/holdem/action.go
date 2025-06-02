package holdem

type ActionType int

const (
	// Player Actions
	ActionFold ActionType = iota
	ActionCheck
	ActionCall
	ActionRaise
	ActionAllIn

	// System Actions
	ActionSystemShuffle     // Deck shuffle
	ActionSystemDealHole    // Deal hole cards
	ActionSystemDealFlop    // Deal flop cards
	ActionSystemDealTurn    // Deal turn card
	ActionSystemDealRiver   // Deal river card
	ActionSystemPhaseChange // Phase transition
)

const SystemPlayerID = -1

type Action struct {
	PlayerID int
	Type     ActionType
	Amount   int
}
