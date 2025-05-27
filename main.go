package main

import (
	"fmt"
	"os"

	"github.com/ljbink/ai-poker/frontend"
)

func main() {
	// Start the TUI application
	if err := frontend.RunTUI(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
