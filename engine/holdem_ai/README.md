# Texas Hold'em AI Engine

This package provides a comprehensive AI system for Texas Hold'em poker, including both automated AI decision makers and human player interfaces.

## Architecture Overview

```
holdem_ai/
├── action.go                       # Action types and validation
├── decision_maker.go              # DecisionMaker interface
├── decision_maker_basic_bot.go    # Basic AI implementation
├── decision_maker_human.go        # Human player interface
├── game_adapter.go               # Game state management
└── examples/                      # Usage examples
```

## Core Components

### DecisionMaker Interface

All decision makers (AI and human) implement this interface:

```go
type DecisionMaker interface {
    MakeDecision(game *holdem.Game, player holdem.IPlayer) Action
    GetName() string
    IsBot() bool
}
```

### Action Types

The system supports all standard poker actions:

- `ActionFold` - Fold current hand
- `ActionCheck` - Check (bet 0 when no bet is required)
- `ActionCall` - Call current bet
- `ActionRaise` - Raise by specified amount
- `ActionAllIn` - Bet all remaining chips

## AI Players

### BasicBotDecisionMaker

A simple AI that makes decisions based on:
- Hand strength evaluation
- Pot odds calculation
- Position awareness
- Basic bluffing logic

```go
// Create an AI player
bot := holdem_ai.NewBasicBotDecisionMaker("AI Player")

// Use in game
action := bot.MakeDecision(game, player)
```

## Human Players

### HumanUserDecisionMaker

Provides interface for human players through callbacks and channels:

```go
// Create human decision maker
human := holdem_ai.NewHumanUserDecisionMaker("Human Player")

// Set up frontend callback
human.SetActionNeededCallback(func(game *holdem.Game, player holdem.IPlayer, validActions []holdem_ai.ActionType) {
    // Update UI with game state using helper functions
    // Show valid actions
    // Wait for user input
})

// Frontend sends user action
action := holdem_ai.CreateRaiseAction(50)
human.SendAction(action)
```

#### Key Features

- **Thread-safe**: Actions can be sent from any goroutine
- **Timeout handling**: Auto-fold if no action within timeout
- **Action validation**: Ensures only valid actions are accepted
- **Direct access**: Use existing `*holdem.Game` and `holdem.IPlayer` types
- **Callback system**: Notifies frontend when action needed
- **Helper functions**: Utilities for formatting display information

## Usage Examples

### Setting Up a Mixed Game

```go
// Create game
game := holdem.NewGame([]string{"Human", "AI1", "AI2"}, 5, 10)

// Create decision makers
human := holdem_ai.NewHumanUserDecisionMaker("Human")
ai1 := holdem_ai.NewBasicBotDecisionMaker("AI1")
ai2 := holdem_ai.NewBasicBotDecisionMaker("AI2")

// Set up human player callback
human.SetActionNeededCallback(func(game *holdem.Game, player holdem.IPlayer, validActions []holdem_ai.ActionType) {
    // Display game state in UI
    // Show available actions
    // Handle user input
})

// Game loop
game.StartHand()
for !game.IsHandComplete() {
    currentPlayer := game.GetCurrentPlayer()
    
    var action holdem_ai.Action
    if currentPlayer.GetName() == "Human" {
        action = human.MakeDecision(game, currentPlayer)
    } else if currentPlayer.GetName() == "AI1" {
        action = ai1.MakeDecision(game, currentPlayer)
    } else {
        action = ai2.MakeDecision(game, currentPlayer)
    }
    
    // Execute action
    adapter := holdem_ai.NewGameAdapter()
    adapter.ExecuteAction(game, action)
}
```

### Frontend Integration

For TUI/GUI applications:

```go
// Create human decision maker
humanDM := holdem_ai.NewHumanUserDecisionMaker("Player")

// Set up UI callback
humanDM.SetActionNeededCallback(func(game *holdem.Game, player holdem.IPlayer, validActions []holdem_ai.ActionType) {
    // Format game information for display
    phase := holdem_ai.FormatGamePhase(game.CurrentPhase)
    playerCards := holdem_ai.FormatPlayerCards(player)
    communityCards := holdem_ai.FormatCommunityCards(game)
    callAmount := holdem_ai.CalculateCallAmount(game, player)
    potOdds := game.CalculatePotOdds(player)
    
    // Update game display
    updateGameDisplay(phase, game.Pot, playerCards, communityCards, callAmount, potOdds)
    
    // Show action buttons
    showActionButtons(validActions)
})

// Handle button clicks
func onFoldClick() {
    humanDM.SendAction(holdem_ai.CreateFoldAction())
}

func onRaiseClick(amount int) {
    humanDM.SendAction(holdem_ai.CreateRaiseAction(amount))
}
```

## Helper Functions

The package provides utility functions for formatting game information:

```go
// Format game phase as string
phase := holdem_ai.FormatGamePhase(game.CurrentPhase) // "Preflop", "Flop", etc.

// Format player cards as string
playerCards := holdem_ai.FormatPlayerCards(player) // "A♠ K♥"

// Format community cards as string
communityCards := holdem_ai.FormatCommunityCards(game) // "A♠ K♥ Q♦"

// Calculate betting amounts
callAmount := holdem_ai.CalculateCallAmount(game, player)
minRaise := holdem_ai.CalculateMinRaise(game)
maxRaise := holdem_ai.CalculateMaxRaise(game, player)

// Get pot odds (from game)
potOdds := game.CalculatePotOdds(player)
```

## Action Validation

The system automatically validates all actions:

```go
validator := holdem_ai.NewActionValidator()

// Check if action is valid
isValid := validator.IsValidAction(action, game, player)

// Get all valid actions
validActions := validator.GetValidActions(game, player)
```

## Timeout Handling

Human decision makers support configurable timeouts:

```go
// Set 30 second timeout
human.SetTimeout(30 * time.Second)

// If no action received within timeout, automatically folds
```

## Testing

Run all tests:

```bash
cd engine/holdem_ai
go test -v
```

## Examples

See the `examples/` directory for complete working examples:

- `human_integration_example.go` - Shows frontend integration
- More examples coming soon...

## Advanced Usage

### Custom AI Implementation

Implement the `DecisionMaker` interface:

```go
type CustomBot struct {
    name string
}

func (c *CustomBot) MakeDecision(game *holdem.Game, player holdem.IPlayer) holdem_ai.Action {
    // Your custom AI logic here
    return holdem_ai.CreateFoldAction()
}

func (c *CustomBot) GetName() string { return c.name }
func (c *CustomBot) IsBot() bool { return true }
```

### Direct Game State Access

Since the system uses existing holdem types directly, you have full access to all game information:

```go
// In your callback function
func onActionNeeded(game *holdem.Game, player holdem.IPlayer, validActions []holdem_ai.ActionType) {
    // Access any game information directly
    pot := game.Pot
    currentBet := game.CurrentBet
    phase := game.CurrentPhase
    communityCards := game.CommunityCards
    playerChips := player.GetChips()
    playerBet := player.GetBet()
    playerCards := player.GetHandCards()
    
    // Access all players
    for i, p := range game.Players {
        fmt.Printf("Player %d: %s has %d chips\n", i, p.GetName(), p.GetChips())
    }
}
```

## Contributing

When adding new features:

1. Implement proper interfaces
2. Add comprehensive tests
3. Update documentation
4. Follow Go conventions
5. Ensure thread safety for concurrent use
6. Reuse existing holdem types instead of creating new ones 