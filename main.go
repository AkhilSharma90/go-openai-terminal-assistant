package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/akhilsharma90/terminal-assistant/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Seed the random number generator with the current time
	rand.Seed(time.Now().UnixNano())

	// Create a new UI input, this function is in input.go
	input, err := ui.NewUIInput()
	if err != nil {
		log.Fatal(err)
	}

	// Create a new UI with the input
	ui := ui.NewUi(input)

	// Run the tea program with the UI
	if _, err := tea.NewProgram(ui).Run(); err != nil {
		log.Fatal(err)
	}
}
