package ui

import (
	"fmt"
	"math/rand"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

//spinning basically has logic for when we want to show loading in the terminak
//this could be shown during different messages, which we have defined as loadingMessages variable below

// loadingMessages is a slice of strings that represent different loading messages.
var loadingMessages = []string{
	"let me think",
	"let me see",
	"thinking",
	"loading",
	"hold on",
	"calculating",
	"processing",
	"please wait",
}

// Spinner is a struct that represents a spinner in the user interface.
type Spinner struct {
	message string        // The message to display while the spinner is spinning.
	spinner spinner.Model // The spinner model, got from the bubbles package of charm.sh
}

// NewSpinner is a function that creates a new Spinner instance.
func NewSpinner() *Spinner {
	// Create a new spinner model.
	spin := spinner.New()
	// Set the spinner style to MiniDot.
	//changing this value in the spin variable only as we assign it to spinner in last line of this function
	spin.Spinner = spinner.MiniDot

	// Return a new Spinner instance with a random loading message and the spinner model.
	return &Spinner{
		message: loadingMessages[rand.Intn(len(loadingMessages))],//loadingMessages is a slice we have created above
		//we select a random messages after first calculating the length of loadngMessages, meaning from 0 to n, a random number is picked
		spinner: spin,
	}
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

// Update is a method on the Spinner struct that updates the spinner model with a message.
func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
	//this method takes in Msg of type tea.Msg and returns the spinner and the command
	var updateCmd tea.Cmd
	// Update the spinner model with the message.
	s.spinner, updateCmd = s.spinner.Update(msg)

	return s, updateCmd
}

// View is a method on the Spinner struct that returns a string representation of the spinner.
func (s *Spinner) View() string {
	// Return a string representation of the spinner with the spinner view and the message.
	return fmt.Sprintf(
		"\n  %s %s...",
		//the View function here is getting called from tea
		s.spinner.View(),
		//we add this view with the string message and return from here
		s.spinner.Style.Render(s.message),
	)
}

// Tick is a method on the Spinner struct that returns a tick message for the spinner model.
func (s *Spinner) Tick() tea.Msg {
	//Tick is also a function from tea, tells the spinner what to do in the next second
	return s.spinner.Tick()
}

//>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>