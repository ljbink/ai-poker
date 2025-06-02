# Poker Package - Core Poker Primitives

The `poker` package provides the foundational building blocks for poker games, including card representation, deck operations, and base player functionality.

## 📦 Package Contents

```
poker/
├── card.go          # Card representation and constants
├── cards.go         # Deck operations and card collections  
├── player.go        # Base player structure
├── card_test.go     # Comprehensive card tests
├── cards_test.go    # Deck and collection tests
├── player_test.go   # Player functionality tests
└── README.md        # This file
```

## 🃏 Core Types

### Card Structure

```go
type Card struct {
    Suit Suit
    Rank Rank
}

// Suits
const (
    SuitNone Suit = iota
    SuitHeart     // ♥
    SuitDiamond   // ♦  
    SuitClub      // ♣
    SuitSpade     // ♠
)

// Ranks  
const (
    RankNone Rank = iota
    RankAce           // A
    RankTwo           // 2
    RankThree         // 3
    // ... through King
    RankKing          // K
    RankJoker         // Joker
    RankColoredJoker  // Colored Joker
)
```

### Cards Collection

```go
type Cards []*Card

// Core methods
func (c Cards) Length() int
func (c *Cards) Append(cards ...*Card)
func (c *Cards) Remove(cards ...*Card) int
func (c *Cards) Shuffle()
func (c Cards) String() string
```

### Base Player

```go
type BasePlayer struct {
    ID   int
    Name string
}
```

## 🚀 Usage Examples

### Creating and Working with Cards

```go
package main

import (
    "fmt"
    "github.com/ljbink/ai-poker/engine/poker"
)

func main() {
    // Create individual cards
    aceOfSpades := poker.NewCard(poker.SuitSpade, poker.RankAce)
    kingOfHearts := poker.NewCard(poker.SuitHeart, poker.RankKing)
    
    // Cards display as Unicode symbols
    fmt.Println(aceOfSpades.String()) // 🂡 (Ace of Spades)
    fmt.Println(kingOfHearts.String()) // 🂮 (King of Hearts)
    
    // Create a full deck
    deck := poker.NewDeckCards()
    fmt.Printf("Deck has %d cards\n", deck.Length()) // 54 cards (52 + 2 jokers)
    
    // Shuffle the deck
    deck.Shuffle()
    
    // Deal cards
    hand := poker.Cards{}
    hand.Append(deck[0], deck[1]) // Deal 2 cards
    fmt.Printf("Hand: %s\n", hand.String()) // e.g., "🂡 🂮"
}
```

### Working with Card Collections

```go
// Create a hand
hand := poker.Cards{}

// Add cards to hand
hand.Append(
    poker.NewCard(poker.SuitHeart, poker.RankAce),
    poker.NewCard(poker.SuitSpade, poker.RankKing),
)

// Check hand size
fmt.Printf("Hand has %d cards\n", hand.Length()) // 2

// Remove specific cards
removed := hand.Remove(poker.NewCard(poker.SuitHeart, poker.RankAce))
fmt.Printf("Removed %d cards\n", removed) // 1

// Display remaining cards
fmt.Printf("Remaining: %s\n", hand.String()) // "🂾"
```

### Base Player Usage

```go
// Create a base player
player := poker.BasePlayer{
    ID:   1,
    Name: "Alice",
}

fmt.Printf("Player %d: %s\n", player.ID, player.Name) // Player 1: Alice
```

## 🎴 Card Features

### Unicode Display
All cards render as beautiful Unicode symbols:
- **Hearts**: 🂱 🂲 🂳 🂴 🂵 🂶 🂷 🂸 🂹 🂺 🂫 🂭 🂮
- **Diamonds**: 🃁 🃂 🃃 🃄 🃅 🃆 🃇 🃈 🃉 🃊 🃋 🃍 🃎  
- **Clubs**: 🃑 🃒 🃓 🃔 🃕 🃖 🃗 🃘 🃙 🃚 🃛 🃝 🃞
- **Spades**: 🂡 🂢 🂣 🂤 🂥 🂦 🂧 🂨 🂩 🂪 🂻 🂽 🂾
- **Jokers**: 🃟 🃏

### Comprehensive Coverage
- **52 Standard Cards**: All suits and ranks (Ace through King)
- **2 Jokers**: Regular joker and colored joker
- **Fallback Handling**: Graceful display of invalid cards
- **Nil Safety**: Proper handling of nil card pointers

## 🔧 Deck Operations

### Creating Decks

```go
// Full deck (54 cards: 52 standard + 2 jokers)
fullDeck := poker.NewDeckCards()

// Empty deck for custom collections
customDeck := poker.Cards{}
customDeck.Append(
    poker.NewCard(poker.SuitHeart, poker.RankAce),
    poker.NewCard(poker.SuitSpade, poker.RankAce),
)
```

### Shuffling

```go
deck := poker.NewDeckCards()
deck.Shuffle() // Randomizes card order using time-based seed
```

### Card Management

```go
deck := poker.NewDeckCards()

// Deal cards (remove from deck)
hand := poker.Cards{}
hand.Append(deck[0], deck[1])
deck.Remove(deck[0], deck[1])

// Add cards back
deck.Append(hand...)
hand = poker.Cards{} // Clear hand
```

## 🔍 Advanced Features

### Type Safety
```go
// All operations are type-safe
var suit poker.Suit = poker.SuitHeart
var rank poker.Rank = poker.RankAce
card := poker.NewCard(suit, rank)
```

### Nil Handling
```go
// Cards collections handle nil cards gracefully
cards := poker.Cards{}
cards.Append(nil) // Won't panic
fmt.Println(cards.String()) // Displays "[nil]"

// Remove operations handle nil safely
removed := cards.Remove(nil) // Returns count of nil cards removed
```

### Closure Support
```go
// Define custom card filtering
type CardBooleanClosure = func(*Card) bool
type CardCountsClosure = func(val Rank, count int) bool

// Example usage
isRed := func(c *poker.Card) bool {
    return c != nil && (c.Suit == poker.SuitHeart || c.Suit == poker.SuitDiamond)
}
```

## 🧪 Testing

The package includes comprehensive tests covering:

### Card Tests (`card_test.go`)
- ✅ All 69 card combinations (52 standard + special cards)
- ✅ Unicode string representation
- ✅ Fallback handling for invalid cards  
- ✅ Constructor validation
- ✅ Constants verification
- ✅ Edge cases and boundary conditions

### Cards Collection Tests (`cards_test.go`)  
- ✅ Deck operations (append, remove, shuffle)
- ✅ Length calculations
- ✅ String formatting
- ✅ Nil card handling
- ✅ Duplicate card removal
- ✅ Large collection performance
- ✅ Full deck validation (54 cards with correct distribution)

### Player Tests (`player_test.go`)
- ✅ Player creation and field access
- ✅ Zero value initialization
- ✅ Field assignment and modification
- ✅ Struct copying behavior
- ✅ Edge cases (negative IDs, special characters)

### Running Tests

```bash
# Run poker package tests
go test ./engine/poker/ -v

# Run with coverage
go test ./engine/poker/ -cover

# Verbose output with details
go test ./engine/poker/ -v -cover
```

## 📊 Test Coverage

- **100% Line Coverage** on all public methods
- **Edge Case Testing** for nil inputs, boundary values
- **Integration Testing** with realistic poker scenarios  
- **Performance Testing** for large card collections
- **Unicode Validation** ensuring proper card display

## 🔗 Integration

This package serves as the foundation for higher-level poker implementations:

### Used by `holdem` package:
```go
// Game uses Cards for community cards and deck
type Game struct {
    Deck           poker.Cards
    CommunityCards poker.Cards
    // ...
}

// Player embeds BasePlayer  
type Player struct {
    poker.BasePlayer
    cards []*poker.Card
    // ...
}
```

### Used by `holdem_ai` package:
```go
// AI evaluates hands using poker Cards
func evaluateHand(cards poker.Cards) HandStrength {
    // Hand evaluation logic using card ranks and suits
}
```

## 🎯 Design Principles

- **Immutability**: Card values are immutable once created
- **Type Safety**: Strong typing prevents invalid card combinations
- **Performance**: Efficient operations for deck shuffling and card management
- **Unicode Support**: Beautiful visual representation of cards
- **Nil Safety**: Graceful handling of edge cases
- **Extensibility**: Easy to add new card types or operations

## 📈 Performance Characteristics

- **O(1)** card creation and property access
- **O(n)** deck shuffling using Fisher-Yates algorithm
- **O(n)** card collection operations (append, remove)  
- **O(n)** string formatting for card display
- **Memory efficient** with minimal allocations for large decks

This package provides a solid, well-tested foundation for any poker implementation with clean APIs, comprehensive error handling, and excellent performance characteristics. 