package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ljbink/ai-poker/engine/holdem"
	"github.com/ljbink/ai-poker/engine/holdem_ai"
)

// Example showing how to integrate HumanUserDecisionMaker with frontend
func main() {
	fmt.Println("=== Human Decision Maker Integration Example ===")

	// Create a game
	game := holdem.NewGame([]string{"Human", "Bot"}, 5, 10)
	game.StartHand()

	// Create human decision maker bound to the first player
	humanPlayer := game.Players[0]
	humanDM := holdem_ai.NewHumanUserDecisionMaker(humanPlayer, game)

	// Set up callback to handle action requests from the game engine
	humanDM.SetActionNeededCallback(func(game *holdem.Game, player holdem.IPlayer, validActions []holdem_ai.ActionType) {
		fmt.Println("\n=== ACTION NEEDED ===")
		displayGameState(game, player)
		displayValidActions(validActions)

		// In a real frontend, this would trigger UI updates
		// For this example, we'll simulate user input after a delay
		go simulateUserInput(humanDM, validActions)
	})

	// Set timeout
	humanDM.SetTimeout(10 * time.Second)

	// Simulate the game asking for human decision
	fmt.Printf("\nAsking human player for decision...\n")

	// Get decision channel
	actionChan := humanDM.MakeDecision()

	// Example of caller-side timeout handling (optional)
	select {
	case action := <-actionChan:
		fmt.Printf("\nHuman player chose: %s\n", action.Type.String())
		if action.Amount > 0 {
			fmt.Printf("Amount: %d\n", action.Amount)
		}
	case <-time.After(15 * time.Second):
		fmt.Println("\nCaller timeout - decision took too long!")
		// In real implementation, you might want to cancel the decision
		// or handle this timeout differently
	}

	// Show how to get current game state
	fmt.Println("\n=== CURRENT GAME STATE ===")
	// Access game and player directly since we created the decision maker with them
	displayGameState(game, humanPlayer)

	// Demonstrate bot decision making as well
	fmt.Println("\n=== BOT DECISION EXAMPLE ===")
	botPlayer := game.Players[1]
	botDM := holdem_ai.NewBasicBotDecisionMaker(botPlayer, game)

	fmt.Printf("Bot %s is thinking...\n", botPlayer.GetName())
	botActionChan := botDM.MakeDecision()

	// Wait for bot decision with timeout
	select {
	case botAction := <-botActionChan:
		fmt.Printf("Bot chose: %s\n", botAction.Type.String())
		if botAction.Amount > 0 {
			fmt.Printf("Bot raise amount: %d\n", botAction.Amount)
		}
	case <-time.After(5 * time.Second):
		fmt.Println("Bot decision timeout!")
	}
}

func displayGameState(game *holdem.Game, player holdem.IPlayer) {
	fmt.Printf("Phase: %s\n", holdem_ai.FormatGamePhase(game.CurrentPhase))
	fmt.Printf("Pot: %d\n", game.Pot)
	fmt.Printf("Current Bet: %d\n", game.CurrentBet)
	fmt.Printf("Your Chips: %d\n", player.GetChips())
	fmt.Printf("Your Bet: %d\n", player.GetBet())
	fmt.Printf("Your Cards: %s\n", holdem_ai.FormatPlayerCards(player))
	fmt.Printf("Community Cards: %s\n", holdem_ai.FormatCommunityCards(game))
	fmt.Printf("Call Amount: %d\n", holdem_ai.CalculateCallAmount(game, player))
	fmt.Printf("Pot Odds: %.2f%%\n", game.CalculatePotOdds(player)*100)

	fmt.Println("\nPlayers:")
	for i, p := range game.Players {
		status := ""
		if game.IsPlayerTurn(p) {
			status += " [CURRENT]"
		}
		if p.IsFolded() {
			status += " [FOLDED]"
		}
		fmt.Printf("  %s (pos %d): %d chips, bet %d%s\n",
			p.GetName(), i, p.GetChips(), p.GetBet(), status)
	}
}

func displayValidActions(actions []holdem_ai.ActionType) {
	fmt.Println("\nValid Actions:")
	for i, action := range actions {
		fmt.Printf("  %d. %s\n", i+1, action.String())
	}
}

func simulateUserInput(humanDM *holdem_ai.HumanUserDecisionMaker, validActions []holdem_ai.ActionType) {
	// Simulate thinking time
	time.Sleep(2 * time.Second)

	// For demo purposes, just pick the first valid action
	if len(validActions) > 0 {
		var action holdem_ai.Action

		switch validActions[0] {
		case holdem_ai.ActionFold:
			action = holdem_ai.CreateFoldAction()
		case holdem_ai.ActionCheck:
			action = holdem_ai.CreateCheckAction()
		case holdem_ai.ActionCall:
			action = holdem_ai.CreateCallAction()
		case holdem_ai.ActionRaise:
			action = holdem_ai.CreateRaiseAction(20) // Raise by 20
		case holdem_ai.ActionAllIn:
			action = holdem_ai.CreateAllInAction()
		}

		fmt.Printf("Simulating user input: %s\n", action.Type.String())

		err := humanDM.SendAction(action)
		if err != nil {
			log.Printf("Error sending action: %v", err)
		}
	}
}

// Frontend Integration Guide:
/*
In a real frontend application (TUI/GUI), you would:

1. Create HumanUserDecisionMaker bound to player and game:
   humanDM := holdem_ai.NewHumanUserDecisionMaker(player, game)

2. Set up callback for action requests:
   humanDM.SetActionNeededCallback(func(game *holdem.Game, player holdem.IPlayer, validActions []holdem_ai.ActionType) {
       // Update UI to show game state using utility functions:
       // - holdem_ai.FormatGamePhase(game.CurrentPhase)
       // - holdem_ai.FormatPlayerCards(player)
       // - holdem_ai.FormatCommunityCards(game)
       // - holdem_ai.CalculateCallAmount(game, player)
       // - game.CalculatePotOdds(player)

       // Show available actions as buttons/menu
       // Wait for user input
   })

3. Handle user input (e.g., button clicks):
   func onUserClickFold() {
       action := holdem_ai.CreateFoldAction()
       humanDM.SendAction(action)
   }

   func onUserClickRaise(amount int) {
       action := holdem_ai.CreateRaiseAction(amount)
       humanDM.SendAction(action)
   }

4. Use in game loop with channel-based decision making:
   actionChan := humanDM.MakeDecision()

   // Option 1: Simple blocking wait
   action := <-actionChan

   // Option 2: Wait with timeout
   select {
   case action := <-actionChan:
       // Process action...
   case <-time.After(30 * time.Second):
       // Handle timeout...
   }

   // Option 3: Non-blocking check
   select {
   case action := <-actionChan:
       // Process action...
   default:
       // Still waiting for decision...
   }

Key Features:
- Ultra-minimal interface with only MakeDecision()
- Asynchronous decision making through channels
- Thread-safe action sending
- Built-in timeout handling in decision maker
- Optional caller-side timeout handling
- Action validation
- Callback mechanism for UI updates
- Utility functions for formatting display information
- Convenience action creators
- Each decision maker is bound to specific player and game
*/
