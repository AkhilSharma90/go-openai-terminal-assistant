package ui

import (
	"fmt"
	"math/rand"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

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
	spinner spinner.Model // The spinner model.
}

// NewSpinner is a function that creates a new Spinner instance.
func NewSpinner() *Spinner {
	// Create a new spinner model.
	spin := spinner.New()
	// Set the spinner style to MiniDot.
	spin.Spinner = spinner.MiniDot

	// Return a new Spinner instance with a random loading message and the spinner model.
	return &Spinner{
		message: loadingMessages[rand.Intn(len(loadingMessages))],
		spinner: spin,
	}
}

// Update is a method on the Spinner struct that updates the spinner model with a message.
func (s *Spinner) Update(msg tea.Msg) (*Spinner, tea.Cmd) {
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
		s.spinner.View(),
		s.spinner.Style.Render(s.message),
	)
}

// Tick is a method on the Spinner struct that returns a tick message for the spinner model.
func (s *Spinner) Tick() tea.Msg {
	return s.spinner.Tick()
}
