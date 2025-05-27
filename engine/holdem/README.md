# Texas Hold'em with Clean Player Architecture

This is a clean Texas Hold'em implementation using a well-structured Player design with private fields and getter methods.

## Structure

### Files
- `game.go` - Main game logic and structure
- `player.go` - Player structure with encapsulated fields
- `action.go` - Basic action types
- `game_test.go` - Comprehensive tests

### Core Types

```go
type Player struct {
    poker.BasePlayer  // ID, Name
    cards    []*poker.Card
    chips    int
    bet      int
    totalBet int
    folded   bool
}

type Game struct {
    Players               []*Player
    Deck                  poker.Cards
    Pot                   int
    CurrentPhase          GamePhase
    CurrentPlayerPosition int
    DealerPosition        int
    SmallBlind            int
    BigBlind              int
    CurrentBet            int
    CommunityCards        poker.Cards
}
```

## Player Methods

### Creation
```go
player := NewPlayer(id, name, startingChips)
```

### Card Management
```go
player.DealCard(card)           // Add a card to hand
cards := player.GetHandCards()  // Get all cards
```

### Chip Management
```go
chips := player.GetChips()      // Get current chips
player.GrandChips(amount)       // Add chips
player.Bet(amount)              // Bet chips
bet := player.GetBet()          // Get current bet
total := player.GetTotalBet()   // Get total bet this hand
```

### Game State
```go
folded := player.IsFolded()     // Check if folded
player.Fold()                   // Fold the hand
player.ResetBet()               // Reset bet for new round
player.ResetForNewHand()        // Reset for new hand
```

## Game Usage

```go
// Create game
game := holdem.NewGame([]string{"Alice", "Bob"}, 5, 10)

// Start hand
game.StartHand()

// Player actions
game.Call()         // Current player calls
game.Raise(20)      // Current player raises by 20
game.Check()        // Current player checks
game.Fold()         // Current player folds

// Advance phases
game.NextPhase()    // Preflop -> Flop -> Turn -> River -> Showdown

// Get info
currentPlayer := game.GetCurrentPlayer()
activePlayers := game.GetActivePlayers()
isComplete := game.IsBettingRoundComplete()
```

## Game Flow

1. **Create Game**: `NewGame()` with player names and blinds
2. **Start Hand**: `StartHand()` deals cards and posts blinds
3. **Player Actions**: `Call()`, `Raise()`, `Check()`, `Fold()`
4. **Check Completion**: `IsBettingRoundComplete()` to check if round is done
5. **Phase Advancement**: `NextPhase()` moves through betting rounds
6. **Next Hand**: `NextHand()` moves dealer button

## Features

✅ **Clean Player Architecture**: Encapsulated fields with getter methods  
✅ **Proper Inheritance**: Uses `poker.BasePlayer` for common fields  
✅ **Complete Game Flow**: All phases from preflop to showdown  
✅ **Betting Logic**: Call, raise, check, fold with proper chip management  
✅ **Blind Posting**: Automatic small/big blind handling  
✅ **Card Dealing**: Proper deck shuffling and community cards  
✅ **Player Tracking**: Folded players, active players, betting completion  
✅ **Position Management**: Dealer button, player positions  
✅ **Comprehensive Tests**: All functionality tested  

## Architecture Benefits

- **Encapsulation**: Private fields prevent direct manipulation
- **Clean Interface**: Clear getter/setter methods
- **Inheritance**: Reuses `BasePlayer` for common functionality
- **Immutable Operations**: Methods return `*Player` for chaining
- **Type Safety**: Strong typing prevents errors

## Testing

```bash
go test ./engine/holdem -v
```

All core functionality is covered:
- Game creation and initialization
- Hand dealing and blind posting
- Player actions (call, raise, check, fold)
- Phase progression
- Betting round completion

## Demo

```bash
go run main.go
```

Shows a complete 2-player game with:
- Hand initialization
- Hole card dealing
- Blind posting
- Player actions
- Phase advancement to flop
- Community card dealing

## Next Steps

This clean foundation makes it easy to add:

1. **Hand Evaluation**: Determine winners at showdown
2. **All-in Logic**: Handle players with insufficient chips
3. **Side Pots**: Multiple pots for all-in scenarios
4. **Tournament Mode**: Blind escalation and elimination
5. **AI Players**: Automated decision making
6. **Web Interface**: REST API or WebSocket interface 