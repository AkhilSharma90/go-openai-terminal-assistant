package ai

//COMPLETE

type EngineMode int

// Constants representing the different modes of AI engine.
// we don't have to write iota everytime for each value when we're defining
//multiple constants, we just have to do it once
//iota is an integer value and commonly used when defining constants
//to generate a series of related values, incremented by 1
const (
	// ExecEngineMode represents the execution mode of the AI engine.
	ExecEngineMode EngineMode = iota
	// ChatEngineMode represents the chat mode of the AI engine.
	ChatEngineMode
)

//Tere are two modes we can operate in - chat mode enables us to chat with the 
//ai model whereas with exec mode we can execute commands
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
