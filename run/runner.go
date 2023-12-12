package run

import (
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand executes a system command and returns its output and any error encountered
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
