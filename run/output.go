package run

import "fmt"

// RunOutput struct holds the error, error message and success message of a run
type RunOutput struct {
	error          error  // error object if any error occurred during the run
	errorMessage   string // custom error message
	successMessage string // custom success message
}

// NewRunOutput is a constructor for RunOutput struct
func NewRunOutput(error error, errorMessage string, successMessage string) RunOutput {
	return RunOutput{
		error:          error,          // set the error
		errorMessage:   errorMessage,   // set the error message
		successMessage: successMessage, // set the success message
	}
}

// HasError checks if the run has an error
func (o RunOutput) HasError() bool {
	return o.error != nil // return true if error is not nil
}

// GetErrorMessage returns the error message of the run
func (o RunOutput) GetErrorMessage() string {
	// format and return the error message with the error
	return fmt.Sprintf("%s: %s", o.errorMessage, o.error)
}

// GetSuccessMessage returns the success message of the run
func (o RunOutput) GetSuccessMessage() string {
	return o.successMessage // return the success message
}
