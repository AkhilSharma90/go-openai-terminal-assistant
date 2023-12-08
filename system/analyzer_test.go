package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSystem is a test function that runs subtests for the system package.
// It includes tests for GetOperatingSystem and Analyse functions.
func TestSystem(t *testing.T) {
	t.Run("GetOperatingSystem", testGetOperatingSystem)
	t.Run("Analyse", testAnalyse)
}

// testGetOperatingSystem tests the GetOperatingSystem function.
func testGetOperatingSystem(t *testing.T) {
	operatingSystem := GetOperatingSystem()
	assert.NotEqual(t, UnknownOperatingSystem, operatingSystem, "The operating system should not be unknown.")
}

// testAnalyse is a unit test function that tests the Analyse function.
func testAnalyse(t *testing.T) {
	analysis := Analyse()

	require.NotNil(t, analysis, "Analysis should not be nil.")
	assert.NotEmpty(t, analysis.GetApplicationName(), "Application name should not be empty.")
	assert.NotEqual(t, UnknownOperatingSystem, analysis.GetOperatingSystem(), "The operating system should not be unknown.")
	assert.NotEmpty(t, analysis.GetHomeDirectory(), "Home directory should not be empty.")
	assert.NotEmpty(t, analysis.GetUsername(), "Username should not be empty.")
	assert.NotEmpty(t, analysis.GetConfigFile(), "Config file should not be empty.")
}
