package run

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {
	// TestRun is a unit test for the Run function.
	t.Run("RunCommand", testRunCommand)
	t.Run("PrepareInteractiveCommand", testPrepareInteractiveCommand)
	t.Run("PrepareEditSettingsCommand", testPrepareEditSettingsCommand)
}

// testRunCommand is a unit test for the RunCommand function.
func testRunCommand(t *testing.T) {
	output, err := RunCommand("echo", "Hello, World!")
	require.NoError(t, err)

	assert.Equal(t, "Hello, World!\n", output, "The command output should be the same.")
}

func testPrepareInteractiveCommand(t *testing.T) {
	// testPrepareInteractiveCommand tests the PrepareInteractiveCommand function.
	// It verifies that the generated command matches the expected command with the given input.
	cmd := PrepareInteractiveCommand("echo 'Hello, World!'")

	expectedCmd := exec.Command(
		"bash",
		"-c",
		"echo \"\n\";echo 'Hello, World!'; echo \"\n\";",
	)

	assert.Equal(t, expectedCmd.Args, cmd.Args, "The command arguments should be the same.")
}

// testPrepareEditSettingsCommand is a unit test function that tests the PrepareEditSettingsCommand function.
func testPrepareEditSettingsCommand(t *testing.T) {
	cmd := PrepareEditSettingsCommand("nano yo.json")

	expectedCmd := exec.Command(
		"bash",
		"-c",
		"nano yo.json; echo \"\n\";",
	)

	assert.Equal(t, expectedCmd.Args, cmd.Args, "The command arguments should be the same.")
}
