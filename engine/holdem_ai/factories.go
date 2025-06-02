package holdem_ai

import "math/rand"

// Factory functions for creating different types of decision makers

// CreateBasicBot creates a basic bot with default settings
func CreateBasicBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.5, 0.1) // Moderate aggressiveness, low bluff frequency
}

// CreateConservativeBot creates a conservative bot
func CreateConservativeBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.2, 0.05) // Low aggressiveness, very low bluff frequency
}

// CreateAggressiveBot creates an aggressive bot
func CreateAggressiveBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.8, 0.25) // High aggressiveness, moderate bluff frequency
}

// CreateTightBot creates a very conservative bot
func CreateTightBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.1, 0.01) // Very low aggressiveness, almost no bluffs
}

// CreateLooseBot creates a loose aggressive bot
func CreateLooseBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.9, 0.4) // Very high aggressiveness, frequent bluffs
}

// CreateRandomBot creates a bot with random settings
func CreateRandomBot() IDecisionMaker {
	// Random aggressiveness between 0.3 and 0.9
	// Random bluff frequency between 0.05 and 0.3
	aggressiveness := 0.3 + (0.6 * rand.Float64())
	bluffFreq := 0.05 + (0.25 * rand.Float64())
	return NewBasicBotDecisionMaker(aggressiveness, bluffFreq)
}

// CreateCustomBot creates a bot with custom settings
func CreateCustomBot(aggressiveness, bluffFrequency float64) IDecisionMaker {
	return NewBasicBotDecisionMaker(aggressiveness, bluffFrequency)
}

// Preset bot personalities for common styles

// CreateNitBot creates an extremely tight/conservative bot
func CreateNitBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.05, 0.0) // Extremely low aggressiveness, never bluffs
}

// CreateManiacBot creates an extremely loose/aggressive bot
func CreateManiacBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.95, 0.5) // Maximum aggressiveness, frequent bluffs
}

// CreateBalancedBot creates a well-balanced bot
func CreateBalancedBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.6, 0.15) // Balanced aggressiveness and bluff frequency
}

// CreateCallingStationBot creates a bot that calls frequently but rarely raises
func CreateCallingStationBot() IDecisionMaker {
	return NewBasicBotDecisionMaker(0.3, 0.02) // Low aggressiveness, almost no bluffs
}
