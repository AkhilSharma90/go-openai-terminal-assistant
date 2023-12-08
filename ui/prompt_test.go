package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/stretchr/testify/assert"
)

func TestUIPrompt(t *testing.T) {
	t.Run("Prompt", testPrompt)
	t.Run("PromptStyle", testPromptStyle)
	t.Run("PromptIcon", testPromptIcon)
	t.Run("PromptPlaceholder", testPromptPlaceholder)
}

func testPrompt(t *testing.T) {
	// testCases is a slice of test cases for different prompt modes and initial values.
	testCases := []struct {
		name         string
		mode         PromptMode
		initialValue string
	}{
		{"Exec", ExecPromptMode, ""},
		{"Config", ConfigPromptMode, ""},
		{"Chat", ChatPromptMode, ""},
	}

	// Iterate over each test case and run subtests.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPrompt(tc.mode)
			assert.Equal(t, tc.mode, p.GetMode(), "The prompt mode should match the expected value.")
			assert.Equal(t, tc.initialValue, p.GetValue(), "The prompt value should match the expected value.")
		})
	}
}

// testPromptStyle tests the prompt style for different prompt modes.
// It verifies that the prompt style is not nil for each mode.
func testPromptStyle(t *testing.T) {
	testCases := []struct {
		name      string
		mode      PromptMode
		styleFunc func(PromptMode) lipgloss.Style
	}{
		{"Exec", ExecPromptMode, getPromptStyle},
		{"Config", ConfigPromptMode, getPromptStyle},
		{"Chat", ChatPromptMode, getPromptStyle},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			style := tc.styleFunc(tc.mode)
			assert.NotNil(t, style, "The prompt style should not be nil.")
		})
	}
}

// testPromptIcon tests the iconFunc function for different PromptModes.
// It ensures that the prompt icon is not empty for each mode.
func testPromptIcon(t *testing.T) {
	testCases := []struct {
		name     string
		mode     PromptMode
		iconFunc func(PromptMode) string
	}{
		{"Exec", ExecPromptMode, getPromptIcon},
		{"Config", ConfigPromptMode, getPromptIcon},
		{"Chat", ChatPromptMode, getPromptIcon},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			icon := tc.iconFunc(tc.mode)
			assert.NotEmpty(t, icon, "The prompt icon should not be empty.")
		})
	}
}

// testPromptPlaceholder tests the placeholderFunc for different PromptModes.
// It verifies that the prompt placeholder is not empty.
func testPromptPlaceholder(t *testing.T) {
	testCases := []struct {
		name            string
		mode            PromptMode
		placeholderFunc func(PromptMode) string
	}{
		{"Exec", ExecPromptMode, getPromptPlaceholder},
		{"Config", ConfigPromptMode, getPromptPlaceholder},
		{"Chat", ChatPromptMode, getPromptPlaceholder},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			placeholder := tc.placeholderFunc(tc.mode)
			assert.NotEmpty(t, placeholder, "The prompt placeholder should not be empty.")
		})
	}
}
