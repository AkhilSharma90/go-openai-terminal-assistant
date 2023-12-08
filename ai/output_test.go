package ai

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestEngineExecOutputGetCommand is a test function for testing the GetCommand method of the EngineExecOutput type
func TestEngineExecOutputGetCommand(t *testing.T) {
	eo := EngineExecOutput{Command: "testCommand"}
	result := eo.GetCommand()

	assert.Equal(t, "testCommand", result)
}

// TestEngineExecOutputGetExplanation is a test function for testing the GetExplanation method of the EngineExecOutput type
func TestEngineExecOutputGetExplanation(t *testing.T) {
	eo := EngineExecOutput{Explanation: "testExplanation"}
	result := eo.GetExplanation()

	assert.Equal(t, "testExplanation", result)
}

// TestEngineExecOutputIsExecutable is a test function for testing the IsExecutable method of the EngineExecOutput type
func TestEngineExecOutputIsExecutable(t *testing.T) {
	eo := EngineExecOutput{Executable: true}
	result := eo.IsExecutable()

	assert.True(t, result)
}

// TestEngineChatStreamOutputGetContent is a test function for testing the GetContent method of the EngineChatStreamOutput type
func TestEngineChatStreamOutputGetContent(t *testing.T) {
	co := EngineChatStreamOutput{content: "testContent"}
	result := co.GetContent()

	assert.Equal(t, "testContent", result)
}

// TestEngineChatStreamOutputGetExplanation is a test function for testing the GetExplanation method of the EngineChatStreamOutput type
func TestEngineChatStreamOutputIsLast(t *testing.T) {
	co := EngineChatStreamOutput{last: true}
	result := co.IsLast()

	assert.True(t, result)
}

// TestEngineChatStreamOutputGetExplanation is a test function for testing the GetExplanation method of the EngineChatStreamOutput type
func TestEngineChatStreamOutputIsInterrupt(t *testing.T) {
	co := EngineChatStreamOutput{interrupt: true}
	result := co.IsInterrupt()

	assert.True(t, result)
}

// TestEngineChatStreamOutputGetExplanation is a test function for testing the GetExplanation method of the EngineChatStreamOutput type
func TestEngineChatStreamOutputIsExecutable(t *testing.T) {
	co := EngineChatStreamOutput{executable: true}
	result := co.IsExecutable()

	assert.True(t, result)
}
