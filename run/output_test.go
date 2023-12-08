package run

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunOutput(t *testing.T) {
	t.Run("HasError", testHasError)
	t.Run("GetErrorMessage", testGetErrorMessage)
	t.Run("GetSuccessMessage", testGetSuccessMessage)
}

func testHasError(t *testing.T) {
	// testHasError tests the HasError method of the RunOutput.
	// It verifies that the method correctly determines whether the RunOutput has an error or not.

	err := errors.New("test error")
	runOutputWithError := NewRunOutput(err, "Error occurred", "Success")
	runOutputWithoutError := NewRunOutput(nil, "Error occurred", "Success")

	assert.True(t, runOutputWithError.HasError(), "RunOutput should have an error.")
	assert.False(t, runOutputWithoutError.HasError(), "RunOutput should not have an error.")
}

// testGetErrorMessage is a unit test function that tests the GetErrorMessage method of the RunOutput.
func testGetErrorMessage(t *testing.T) {
	err := errors.New("test error")
	runOutput := NewRunOutput(err, "Error occurred", "Success")

	expectedErrorMessage := "Error occurred: test error"
	actualErrorMessage := runOutput.GetErrorMessage()

	assert.Equal(t, expectedErrorMessage, actualErrorMessage, "The error messages should be the same.")
}

// testGetSuccessMessage is a unit test function that tests the GetSuccessMessage method of the RunOutput.
func testGetSuccessMessage(t *testing.T) {
	runOutput := NewRunOutput(nil, "Error occurred", "Success")

	expectedSuccessMessage := "Success"
	actualSuccessMessage := runOutput.GetSuccessMessage()

	assert.Equal(t, expectedSuccessMessage, actualSuccessMessage, "The success messages should be the same.")
}
