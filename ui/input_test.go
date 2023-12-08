package ui

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUiInput tests the functionality of the UiInput function.
func TestUiInput(t *testing.T) {
	t.Run("NewUIInput", testNewUIInput)
	t.Run("GetRunMode", testGetRunMode)
	t.Run("GetPromptMode", testGetPromptMode)
	t.Run("GetArgs", testGetArgs)
}

// testNewUIInput is a unit test function that tests the NewUIInput function.
// It verifies that NewUIInput does not return an error and that the returned UIInput is not nil.
func testNewUIInput(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-e"}
	uiInput, err := NewUIInput()
	assert.NoError(t, err, "NewUIInput should not return an error.")
	assert.NotNil(t, uiInput, "UiInput should not be nil.")
}

// testGetRunMode is a unit test function that tests the GetRunMode method of the UIInput struct.
// It sets up the necessary environment for the test, including modifying the os.Args slice,
// and then asserts that the returned RunMode is CliMode.
func testGetRunMode(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "arg1", "arg2"}
	uiInput, _ := NewUIInput()
	assert.Equal(t, CliMode, uiInput.GetRunMode(), "RunMode should be CliMode.")
}

// testGetPromptMode is a unit test function that tests the GetPromptMode method of the UIInput struct.
// It verifies that the PromptMode is set to ExecPromptMode when the command-line argument "-e" is provided.
func testGetPromptMode(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "-e"}
	uiInput, _ := NewUIInput()
	assert.Equal(t, ExecPromptMode, uiInput.GetPromptMode(), "PromptMode should be ExecPromptMode.")
}

// testGetArgs is a unit test function that tests the GetArgs method of the UIInput struct.
// It sets the os.Args to a specific value, creates a new UIInput instance, and asserts that the GetArgs method returns the expected result.
func testGetArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"cmd", "arg1", "arg2"}
	uiInput, _ := NewUIInput()
	assert.Equal(t, "arg1 arg2", uiInput.GetArgs(), "Args should be 'arg1 arg2'.")
}
