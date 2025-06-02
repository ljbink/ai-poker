# Frontend Texas Hold'em Poker Game

This frontend implements a fully functional Texas Hold'em poker game using the `holdem_ai` engine. The game features a human player vs AI bot in a terminal user interface (TUI) built with Bubble Tea.

## Features

### 🎮 Game Engine Integration
- **Full Texas Hold'em Rules**: Preflop, Flop, Turn, River, and Showdown phases
- **AI Bot Opponent**: Intelligent AI bot with basic poker strategy
- **Real-time Game State**: Live updates of pot, bets, community cards, and player status
- **Action Validation**: Only valid actions are available at any time

### 🎯 Player Actions
- **Fold** (f): Forfeit your hand
- **Check** (c): Pass the action (when no bet to call)
- **Call** (k): Match the current bet
- **Raise** (r): Increase the bet with interactive amount selection
- **All-in** (a): Bet all remaining chips

### 🃏 Game Display
- **Community Cards**: Shows board cards as they're dealt
- **Player Information**: Chip counts, current bets, and status
- **Your Hand**: Your hole cards (bot's cards are hidden)
- **Game Phase**: Current betting round (Preflop, Flop, Turn, River)
- **Action History**: Last action taken by each player

### ⌨️ Controls

#### Main Game Controls
- `f` - Fold
- `c` - Check
- `k` - Call
- `r` - Raise (opens raise amount selector)
- `a` - All-in
- `esc` - Back to menu
- `q` - Quit

#### Raise Amount Selection
- `↑` or `+` - Increase raise amount
- `↓` or `-` - Decrease raise amount
- `enter` - Confirm raise
- `esc` - Cancel raise

## Game Flow

1. **Hand Start**: New hand begins with blinds posted automatically
2. **Player Turn**: When it's your turn, available actions are displayed
3. **Action Selection**: Choose your action using keyboard shortcuts
4. **Bot Turn**: AI bot makes decisions automatically with thinking delay
5. **Phase Progression**: Game advances through betting rounds
6. **Hand Completion**: Winner is determined and chips are distributed
7. **New Hand**: Press any key to start the next hand

## Technical Implementation

### Engine Integration
The view integrates with several key components from the `holdem_ai` package:

- **`holdem.Game`**: Core game state and logic
- **`holdem_ai.GameAdapter`**: Handles action execution
- **`holdem_ai.HumanUserDecisionMaker`**: Human player interface with callback mechanism
- **`holdem_ai.BasicBotDecisionMaker`**: AI bot implementation
- **`holdem_ai.ActionValidator`**: Validates player actions

### Clean Direct Action Architecture
The implementation uses a simplified, direct action architecture:

1. **Turn Detection**: Frontend detects when it's the human player's turn by checking game state
2. **Action Setup**: Frontend gets valid actions directly from `ActionValidator`
3. **User Input**: Frontend captures keyboard input and calls `SendAction()` to send the chosen action
4. **Action Processing**: The decision maker validates and processes the action asynchronously
5. **Game Progression**: Actions are executed through the `GameAdapter` and the game state updates

This eliminates unnecessary complexity and provides a cleaner, more direct interaction pattern.

### Unified Timeout Handling
All timeouts (both human and bot) are handled consistently by the frontend:
- **Human Player**: 30 second timeout for decision making
- **Bot Player**: 5 second timeout for AI decision processing
- **Timeout Action**: Auto-fold if no decision is made within the timeout period

This centralized approach eliminates timeout logic duplication across different decision makers.

### Game States
The view manages several game states:

- **`GameStateWaitingForPlayer`**: Player's turn to act
- **`GameStateRaiseInput`**: Selecting raise amount
- **`GameStateGameOver`**: Hand completed, waiting for new hand
- **`GameStateBotTurn`**: AI bot is making decision

## Example Gameplay

```
Phase: Preflop | Pot: 15 | Current Bet: 10

Community Cards: (hidden)

Players:
  Player: 985 chips, bet 5 [CURRENT TURN]
  AI Bot: 990 chips, bet 10

Your Cards: A♠ K♦

Status: Your turn! Choose an action.
Last Action: AI Bot: call (10)

Available Actions:
  (f) Fold
  (k) Call 5
  (r) Raise
  (a) All-in (985)
```

## Configuration

### Game Settings
- **Small Blind**: 5 chips
- **Big Blind**: 10 chips  
- **Starting Chips**: 1000 chips per player
- **Bot Timeout**: 5 seconds maximum thinking time

### Bot Behavior
The AI bot uses a basic strategy that considers:
- Hand strength evaluation
- Position awareness
- Pot odds calculations
- Betting patterns

## Future Enhancements

Potential improvements for the poker game:
- Multiple bot difficulty levels
- Tournament mode with multiple hands
- Hand history tracking
- Statistical analysis
- Multi-player support
- Custom blind structures
- Save/load game state 