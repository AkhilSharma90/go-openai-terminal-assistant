package ui

type PromptMode int

// These are the constants representing different prompt modes.
const (
	ExecPromptMode PromptMode = iota
	ConfigPromptMode
	ChatPromptMode
	DefaultPromptMode
)

// String is a method on the PromptMode type that returns a string representation of the prompt mode.
func (m PromptMode) String() string {
	switch m {
	case ExecPromptMode:
		return "exec"
	case ConfigPromptMode:
		return "config"
	case ChatPromptMode:
		return "chat"
	default:
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
const (
	// CliMode is used when the run mode is command-line interface.
	CliMode RunMode = iota
	// ReplMode is used when the run mode is read-eval-print loop.
	ReplMode
)

// String is a method on the RunMode type that returns a string representation of the run mode.
func (m RunMode) String() string {
	if m == CliMode {
		return "cli"
	} else {
		return "repl"
	}
}
