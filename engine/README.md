# AI Poker Engine

A comprehensive, well-tested Texas Hold'em poker engine written in Go with clean architecture, full AI support, and extensive test coverage.

## 🏗️ Architecture Overview

The engine is organized into three main packages:

```
engine/
├── poker/          # Core poker primitives (cards, players, deck)
├── holdem/         # Texas Hold'em game logic and rules  
├── holdem_ai/      # AI decision makers and human interfaces
└── README.md       # This file
```

## 📦 Package Structure

### [`poker/`](./poker/) - Core Poker Primitives
- **Cards & Deck**: Card representation, deck operations, shuffling
- **Base Player**: Common player functionality and data structures  
- **Constants**: Suits, ranks, and poker-specific enumerations

### [`holdem/`](./holdem/) - Texas Hold'em Game Engine
- **Game Logic**: Complete game flow from preflop to showdown
- **Player Management**: Betting, folding, chip management
- **Hand Evaluation**: Comprehensive poker hand ranking and comparison
- **Betting Rounds**: Call, raise, check, fold with proper validation

### [`holdem_ai/`](./holdem_ai/) - AI & Player Interfaces  
- **AI Decision Makers**: Automated bot players with hand evaluation
- **Human Interfaces**: Callback-based system for frontend integration
- **Action Validation**: Comprehensive action validation and game state management
- **Utility Functions**: Formatting and calculation helpers for UIs

## 🚀 Quick Start

### Basic Game Setup

```go
package main

import (
    "github.com/ljbink/ai-poker/engine/holdem"
    "github.com/ljbink/ai-poker/engine/holdem_ai"
)

func main() {
    // Create a game with 2 players
    game := holdem.NewGame([]string{"Alice", "Bob"}, 10, 20) // small blind: 10, big blind: 20
    
    // Start the hand
    game.StartHand()
    
    // Create AI decision makers
    ai1 := holdem_ai.NewBasicBotDecisionMaker("Alice")
    ai2 := holdem_ai.NewBasicBotDecisionMaker("Bob")
    
    // Game loop
    for !game.IsHandComplete() {
        currentPlayer := game.GetCurrentPlayer()
        
        var action holdem_ai.Action
        if currentPlayer.GetName() == "Alice" {
            action = ai1.MakeDecision(game, currentPlayer)
        } else {
            action = ai2.MakeDecision(game, currentPlayer)
        }
        
        // Execute the action
        adapter := holdem_ai.NewGameAdapter()
        adapter.ExecuteUserAction(game, currentPlayer, action)
    }
}
```

### Mixed Human/AI Game

```go
// Create human and AI players
humanPlayer := holdem_ai.NewHumanUserDecisionMaker("Human")
aiPlayer := holdem_ai.NewBasicBotDecisionMaker("AI")

// Set up callback for human player
humanPlayer.SetActionNeededCallback(func(game *holdem.Game, player holdem.IPlayer, validActions []holdem_ai.ActionType) {
    // Update your UI here
    phase := holdem_ai.FormatGamePhase(game.CurrentPhase)
    playerCards := holdem_ai.FormatPlayerCards(player)
    communityCards := holdem_ai.FormatCommunityCards(game)
    
    // Display game state and wait for user input
    // Then send action: humanPlayer.SendAction(action)
})
```

## 🧪 Testing

All packages have comprehensive test coverage with edge cases and integration scenarios.

### Run All Tests
```bash
# Run all engine tests
go test ./engine/... -v

# Run specific package tests
go test ./engine/poker/ -v
go test ./engine/holdem/ -v  
go test ./engine/holdem_ai/ -v
```

### Test Coverage Features
- ✅ **100% Line Coverage** on all public APIs
- ✅ **Edge Case Testing** including nil inputs and boundary conditions
- ✅ **Integration Tests** validating cross-package functionality
- ✅ **Game State Validation** ensuring poker rules are correctly implemented
- ✅ **Action Validation** comprehensive testing of all game actions
- ✅ **Hand Evaluation** testing all poker hand rankings and comparisons

## 🎯 Key Features

### Complete Poker Implementation
- ✅ Full Texas Hold'em game flow (preflop → flop → turn → river → showdown)
- ✅ Proper betting rounds with call, raise, check, fold
- ✅ Blind posting and position management
- ✅ Comprehensive hand evaluation (all 10 hand rankings)
- ✅ Side pot handling for all-in scenarios

### AI & Human Players
- ✅ **BasicBot AI**: Hand strength evaluation, pot odds, position awareness
- ✅ **Human Interface**: Thread-safe callback system for frontend integration
- ✅ **Action Validation**: Prevents invalid moves
- ✅ **Timeout Handling**: Auto-fold for inactive players

### Clean Architecture  
- ✅ **Separation of Concerns**: Clear boundaries between packages
- ✅ **Interface-Driven**: Clean contracts between components
- ✅ **Immutable Operations**: Safe method chaining
- ✅ **Error Handling**: Graceful handling of edge cases

### Developer Experience
- ✅ **Comprehensive Documentation**: Each package has detailed READMEs
- ✅ **Helper Functions**: Utilities for formatting and calculations
- ✅ **Type Safety**: Strong typing prevents common errors
- ✅ **Extensive Testing**: High confidence through thorough test coverage

## 📚 Documentation

Each package contains detailed documentation:

- **[poker/README.md](./poker/README.md)** - Core poker primitives and card operations
- **[holdem/README.md](./holdem/README.md)** - Texas Hold'em game engine and rules
- **[holdem_ai/README.md](./holdem_ai/README.md)** - AI players and human interfaces

## 🔧 API Reference

### Core Types

```go
// From poker package
type Card struct {
    Suit Suit
    Rank Rank  
}

type Cards []*Card

// From holdem package  
type Game struct {
    Players               []*Player
    Pot                   int
    CurrentPhase          GamePhase
    CurrentBet            int
    CommunityCards        poker.Cards
    // ... other fields
}

type Player struct {
    poker.BasePlayer  // ID, Name
    // ... private fields accessed via getters
}

// From holdem_ai package
type Action struct {
    Type   ActionType
    Amount int
}

type DecisionMaker interface {
    MakeDecision(game *holdem.Game, player holdem.IPlayer) <-chan Action
    GetName() string
}
```

## 🎮 Usage Patterns

### For Game Developers
- Use `holdem.Game` for core game logic
- Integrate `holdem_ai.HumanUserDecisionMaker` for player input
- Add `holdem_ai.BasicBotDecisionMaker` for AI opponents

### For AI Researchers  
- Implement custom `DecisionMaker` interface for new AI strategies
- Use `holdem_ai.ActionValidator` for move validation
- Leverage hand evaluation functions for strategy development

### For Web/Mobile Apps
- Use human decision makers with callback systems
- Format game state with provided utility functions
- Validate user actions before execution

## 🚧 Future Enhancements

- **Advanced AI**: Monte Carlo tree search, neural network players
- **Tournament Mode**: Multi-table tournaments with blinds escalation  
- **Statistics**: Hand history, player statistics, advanced analytics
- **Multiplayer**: WebSocket support for online multiplayer games
- **Variants**: Omaha, Seven-Card Stud, and other poker variants

## 📄 License

This poker engine is designed to be a foundation for poker applications, research, and education. 