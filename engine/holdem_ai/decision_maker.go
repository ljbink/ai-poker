package holdem_ai

// DecisionMaker interface that both human players and AI bots must implement
// Kept truly minimal with only the essential decision-making method
type DecisionMaker interface {
	// MakeDecision returns a channel that will receive the chosen action
	// This allows for asynchronous decision making and timeout handling
	MakeDecision() <-chan Action
}
