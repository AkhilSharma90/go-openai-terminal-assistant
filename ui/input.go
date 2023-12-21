package ui
//COMPLETE
import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type UiInput struct {
	runMode    RunMode //runmode has REPL and CLI
	promptMode PromptMode //promptmode has execute, chat and default
	args       string
	pipe       string
}

// NewUIInput is a function that creates a new UiInput instance.
func NewUIInput() (*UiInput, error) {
	// Create a new flag set with the application's name and an error handling.
	//the application's name is available through Args[0] and we set an error policy
	//ExitOnError for this

	//we are using the flag package and flatSet to understand from the user what the user wants
	//to do, if user selects -e flag, means he wants the terminal to enter execute mode
	//whereas if user selects -c, means he wants the terminal to enter chat mode
	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Declare boolean variables for the exec and chat flags.
	var exec, chat bool

	//for our flagSet (which is a NewFlagSet from the flag package), we set exec and chat flags
	// Register the exec and chat flags with the flag set.
	flagSet.BoolVar(&exec, "e", false, "exec prompt mode")
	flagSet.BoolVar(&chat, "c", false, "chat prompt mode")

	//parse is a function present in the flag package, it parses arguments from the command list
	//according to the flag package, this should not contain any commands and this is why the command
	//flag we have set earlier with os.Args[0], now we're setting all the other values other than the command
	//hence starting 1: or 1 onwards
	//It parses the command-line arguments starting from the second argument (os.Args[1:]). 
	//If there's an error parsing the flags, it prints an error message and returns the error.
	err := flagSet.Parse(os.Args[1:])
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return nil, err
	}

//according to the docs - this returns the "non-flag" arguments
//meaning strings etc. that we don't want to process as flags

//if the user enters a non-flag character, meaning neither -c, -e, then we will store it
	args := flagSet.Args()

	// Get the file info for the standard input.
	//stat has the name, size, permission, time etc. 
//input is being considered as a file
	stat, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println("Error getting stat:", err)
		return nil, err
	}

	// Declare a variable for the pipe input.
	//if the user sets a new prompt, we don't need the old one, so we're setting
	//the pipe as empty
	//since this function initializes a new UI for us, we're clearing the pipe
	pipe := ""
	// Check if the standard input is not a named pipe and is empty.
	//if there's some existing data and is not empty, then we will create a reader
	//with bufio.NewReader and read it
	//if size is not zero and pipe is also not zero, then we will read this
	if !(stat.Mode()&os.ModeNamedPipe == 0 && stat.Size() == 0) {
		// Create a new reader for the standard input.
		//the new reader has access to the os.stdin
		reader := bufio.NewReader(os.Stdin)
		// Create a new string builder for operations such as append
		var builder strings.Builder

		// Read runes from the reader until EOF is reached.
		//read the entire data from the file using the reader
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
		//so we can access this later (the input from the user now)
		pipe = strings.TrimSpace(builder.String())
	}

	// Set the run mode to REPL mode by default. chat mode of the engine equates to repl mode of UI
	//if the command is not executable, we stay in repl mode, if it's executable we go to CLI mode
	runMode := ReplMode
	//the non-flag args stored earlier, we will check for them now
	// If there are non-flag arguments, set the run mode to CLI mode. exec mode of engine equates to CLI mode here in UI
	if len(args) > 0 {
		runMode = CliMode
	}

	// Set the prompt mode to default mode by default.
	//prompt mode is a field in the UiInput and we set it here
	//we check if it's exec and not chat, then we set it to execPromptMode
	//else we set it to chatPromptMode
	
	promptMode := DefaultPromptMode
	if exec && !chat {
		promptMode = ExecPromptMode
	} else if !exec && chat {
		promptMode = ChatPromptMode
	}

	// Return a new UiInput instance with the run mode, prompt mode, arguments, and pipe input.
	//all the values set for the fields of this struct above are now being set and returned
	return &UiInput{
		runMode:    runMode,
		promptMode: promptMode,
		args:       strings.Join(args, " "),
		pipe:       pipe,
	}, nil
}

//four helper functions below to get the fields from the struct UiInput defined on top
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
