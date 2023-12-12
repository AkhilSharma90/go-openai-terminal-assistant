package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUserConfig is the main testing function for UserConfig
func TestUserConfig(t *testing.T) {
	// Run the test for GetDefaultPromptMode
	t.Run("GetDefaultPromptMode", testGetDefaultPromptMode)
	// Run the test for GetPreferences
	t.Run("GetPreferences", testGetPreferences)
}

// testGetDefaultPromptMode tests the GetDefaultPromptMode method of UserConfig
func testGetDefaultPromptMode(t *testing.T) {
	// Define the expected default prompt mode
	expectedDefaultPromptMode := "test_mode"
	// Create a UserConfig with the expected default prompt mode
	userConfig := UserConfig{defaultPromptMode: expectedDefaultPromptMode}

	actualDefaultPromptMode := userConfig.GetDefaultPromptMode()

	assert.Equal(t, expectedDefaultPromptMode, actualDefaultPromptMode, "The two default prompt modes should be the same.")
}

// testGetPreferences tests the GetPreferences method of UserConfig
func testGetPreferences(t *testing.T) {
	// Define the expected preferences
	expectedPreferences := "test_preferences"
	// Create a UserConfig with the expected preferences
	userConfig := UserConfig{preferences: expectedPreferences}

	actualPreferences := userConfig.GetPreferences()

	assert.Equal(t, expectedPreferences, actualPreferences, "The two preferences should be the same.")
}
