package ai

type EngineMode int

// Constants representing the different modes of AI engine.
const (
	// ExecEngineMode represents the execution mode of the AI engine.
	ExecEngineMode EngineMode = iota
	// ChatEngineMode represents the chat mode of the AI engine.
	ChatEngineMode
)

// String method returns the string representation of the EngineMode.
func (m EngineMode) String() string {
	// If the mode is ExecEngineMode, return "exec".
	if m == ExecEngineMode {
		return "exec"
	} else {
		// Otherwise, return "chat".
		return "chat"
	}
}
