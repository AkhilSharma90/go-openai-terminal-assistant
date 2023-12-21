package ui

//COMPLETE

//defining a prompmode as an integer
type PromptMode int

// These are the constants representing different prompt modes.
//four prompt modes defined here (exec, chat, config and default)
const (
	ExecPromptMode PromptMode = iota
	ConfigPromptMode
	ChatPromptMode
	DefaultPromptMode
)

// String is a method on the PromptMode type that returns a string representation of the prompt mode.
//this takes in the mode and returns it in string format, the modes are returned by the
//function below GetPromptModelFromString
func (m PromptMode) String() string {
	switch m {
	case ExecPromptMode:
		return "exec"
	case ConfigPromptMode:
		//when we want to set config
		return "config"
	case ChatPromptMode:
		return "chat"
	default:
		//whatever we set as the default in the config file
		return "default"
	}
}

// GetPromptModeFromString is a function that returns the PromptMode from a string.
func GetPromptModeFromString(s string) PromptMode {
	switch s {
	case "exec":
		return ExecPromptMode
	case "config":
		return ConfigPromptMode
	case "chat":
		return ChatPromptMode
	default:
		return DefaultPromptMode
	}
}

// RunMode is an enumerated type that represents different modes of running the application.
type RunMode int

// These are the constants representing different run modes.
//when the user initiates the terminal tool, we are in CLI mode
//and when we start interacting with the terminal when the user is asking if it wants us to
//execute some command on our behalf or not, then it's in REPL mode
const (
	// CliMode is used when the run mode is command-line interface.
	CliMode RunMode = iota
	// ReplMode is used when the run mode is read-eval-print loop.
	ReplMode
)

// String is a method on the RunMode type that returns a string representation of the run mode.
//string method is invoked for a RunMode by writing m.String()
func (m RunMode) String() string {
	if m == CliMode {
		return "cli"
	} else {
		return "repl"
	}
}
