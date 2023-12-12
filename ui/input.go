package ui

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type UiInput struct {
	runMode    RunMode
	promptMode PromptMode
	args       string
	pipe       string
}

// NewUIInput is a function that creates a new UiInput instance.
func NewUIInput() (*UiInput, error) {
	// Create a new flag set with the application's name and an error handling.
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Declare boolean variables for the exec and chat flags.
	var exec, chat bool

	// Register the exec and chat flags with the flag set.
	flagSet.BoolVar(&exec, "e", false, "exec prompt mode")
	flagSet.BoolVar(&chat, "c", false, "chat prompt mode")

	// Parse the command-line arguments starting from the second argument.
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return nil, err
	}

	args := flagSet.Args()

	// Get the file info for the standard input.
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error getting stat:", err)
		return nil, err
	}

	// Declare a variable for the pipe input.
	pipe := ""
	// Check if the standard input is not a named pipe and is empty.
	if !(stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0) {
		// Create a new reader for the standard input.
		reader := bufio.NewReader(os.Stdin)
		// Create a new string builder.
		var builder strings.Builder

		// Read runes from the reader until EOF is reached.
		for {
			r, _, err := reader.ReadRune()
			if err != nil && err == io.EOF {
				break
			}
			// Write the rune to the string builder.
			_, err = builder.WriteRune(r)
			if err != nil {
				// Print an error message and return if there is an error writing the rune.
				fmt.Println("Error getting input:", err)
				return nil, err
			}
		}

		// Trim the whitespace from the string builder's string and assign it to the pipe variable.
		pipe = strings.TrimSpace(builder.String())
	}

	// Set the run mode to REPL mode by default.
	runMode := ReplMode
	// If there are non-flag arguments, set the run mode to CLI mode.
	if len(args) > 0 {
		runMode = CliMode
	}

	// Set the prompt mode to default mode by default.
	promptMode := DefaultPromptMode
	if exec && !chat {
		promptMode = ExecPromptMode
	} else if !exec && chat {
		promptMode = ChatPromptMode
	}

	// Return a new UiInput instance with the run mode, prompt mode, arguments, and pipe input.
	return &UiInput{
		runMode:    runMode,
		promptMode: promptMode,
		args:       strings.Join(args, " "),
		pipe:       pipe,
	}, nil
}

// GetRunMode is a method that returns the run mode of the UiInput instance.
func (i *UiInput) GetRunMode() RunMode {
	return i.runMode
}

// GetPromptMode is a method that returns the prompt mode of the UiInput instance.
func (i *UiInput) GetPromptMode() PromptMode {
	return i.promptMode
}

// GetArgs is a method that returns the arguments of the UiInput instance.
func (i *UiInput) GetArgs() string {
	return i.args
}

// GetPipe is a method that returns the pipe input of the UiInput instance.
func (i *UiInput) GetPipe() string {
	return i.pipe
}
