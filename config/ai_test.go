package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestAiConfig is a test function for testing the AiConfig type
func TestAiConfig(t *testing.T) {
	// Run subtests
	t.Run("GetKey", testGetKey)
	t.Run("GetProxy", testGetProxy)
	t.Run("GetTemperature", testGetTemperature)
	t.Run("GetMaxTokens", testGetMaxTokens)
}

// testGetKey is a subtest function for testing the GetKey method of the AiConfig type
func testGetKey(t *testing.T) {
	// Define the expected key
	expectedKey := "test_key"
	// Create an AiConfig instance with the expected key
	aiConfig := AiConfig{key: expectedKey}

	// Call the GetKey method
	actualKey := aiConfig.GetKey()

	// Assert that the actual key is equal to the expected key
	assert.Equal(t, expectedKey, actualKey, "The two keys should be the same.")
}

// testGetProxy is a subtest function for testing the GetProxy method of the AiConfig type
func testGetProxy(t *testing.T) {
	// Define the expected proxy
	expectedProxy := "test_proxy"
	// Create an AiConfig instance with the expected proxy
	aiConfig := AiConfig{proxy: expectedProxy}

	// Call the GetProxy method
	actualProxy := aiConfig.GetProxy()

	// Assert that the actual proxy is equal to the expected proxy
	assert.Equal(t, expectedProxy, actualProxy, "The two proxies should be the same.")
}

// testGetTemperature is a subtest function for testing the GetTemperature method of the AiConfig type
func testGetTemperature(t *testing.T) {
	// Define the expected temperature
	expectedTemperature := 0.7
	// Create an AiConfig instance with the expected temperature
	aiConfig := AiConfig{temperature: expectedTemperature}

	// Call the GetTemperature method
	actualTemperature := aiConfig.GetTemperature()

	// Assert that the actual temperature is equal to the expected temperature
	assert.Equal(t, expectedTemperature, actualTemperature, "The two temperatures should be the same.")
}

// testGetMaxTokens is a subtest function for testing the GetMaxTokens method of the AiConfig type
func testGetMaxTokens(t *testing.T) {
	expectedMaxTokens := 2000
	aiConfig := AiConfig{maxTokens: expectedMaxTokens}

	actualMaxTokens := aiConfig.GetMaxTokens()

	assert.Equal(t, expectedMaxTokens, actualMaxTokens, "The two maxTokens should be the same.")
}
