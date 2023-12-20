package ai

//COMPLETE

//the output from execution and chat stream both look different, they're the two modes of the engine
//the exec output has it's own struct defined as EngineExecOutput and the chatStreamOutput has its
//own struct defined below
// EngineExecOutput represents the output of an AI engine execution.
type EngineExecOutput struct {
	Command     string `json:"cmd"`  // Command executed by the AI engine
	Explanation string `json:"exp"`  // Explanation of the command
	Executable  bool   `json:"exec"` // Indicates if the command is executable.
}

//eo being accepted in these methods below is engine exec output
//while co is chat output

// GetCommand returns the command executed by the AI engine.
func (eo EngineExecOutput) GetCommand() string {
	return eo.Command
}

// GetExplanation returns the explanation of the command executed by the AI engine.
func (eo EngineExecOutput) GetExplanation() string {
	return eo.Explanation
}

// IsExecutable returns a boolean indicating if the command is executable.
func (eo EngineExecOutput) IsExecutable() bool {
	return eo.Executable
}

// EngineChatStreamOutput represents the output of an AI engine chat stream.
type EngineChatStreamOutput struct {
	content    string // The content of the chat stream.
	last       bool   // Indicates if this is the last output in the chat stream.
	interrupt  bool   // Indicates if the chat stream was interrupted.
	executable bool   // Indicates if the content is executable.
}

// GetContent returns the content of the chat stream.
func (co EngineChatStreamOutput) GetContent() string {
	return co.content
}

// IsLast returns a boolean indicating if this is the last output in the chat stream.
func (co EngineChatStreamOutput) IsLast() bool {
	return co.last
}

// IsInterrupt returns a boolean indicating if the chat stream was interrupted.
func (co EngineChatStreamOutput) IsInterrupt() bool {
	return co.interrupt
}

// IsExecutable returns a boolean indicating if the content is executable.
func (co EngineChatStreamOutput) IsExecutable() bool {
	return co.executable
}
