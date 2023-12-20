package run

import (
	"fmt"
	"os/exec"
	"strings"
)

//We can use open AI to generate commands but to actually run it, we need to use the exec.Command function
//and that's essentially what we're calling here
// RunCommand executes the command and returns the output from execution and any error encountered
func RunCommand(cmd string, arg ...string) (string, error) {
	// Execute the command and get the output
	out, err := exec.Command(cmd, arg...).Output()
	if err != nil {
		return fmt.Sprintf("error: %v", err), err
	}

	// If no error occurred, return the output of the command and nil for the error
	return string(out), nil
}

// PrepareInteractiveCommand prepares a bash command for interactive execution
func PrepareInteractiveCommand(input string) *exec.Cmd {
	// Return a bash command that echoes a newline, executes the input command, and then echoes another newline
	return exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("echo \"\n\";%s; echo \"\n\";", strings.TrimRight(input, ";")),
	)
}

// PrepareEditSettingsCommand prepares a bash command for editing settings
func PrepareEditSettingsCommand(input string) *exec.Cmd {
	// Return a bash command that executes the input command and then echoes a newline
	return exec.Command(
		"bash",
		"-c",
		fmt.Sprintf("%s; echo \"\n\";", strings.TrimRight(input, ";")),
	)
}
