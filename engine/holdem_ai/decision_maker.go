package holdem_ai

import "github.com/ljbink/ai-poker/engine/holdem"

// IDecisionMaker interface that both human players and AI bots must implement
// Kept truly minimal with only the essential decision-making method
type IDecisionMaker interface {
	// MakeDecision returns a channel that will receive the chosen action
	// This allows for asynchronous decision making and timeout handling
	// Takes game and player as parameters to make IDecisionMakers stateless
	MakeDecision(game *holdem.Game, player holdem.IPlayer) <-chan holdem.Action
}
