package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEngineModeString is a test function for testing the String method of the EngineMode type
func TestEngineModeString(t *testing.T) {
	// Define the test cases
	tests := []struct {
		name     string
		mode     EngineMode
		expected string
	}{
		{
			name:     "ExecEngineMode",
			mode:     ExecEngineMode,
			expected: "exec",
		},
		{
			name:     "ChatEngineMode",
			mode:     ChatEngineMode,
			expected: "chat",
		},
		{
			name:     "UnknownEngineMode",
			mode:     EngineMode(42),
			expected: "chat",
		},
	}

	// Iterate over the test cases
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Call the method under test
			result := test.mode.String()
			// Assert that the result is as expected
			assert.Equal(t, test.expected, result)
		})
	}
}
