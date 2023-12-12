package ui

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUI is a test function that tests the UI package.
func TestUI(t *testing.T) {
	t.Run("PromptModeString", testPromptModeString)
	t.Run("GetPromptModeFromString", testGetPromptModeFromString)
	t.Run("RunModeString", testRunModeString)
}

// testPromptModeString tests the String method of the PromptMode type.
// It verifies that the string representation of each PromptMode value matches the expected value.
func testPromptModeString(t *testing.T) {
	testCases := []struct {
		name       string
		promptMode PromptMode
		expected   string
	}{
		{"Exec", ExecPromptMode, "exec"},
		{"Config", ConfigPromptMode, "config"},
		{"Chat", ChatPromptMode, "chat"},
		{"Default", DefaultPromptMode, "default"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.promptMode.String(), "The string representation should match the expected value.")
		})
	}
}

func testGetPromptModeFromString(t *testing.T) {
	// testCases is a slice of structs that define the test cases for GetPromptModeFromString function.
	testCases := []struct {
		name     string
		input    string
		expected PromptMode
	}{
		{"Exec", "exec", ExecPromptMode},
		{"Config", "config", ConfigPromptMode},
		{"Chat", "chat", ChatPromptMode},
		{"Default", "unknown", DefaultPromptMode},
	}

	// Iterate over each test case and run the sub-test.
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Assert that the actual prompt mode returned by GetPromptModeFromString matches the expected prompt mode.
			assert.Equal(t, tc.expected, GetPromptModeFromString(tc.input), "The prompt mode should match the expected value.")
		})
	}
}

// testRunModeString tests the String method of the RunMode type.
// It verifies that the string representation of each run mode matches the expected value.
func testRunModeString(t *testing.T) {
	testCases := []struct {
		name     string
		runMode  RunMode
		expected string
	}{
		{"CLI", CliMode, "cli"},
		{"REPL", ReplMode, "repl"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.runMode.String(), "The string representation should match the expected value.")
		})
	}
}
